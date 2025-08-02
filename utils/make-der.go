package utils

import (
	"fmt"
)

// ASN1Object represents an ASN.1 object
type ASN1Object struct {
	Tag           byte
	Value         []byte
	IsConstructed bool
	SubObjects    []ASN1Object
}

// MakeDER converts non-DER (i.e., any encoding that contains an indefinite-length type) encoded asn1 struct bytes to DER encoded bytes - encodings with an indefinite-length type are not supported by encoding/asn1, so prior conversion is required
// If the given bytes are already DER encoded, the bytes will remain the same without any change
func MakeDER(bytes []byte) ([]byte, error) {
	if len(bytes) == 0 {
		return nil, fmt.Errorf("empty input")
	}

	obj, _, err := parseBytesWithIndefiniteLength(bytes, 0)
	if err != nil {
		return nil, err
	}

	return encodeDER(obj)
}

func parseBytesWithIndefiniteLength(data []byte, pos int) (ASN1Object, int, error) {
	if pos >= len(data) {
		return ASN1Object{}, pos, fmt.Errorf("unexpected end of data")
	}

	tag := data[pos]
	pos++
	isConstructed := (tag & 0x20) != 0

	// Handle length
	if pos >= len(data) {
		return ASN1Object{}, pos, fmt.Errorf("unexpected end of data")
	}

	var length int
	var err error
	length, pos, err = parseLength(data, pos)
	if err != nil {
		return ASN1Object{}, pos, err
	}

	// Handle indefinite length
	if length == -1 {
		// For indefinite length, we need to find EOC (End-Of-Contents) marker
		return parseIndefiniteLength(data, pos, tag)
	}

	if pos+length > len(data) {
		return ASN1Object{}, pos, fmt.Errorf("length exceeds data bounds")
	}

	obj := ASN1Object{
		Tag:           tag,
		IsConstructed: isConstructed,
	}

	if isConstructed {
		endPos := pos + length
		for pos < endPos {
			subObj, newPos, err := parseBytesWithIndefiniteLength(data, pos)
			if err != nil {
				return ASN1Object{}, pos, err
			}
			obj.SubObjects = append(obj.SubObjects, subObj)
			pos = newPos
		}
	} else {
		obj.Value = data[pos : pos+length]
		pos += length
	}

	return obj, pos, nil
}

func parseLength(data []byte, pos int) (int, int, error) {
	if pos >= len(data) {
		return 0, pos, fmt.Errorf("unexpected end of data")
	}

	length := int(data[pos])
	pos++

	if length&0x80 == 0 {
		return length, pos, nil
	}

	// Handle multi-byte length
	numLengthBytes := length & 0x7F
	if numLengthBytes == 0 {
		// Indefinite length
		return -1, pos, nil
	}

	if pos+int(numLengthBytes) > len(data) {
		return 0, pos, fmt.Errorf("unexpected end of data")
	}

	length = 0
	for i := 0; i < int(numLengthBytes); i++ {
		length = length<<8 | int(data[pos])
		pos++
	}

	return length, pos, nil
}

func parseIndefiniteLength(data []byte, startPos int, tag byte) (ASN1Object, int, error) {
	obj := ASN1Object{
		Tag:           tag,
		IsConstructed: true,
	}

	pos := startPos
	for pos+1 < len(data) {
		// Check for EOC
		if data[pos] == 0x00 && data[pos+1] == 0x00 {
			pos += 2
			return obj, pos, nil
		}

		subObj, newPos, err := parseBytesWithIndefiniteLength(data, pos)
		if err != nil {
			return ASN1Object{}, pos, err
		}
		obj.SubObjects = append(obj.SubObjects, subObj)
		pos = newPos
	}

	return ASN1Object{}, pos, fmt.Errorf("no EOC marker found")
}

func encodeDER(obj ASN1Object) ([]byte, error) {
	var content []byte
	var err error

	if obj.IsConstructed {
		content, err = encodeConstructedDER(obj)
		if err != nil {
			return nil, err
		}
	} else {
		// Handle boolean normalization for DER
		// In TCAP, context-specific tag 1 (0x81) is commonly used for boolean values
		// and must be normalized according to DER rules
		if len(obj.Value) == 1 && (obj.Tag == 0x01 || obj.Tag == 0x81) {
			// Boolean values must be 0x00 (false) or 0xFF (true) in DER
			// Any non-zero value is normalized to 0xFF
			if obj.Value[0] != 0x00 {
				content = []byte{0xFF}
			} else {
				content = []byte{0x00}
			}
		} else {
			content = obj.Value
		}
	}

	// Encode length
	lengthBytes := encodeDERLength(len(content))

	// Combine all parts
	result := make([]byte, 0, 1+len(lengthBytes)+len(content))
	result = append(result, obj.Tag)
	result = append(result, lengthBytes...)
	result = append(result, content...)

	return result, nil
}

func encodeConstructedDER(obj ASN1Object) ([]byte, error) {
	var content []byte

	//// For SET OF, we need to sort the elements
	//if obj.Tag == 0x31 {
	//	// Sort SubObjects based on their DER encoding
	//	// Implementation of sorting omitted for brevity
	//}

	// Encode each sub-object
	for _, subObj := range obj.SubObjects {
		subContent, err := encodeDER(subObj)
		if err != nil {
			return nil, err
		}
		content = append(content, subContent...)
	}

	return content, nil
}

func encodeDERLength(length int) []byte {
	if length < 128 {
		return []byte{byte(length)}
	}

	// Convert length to minimum number of bytes
	var lengthBytes []byte
	for length > 0 {
		lengthBytes = append([]byte{byte(length & 0xFF)}, lengthBytes...)
		length >>= 8
	}

	// Add length of length bytes with high bit set
	return append([]byte{byte(0x80 | len(lengthBytes))}, lengthBytes...)
}
