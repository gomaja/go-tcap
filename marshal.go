package tcap

import (
	"encoding/asn1"

	"github.com/gomaja/go-tcap/asn1tcapmodel"
)

// care for int optional fields in asn1 structs, to omit a field, its value should be = 255 (FieldOmissionValue)

func marshalAsn1TcapModel(asn1Tcap asn1tcapmodel.TCMessage) ([]byte, error) {
	// Marshal to bytes
	bytes, err := asn1.Marshal(asn1Tcap)
	if err != nil {
		return nil, err
	}

	// reverse operation of parse for "Length has a special treatment"
	var rv asn1.RawValue
	_, err = asn1.Unmarshal(bytes, &rv)
	if err != nil {
		return nil, err
	}

	return rv.Bytes, nil
}

func convertUnidirectionalTCAPToUnidirectional(ud *UnidirectionalTCAP) asn1tcapmodel.Unidirectional {
	var result asn1tcapmodel.Unidirectional

	if ud.Dialogue != nil {
		result.DialoguePortion = convertDialogueTCAPToDialoguePortion(ud.Dialogue)
	}
	result.Components = convertComponentsTCAPToComponentPortion(&ud.Components)

	return result
}

func convertBeginTCAPToBegin(bn *BeginTCAP) asn1tcapmodel.Begin {
	var result asn1tcapmodel.Begin

	// Convert Transaction ID using BigEndian
	result.Otid = asn1tcapmodel.TransactionID(bn.Otid)

	if bn.Dialogue != nil {
		result.DialoguePortion = convertDialogueTCAPToDialoguePortion(bn.Dialogue)
	}
	if bn.Components != nil {
		result.Components = convertComponentsTCAPToComponentPortion(bn.Components)
	}

	return result
}

func convertEndTCAPToEnd(ed *EndTCAP) asn1tcapmodel.End {
	var result asn1tcapmodel.End

	result.Dtid = asn1tcapmodel.TransactionID(ed.Dtid)

	if ed.Dialogue != nil {
		result.DialoguePortion = convertDialogueTCAPToDialoguePortion(ed.Dialogue)
	}
	if ed.Components != nil {
		result.Components = convertComponentsTCAPToComponentPortion(ed.Components)
	}

	return result
}

func convertContinueTCAPToContinue(ct *ContinueTCAP) asn1tcapmodel.Continue {
	var result asn1tcapmodel.Continue

	result.Otid = asn1tcapmodel.TransactionID(ct.Otid)
	result.Dtid = asn1tcapmodel.TransactionID(ct.Dtid)

	if ct.Dialogue != nil {
		result.DialoguePortion = convertDialogueTCAPToDialoguePortion(ct.Dialogue)
	}
	if ct.Components != nil {
		result.Components = convertComponentsTCAPToComponentPortion(ct.Components)
	}

	return result
}

func convertAbortTCAPToAbort(ab *AbortTCAP) asn1tcapmodel.Abort {
	var result asn1tcapmodel.Abort

	result.Dtid = asn1tcapmodel.TransactionID(ab.Dtid)

	if ab.PAbortCause != nil {
		result.PAbortCause = int(*ab.PAbortCause)
	} else {
		result.PAbortCause = FieldOmissionValue // Default value to omit field
	}

	if ab.UAbortCause != nil {
		result.UAbortCause = convertDialogueTCAPToDialoguePortion(ab.UAbortCause)
	}

	return result
}

func convertComponentsTCAPToComponentPortion(cp *ComponentTCAP) asn1tcapmodel.ComponentPortion {
	var result asn1tcapmodel.ComponentPortion

	if cp.Invoke != nil {
		result.Invoke = convertInvokeTCAPToInvoke(cp.Invoke)
	} else if cp.ReturnResultLast != nil {
		result.ReturnResultLast = convertReturnResultTCAPToReturnResult(cp.ReturnResultLast)
	} else if cp.ReturnError != nil {
		result.ReturnError = convertReturnErrorTCAPToReturnError(cp.ReturnError)
	} else if cp.Reject != nil {
		result.Reject = convertRejectTCAPToReject(cp.Reject)
	} else if cp.ReturnResultNotLast != nil {
		result.ReturnResultNotLast = convertReturnResultTCAPToReturnResult(cp.ReturnResultNotLast)
	}

	if cp.MoreComponents != nil {
		moreBytes, err := convertComponentsTCAPToBytes(cp.MoreComponents)
		if err == nil {
			result.MoreComponents = asn1.RawValue{FullBytes: moreBytes}
		}
	}

	return result
}

func convertComponentsTCAPToBytes(cp *ComponentTCAP) ([]byte, error) {
	comp := convertComponentsTCAPToComponentPortion(cp)

	// Marshal to bytes
	bytes, err := asn1.Marshal(comp)
	if err != nil {
		return nil, err
	}

	// reverse operation of parse for "Length has a special treatment"
	var rv asn1.RawValue
	_, err = asn1.Unmarshal(bytes, &rv)
	if err != nil {
		return nil, err
	}

	return rv.Bytes, nil
}

func convertInvokeTCAPToInvoke(inv *InvokeTCAP) asn1tcapmodel.Invoke {
	var result asn1tcapmodel.Invoke

	result.InvokeID = int(inv.InvokeID)
	if inv.LinkedID != nil {
		result.LinkedID = int(*inv.LinkedID)
	} else {
		result.LinkedID = FieldOmissionValue // Default value to omit field
	}
	result.OpCode = int(inv.OpCode)
	if inv.Parameter != nil {
		result.Parameter = asn1.RawValue{FullBytes: inv.Parameter}
	}

	return result
}

func convertReturnResultTCAPToReturnResult(rr *ReturnResultTCAP) asn1tcapmodel.ReturnResult {
	var result asn1tcapmodel.ReturnResult

	var b []byte
	b = append(b, byte(rr.InvokeID))
	result.InvokeID = asn1.RawValue{Tag: asn1.TagInteger, Bytes: b}
	if rr.OpCode != nil {
		result.ResultRetRes = asn1tcapmodel.ResultRetRes{
			OpCode:    int(*rr.OpCode),
			Parameter: asn1.RawValue{FullBytes: rr.Parameter},
		}
	}

	return result
}

func convertReturnErrorTCAPToReturnError(re *ReturnErrorTCAP) asn1tcapmodel.ReturnError {
	var result asn1tcapmodel.ReturnError

	result.InvokeID = int(re.InvokeID)
	result.ErrorCode = int(re.ErrorCode)
	if re.Parameter != nil {
		result.Parameter = asn1.RawValue{FullBytes: re.Parameter}
	}

	return result
}

func convertRejectTCAPToReject(rj *RejectTCAP) asn1tcapmodel.Reject {
	var result asn1tcapmodel.Reject

	if rj.DerivableOrNotDerivable != nil {
		result.DerivableOrNotDerivable = asn1.RawValue{
			Class: asn1.ClassUniversal,
			Tag:   asn1.TagInteger,
			Bytes: []byte{*rj.DerivableOrNotDerivable},
		}
	} else {
		result.DerivableOrNotDerivable = asn1.RawValue{
			Class: asn1.ClassUniversal,
			Tag:   asn1.TagNull,
			Bytes: []byte{},
		}
	}

	if rj.GeneralProblem != nil {
		result.GeneralProblem = int(*rj.GeneralProblem)
	} else {
		result.GeneralProblem = FieldOmissionValue
	}
	if rj.InvokeProblem != nil {
		result.InvokeProblem = int(*rj.InvokeProblem)
	} else {
		result.InvokeProblem = FieldOmissionValue
	}
	if rj.ReturnResultProblem != nil {
		result.ReturnResultProblem = int(*rj.ReturnResultProblem)
	} else {
		result.ReturnResultProblem = FieldOmissionValue
	}
	if rj.ReturnErrorProblem != nil {
		result.ReturnErrorProblem = int(*rj.ReturnErrorProblem)
	} else {
		result.ReturnErrorProblem = FieldOmissionValue
	}

	return result
}

func convertDialogueTCAPToDialoguePortion(dp *DialogueTCAP) asn1tcapmodel.DialoguePortion {
	var diagAll asn1tcapmodel.DialogueAll

	if dp.DialogAsId != nil {
		diagAll.DialogueAsId = asn1.ObjectIdentifier(dp.DialogAsId)
	}

	if dp.DialogueRequest != nil {
		diagAll.DialoguePDU.DialogueRequest = convertAARQapduTCAPToAARQapdu(dp.DialogueRequest)
	} else if dp.DialogueResponse != nil {
		diagAll.DialoguePDU.DialogueResponse = convertAAREapduTCAPToAAREapdu(dp.DialogueResponse)
	} else if dp.DialogueAbort != nil {
		diagAll.DialoguePDU.DialogueAbort = convertABRTapduTCAPToABRTapdu(dp.DialogueAbort)
	}

	bytes, _ := asn1.Marshal(diagAll)
	bytes[0] = ExternalConstructorTag // overwrite "sequence constructor" tag to the "external constructor" tag
	return asn1tcapmodel.DialoguePortion{
		Data: asn1.RawValue{FullBytes: bytes},
	}
}

func convertAARQapduTCAPToAARQapdu(aarq *AARQapduTCAP) asn1tcapmodel.AARQapdu {
	var result asn1tcapmodel.AARQapdu

	if aarq.ProtocolVersionPadded != nil {
		result.ProtocolVersionPadded = asn1.BitString{
			Bytes:     []byte{*aarq.ProtocolVersionPadded},
			BitLength: 1,
		}
	}

	if aarq.AcnVersion != nil {
		result.ApplicationContextName = asn1.ObjectIdentifier(aarq.AcnVersion)
	}

	if aarq.UserInformation != nil {
		aarq.UserInformation[0] = ExternalConstructorTag // overwrite "sequence constructor" tag to the "external constructor" tag
		result.UserInformation = asn1tcapmodel.UserInformation{
			Data: asn1.RawValue{FullBytes: aarq.UserInformation},
		}
	}

	return result
}

func convertAAREapduTCAPToAAREapdu(aare *AAREapduTCAP) asn1tcapmodel.AAREapdu {
	var result asn1tcapmodel.AAREapdu

	if aare.ProtocolVersionPadded != nil {
		result.ProtocolVersionPadded = asn1.BitString{
			Bytes:     []byte{*aare.ProtocolVersionPadded},
			BitLength: 1,
		}
	}

	if aare.AcnVersion != nil {
		result.ApplicationContextName = asn1.ObjectIdentifier(aare.AcnVersion)
	}

	result.Result = asn1tcapmodel.AssociateResult{Data: int(aare.Result)}

	if aare.ResultSourceDiagnostic.DialogueServiceUser != nil {
		result.ResultSourceDiagnostic = asn1tcapmodel.AssociateSourceDiagnostic{
			DialogueServiceUser:     int(*aare.ResultSourceDiagnostic.DialogueServiceUser),
			DialogueServiceProvider: FieldOmissionValue, // omit the field with this default number (magic number used in asn1 structs)
		}
	}

	if aare.ResultSourceDiagnostic.DialogueServiceProvider != nil {
		result.ResultSourceDiagnostic = asn1tcapmodel.AssociateSourceDiagnostic{
			DialogueServiceUser:     FieldOmissionValue, // omit the field with this default number (magic number used in asn1 structs)
			DialogueServiceProvider: int(*aare.ResultSourceDiagnostic.DialogueServiceProvider),
		}
	}

	if aare.UserInformation != nil {
		aare.UserInformation[0] = ExternalConstructorTag // overwrite the "sequence constructor" tag to the "external constructor" tag
		result.UserInformation = asn1tcapmodel.UserInformation{
			Data: asn1.RawValue{FullBytes: aare.UserInformation},
		}
	}

	return result
}

func convertABRTapduTCAPToABRTapdu(abrt *ABRTapduTCAP) asn1tcapmodel.ABRTapdu {
	var result asn1tcapmodel.ABRTapdu

	result.AbortSource = int(abrt.AbortSource)

	if abrt.UserInformation != nil {
		abrt.UserInformation[0] = ExternalConstructorTag // overwrite "sequence constructor" tag to the "external constructor" tag
		result.UserInformation = asn1tcapmodel.UserInformation{
			Data: asn1.RawValue{FullBytes: abrt.UserInformation},
		}
	}

	return result
}

/*
// uint32ToBytes converts a Go uint to a big-endian byte slice,
// omitting all left (most significant) zero bytes except when the value is zero.
func uint32ToBytes(u uint32) []byte {
	// Special case: if u == 0, return a single zero byte
	if u == 0 {
		return []byte{0}
	}

	// Convert to uint32 for consistent 4-byte handling,
	// then serialize in BigEndian.
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], u)

	// Find the first non-zero byte
	idx := 0
	for idx < len(buf) && buf[idx] == 0 {
		idx++
	}

	// Slice from the first non-zero byte to the end
	return buf[idx:]
}
*/
