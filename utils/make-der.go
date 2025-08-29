package utils

import (
	"bytes"
	"fmt"
	"sort"
)

// ASN1Object represents an ASN.1 element (potentially constructed with sub-elements)
type ASN1Object struct {
	Class         int // 0=Universal, 1=Application, 2=Context-specific, 3=Private
	TagNumber     int // The tag number (could be > 30 for extended tags)
	IsConstructed bool
	Value         []byte       // Used if primitive (holds the content bytes)
	SubObjects    []ASN1Object // Used if constructed (holds the parsed sub-elements)
}

// MakeDER converts ASN.1 bytes that may use BER indefinite lengths (or other non-DER forms) into DER-encoded bytes.
// If the input is already in DER form, it is returned unchanged (aside from necessary canonicalizations like boolean normalization or set ordering).
func MakeDER(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	// Parse the ASN.1 structure (supporting indefinite lengths)
	obj, pos, err := parseElement(data, 0)
	if err != nil {
		return nil, err
	}
	if pos != len(data) {
		return nil, fmt.Errorf("extra data found after parsing ASN.1 structure")
	}
	// Encode the parsed ASN.1 object back into DER form
	return encodeDER(obj)
}

// parseElement parses an ASN.1 element starting at data[pos], handling definite and indefinite lengths.
// It returns the ASN1Object, the new position (after this element), or an error.
func parseElement(data []byte, pos int) (ASN1Object, int, error) {
	if pos >= len(data) {
		return ASN1Object{}, pos, fmt.Errorf("unexpected end of data")
	}
	// --- Parse Tag ---
	firstByte := data[pos]
	pos++
	// Determine class and constructed flag from first byte
	tagClass := int((firstByte & 0xC0) >> 6) // class bits (bits 8-7)
	isConstructed := (firstByte & 0x20) != 0 // bit 6
	tagNumber := int(firstByte & 0x1F)       // lower 5 bits
	if tagNumber == 0x1F {
		// Long-form tag encoding
		tagNumber = 0
		var octetCount int
		for {
			if pos >= len(data) {
				return ASN1Object{}, pos, fmt.Errorf("truncated tag bytes")
			}
			octet := data[pos]
			pos++
			octetCount++
			// Append lower 7 bits to tagNumber
			tagNumber = (tagNumber << 7) | int(octet&0x7F)
			// If MSB is 0, this was the last tag byte
			if octet&0x80 == 0 {
				break
			}
			// Safety check: ASN.1 tag bytes are limited (to avoid overflow)
			if octetCount > 6 {
				// (6 bytes => 42 bits for tag number, which is extremely large; adjust as needed)
				return ASN1Object{}, pos, fmt.Errorf("tag number is too large")
			}
		}
	}
	// Create object shell
	obj := ASN1Object{
		Class:         tagClass,
		TagNumber:     tagNumber,
		IsConstructed: isConstructed,
	}
	// --- Parse Length ---
	if pos >= len(data) {
		return ASN1Object{}, pos, fmt.Errorf("unexpected end of data after tag")
	}
	lengthByte := data[pos]
	pos++
	var length int
	if lengthByte&0x80 == 0 {
		// Short form length
		length = int(lengthByte)
	} else {
		// Long form or indefinite
		numLenBytes := int(lengthByte & 0x7F)
		if numLenBytes == 0 {
			// Indefinite length marker (0x80)
			length = -1
		} else {
			if pos+numLenBytes > len(data) {
				return ASN1Object{}, pos, fmt.Errorf("length bytes exceed data size")
			}
			// Decode big-endian length
			length = 0
			for i := 0; i < numLenBytes; i++ {
				length = (length << 8) | int(data[pos+i])
			}
			pos += numLenBytes
		}
	}
	// --- Parse Value / SubObjects ---
	if length == -1 {
		// Indefinite length: must be constructed (per BER rules):contentReference[oaicite:9]{index=9}.
		if !isConstructed {
			return ASN1Object{}, pos, fmt.Errorf("indefinite length used for a primitive tag (tag class %d, number %d)", tagClass, tagNumber)
		}
		// Parse sub-elements until EOC (0x00 0x00)
		for {
			if pos+1 <= len(data) && data[pos] == 0x00 && data[pos+1] == 0x00 {
				// End-of-Contents marker
				pos += 2
				break
			}
			subObj, newPos, err := parseElement(data, pos)
			if err != nil {
				return ASN1Object{}, pos, err
			}
			obj.SubObjects = append(obj.SubObjects, subObj)
			pos = newPos
		}
	} else {
		// Definite length
		if pos+length > len(data) {
			return ASN1Object{}, pos, fmt.Errorf("specified length %d exceeds remaining data", length)
		}
		if isConstructed {
			// Parse `length` bytes as sub-elements
			endPos := pos + length
			for pos < endPos {
				subObj, newPos, err := parseElement(data, pos)
				if err != nil {
					return ASN1Object{}, pos, err
				}
				obj.SubObjects = append(obj.SubObjects, subObj)
				pos = newPos
			}
			// Sanity: pos should now equal endPos
			if pos != endPos {
				return ASN1Object{}, pos, fmt.Errorf("length mismatch for constructed tag (expected end at %d, got %d)", endPos, pos)
			}
		} else {
			// Primitive: take the value bytes
			obj.Value = data[pos : pos+length]
			pos += length
		}
	}
	return obj, pos, nil
}

// encodeDER encodes an ASN1Object (with all substructure parsed) into DER-compliant bytes.
func encodeDER(obj ASN1Object) ([]byte, error) {
	// If the object is constructed, first handle any special DER rules for constructed types
	if obj.IsConstructed {
		// **DER Flattening**: If this is a constructed form of a normally primitive type (e.g. Octet String, Bit String, or other strings),
		// we need to flatten it into a primitive encoding.
		if obj.Class == 0 { // universal class
			switch obj.TagNumber {
			case 3: // BIT STRING
				// All subobjects should be primitive bit string fragments
				var bitData []byte
				var finalUnusedBits byte = 0
				for i, sub := range obj.SubObjects {
					if !(sub.Class == 0 && sub.TagNumber == 3 && !sub.IsConstructed) {
						return nil, fmt.Errorf("BIT STRING constructed with invalid sub-element at index %d", i)
					}
					if len(sub.Value) == 0 {
						return nil, fmt.Errorf("BIT STRING sub-element %d has no content", i)
					}
					// The first byte of each BIT STRING content is the 'unused bits' count
					unused := sub.Value[0]
					dataBytes := sub.Value[1:]
					if i < len(obj.SubObjects)-1 {
						// Not the last chunk
						if unused != 0 {
							return nil, fmt.Errorf("intermediate BIT STRING chunk %d has non-zero unused bits (%d)", i, unused)
						}
						// Append all bits (full bytes) from this chunk
						bitData = append(bitData, dataBytes...)
					} else {
						// Last chunk: record its unused bits and append its data
						finalUnusedBits = unused
						bitData = append(bitData, dataBytes...)
					}
				}
				// Construct the flattened BIT STRING content: one byte for final unused bits count + concatenated data
				obj.IsConstructed = false
				obj.Value = append([]byte{finalUnusedBits}, bitData...)
				obj.SubObjects = nil // clear subobjects since now primitive
			case 4: // OCTET STRING
				var octets []byte
				for i, sub := range obj.SubObjects {
					if !(sub.Class == 0 && sub.TagNumber == 4 && !sub.IsConstructed) {
						return nil, fmt.Errorf("OCTET STRING constructed with invalid sub-element at index %d", i)
					}
					octets = append(octets, sub.Value...)
				}
				obj.IsConstructed = false
				obj.Value = octets
				obj.SubObjects = nil
			case 12, 18, 19, 20, 21, 22, 25, 26, 27, 28, 30:
				// Various string types: UTF8String(12), NumericString(18), PrintableString(19), TeletexString(20),
				// VideotexString(21), IA5String(22), GraphicString(25), VisibleString(26), GeneralString(27),
				// UniversalString(28), BMPString(30) – flatten similarly
				var strContent []byte
				for i, sub := range obj.SubObjects {
					if !(sub.Class == 0 && sub.TagNumber == obj.TagNumber && !sub.IsConstructed) {
						return nil, fmt.Errorf("Constructed string (tag %d) has invalid sub-element at index %d", obj.TagNumber, i)
					}
					strContent = append(strContent, sub.Value...)
				}
				obj.IsConstructed = false
				obj.Value = strContent
				obj.SubObjects = nil
			case 17:
				// SET (or SET OF) – inherently constructed, do not flatten, but will sort below
				// (Handled in sorting section)
			default:
				// No special handling for other constructed universal types (e.g., SEQUENCE is fine as is)
			}
		}
	}

	var contentBytes []byte

	if obj.IsConstructed {
		// Handle constructed content encoding (e.g., Sequence, Set, or context/application containers)
		// If SET OF/SET, sort the sub-encodings for DER:contentReference[oaicite:10]{index=10}:
		if obj.Class == 0 && obj.TagNumber == 17 {
			// Encode each subobject to DER and sort them lexicographically
			encodedSubs := make([][]byte, len(obj.SubObjects))
			for i, sub := range obj.SubObjects {
				subEnc, err := encodeDER(sub)
				if err != nil {
					return nil, err
				}
				encodedSubs[i] = subEnc
			}
			sort.Slice(encodedSubs, func(i, j int) bool {
				return bytes.Compare(encodedSubs[i], encodedSubs[j]) < 0
			})
			// Concatenate sorted encodings
			for _, enc := range encodedSubs {
				contentBytes = append(contentBytes, enc...)
			}
		} else {
			// Not a SET, just encode subobjects in given order
			for _, sub := range obj.SubObjects {
				subEnc, err := encodeDER(sub)
				if err != nil {
					return nil, err
				}
				contentBytes = append(contentBytes, subEnc...)
			}
		}
	} else {
		// Primitive content
		contentBytes = obj.Value
		// **Boolean normalization**: If this is a boolean, ensure value is 0x00 or 0xFF:contentReference[oaicite:11]{index=11}.
		// Universal BOOLEAN (tag 1) or context-specific boolean (commonly tag 1 in context).
		if len(contentBytes) == 1 {
			if (obj.Class == 0 && obj.TagNumber == 1) || (obj.Class == 2 && obj.TagNumber == 1) {
				// Normalize TRUE to 0xFF, FALSE is 0x00
				if contentBytes[0] != 0x00 {
					contentBytes = []byte{0xFF}
				} else {
					contentBytes = []byte{0x00}
				}
			}
		}
	}

	// --- Encode Tag ---
	var tagBytes []byte
	// Construct the first identifier octet from class and constructed flag
	var firstOctet byte = byte((obj.Class & 0x3) << 6)
	if obj.IsConstructed {
		firstOctet |= 0x20 // set constructed bit
	}
	if obj.TagNumber < 31 {
		// Fits in short form
		firstOctet |= byte(obj.TagNumber & 0x1F)
		tagBytes = []byte{firstOctet}
	} else {
		// Extended tag form
		firstOctet |= 0x1F
		tagBytes = []byte{firstOctet}
		// Encode tag number in base-128 (7-bit) big-endian, with MSB = 1 on all but last byte
		tagNum := obj.TagNumber
		// Determine bytes needed:
		var tmp []byte
		for {
			tmp = append([]byte{byte(tagNum & 0x7F)}, tmp...)
			tagNum >>= 7
			if tagNum == 0 {
				break
			}
		}
		// Set MSB on all but last in tmp
		for i := 0; i < len(tmp)-1; i++ {
			tmp[i] |= 0x80
		}
		// Append the tag number bytes
		tagBytes = append(tagBytes, tmp...)
	}

	// --- Encode Length ---
	length := len(contentBytes)
	var lengthBytes []byte
	if length < 128 {
		lengthBytes = []byte{byte(length)}
	} else {
		// Encode length in big-endian bytes, without leading zeros
		// Determine number of bytes needed for length
		temp := length
		var lenOctets []byte
		for temp > 0 {
			lenOctets = append([]byte{byte(temp & 0xFF)}, lenOctets...)
			temp >>= 8
		}
		// Prefix with length of length, with MSB = 1
		numLenBytes := len(lenOctets)
		lengthBytes = append([]byte{0x80 | byte(numLenBytes)}, lenOctets...)
	}

	// Combine tag, length, and content
	encoded := make([]byte, 0, len(tagBytes)+len(lengthBytes)+len(contentBytes))
	encoded = append(encoded, tagBytes...)
	encoded = append(encoded, lengthBytes...)
	encoded = append(encoded, contentBytes...)
	return encoded, nil
}
