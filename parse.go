package tcap

import (
	"encoding/asn1"
	"reflect"
	"strings"

	"github.com/gomaja/go-tcap/asn1tcapmodel"
	"github.com/gomaja/go-tcap/utils"
)

func uint8Ptr(i int) *uint8 {
	u := uint8(i)
	return &u
}

func ParseAny(b []byte) (TCAP, error) {
	// Parse first with ParseDER and check for error
	tcap, err := ParseDER(b)
	if err != nil && strings.Contains(err.Error(), IndefiniteLengthErrorString) {
		// If ParseDER showed an error due to indefinite-length type (not DER), then convert to DER
		derBytes, err := utils.MakeDER(b) // bytes will remain the same if the provided bytes are already DER, otherwise they are converted to from non-DER (indefinite length) to DER
		if err != nil {
			return nil, newParseError("ParseAny", "MakeDER", err)
		}

		return ParseDER(derBytes)
	}

	return tcap, err
}

// ParseDER takes a slice of byte, without a tag and length. A tag byte and length byte are added at the beginning internally to satisfy asn1 Unmarshal behavior for completeness
// It parses on DER encoded asn1 structs, the encoding/asn1 package parses on DER (indefinite length is not supported)
func ParseDER(b []byte) (TCAP, error) {
	if len(b) == 0 {
		return nil, ErrEmptyData
	}

	// Length have a special treatment, either we can write a custom function to treat short form and long form of lengths as per ITU-T Q.773 (06/97) page 10,
	//or we can make a marshal to let the Go "encoding/binary" package do it, hence encapsulating the input byte to the proper one that can be understood by "encoding/binary"
	newByteString := asn1.RawValue{Tag: asn1.TagSequence, IsCompound: true, Bytes: b} // Tag = 16 with Constructor = 0x30

	// marshal for proper length (short form or long form) and tag treatment (to add struct tag)
	b, err := asn1.Marshal(newByteString)
	if err != nil {
		return nil, newParseError("ParseDER", "Marshal", err)
	}

	var asn1Tcap asn1tcapmodel.TCMessage
	_, err = asn1.Unmarshal(b, &asn1Tcap)
	if err != nil {
		return nil, newParseError("ParseDER", "Unmarshal", err)
	}

	return convertTCMessageToTCAP(&asn1Tcap), nil
}

func convertTCMessageToTCAP(tcMessage *asn1tcapmodel.TCMessage) TCAP {
	var tcTCAP TCAP

	// CHOICE
	// fulfill one of the: Unidirectional, Begin, End, Continue, Abort

	// check Unidirectional and fill tc
	if !reflect.ValueOf(tcMessage.Unidirectional).IsZero() { // check if Unidirectional is not empty
		tcTCAP = convertUnidirectionalToUnidirectionalTCAP(&tcMessage.Unidirectional)
		return tcTCAP
	}

	// check Begin and fill tc
	if !reflect.ValueOf(tcMessage.Begin).IsZero() { // check if Begin is not empty
		tcTCAP = convertBeginToBeginTCAP(&tcMessage.Begin)
		return tcTCAP
	}

	// check End and fill tc
	if !reflect.ValueOf(tcMessage.End).IsZero() { // check if End is not empty
		tcTCAP = convertEndToEndTCAP(&tcMessage.End)
		return tcTCAP
	}

	// check Continue and fill tc
	if !reflect.ValueOf(tcMessage.Continue).IsZero() { // check if Continue is not empty
		tcTCAP = convertContinueToContinueTCAP(&tcMessage.Continue)
		return tcTCAP
	}

	// check Abort and fill tc
	if !reflect.ValueOf(tcMessage.Abort).IsZero() { // check if Abort is not empty
		tcTCAP = convertAbortToAbortTCAP(&tcMessage.Abort)
		return tcTCAP
	}

	return nil
}

func convertUnidirectionalToUnidirectionalTCAP(ud *asn1tcapmodel.Unidirectional) *UnidirectionalTCAP {
	var udTcap UnidirectionalTCAP

	// fill Dialogue
	if !reflect.ValueOf(ud.DialoguePortion).IsZero() { // checks and fulfills optional field
		udTcap.Dialogue = convertDialoguePortionToDialogueTCAP(&ud.DialoguePortion)
	}

	// fill Components
	udTcap.Components = *convertComponentPortionToComponentsTCAP(&ud.Components)

	return &udTcap
}

func convertBeginToBeginTCAP(bn *asn1tcapmodel.Begin) *BeginTCAP {
	var bnTcap BeginTCAP

	// fill Transaction
	bnTcap.Otid = TransactionID(bn.Otid)

	// fill Dialogue
	if !reflect.ValueOf(bn.DialoguePortion).IsZero() { // checks and fulfills optional field
		bnTcap.Dialogue = convertDialoguePortionToDialogueTCAP(&bn.DialoguePortion)
	}

	// fill Components
	if !reflect.ValueOf(bn.Components).IsZero() { // checks and fulfills optional field
		bnTcap.Components = convertComponentPortionToComponentsTCAP(&bn.Components)
	}

	return &bnTcap
}

func convertEndToEndTCAP(ed *asn1tcapmodel.End) *EndTCAP {
	var edTcap EndTCAP

	// fill Transaction
	edTcap.Dtid = TransactionID(ed.Dtid)

	// fill Dialogue
	if !reflect.ValueOf(ed.DialoguePortion).IsZero() { // checks and fulfills optional field
		edTcap.Dialogue = convertDialoguePortionToDialogueTCAP(&ed.DialoguePortion)
	}

	// fill Components
	if !reflect.ValueOf(ed.Components).IsZero() { // checks and fulfills optional field
		edTcap.Components = convertComponentPortionToComponentsTCAP(&ed.Components)
	}

	return &edTcap
}

func convertContinueToContinueTCAP(ctn *asn1tcapmodel.Continue) *ContinueTCAP {
	var ctnTcap ContinueTCAP

	// fill Transaction
	ctnTcap.Dtid = TransactionID(ctn.Dtid)
	ctnTcap.Otid = TransactionID(ctn.Otid)

	// fill Dialogue
	if !reflect.ValueOf(ctn.DialoguePortion).IsZero() { // checks and fulfills optional field
		ctnTcap.Dialogue = convertDialoguePortionToDialogueTCAP(&ctn.DialoguePortion)
	}

	// fill Components
	if !reflect.ValueOf(ctn.Components).IsZero() { // checks and fulfills optional field
		ctnTcap.Components = convertComponentPortionToComponentsTCAP(&ctn.Components)
	}

	return &ctnTcap
}

func convertAbortToAbortTCAP(abrt *asn1tcapmodel.Abort) *AbortTCAP {
	var abrtTcap AbortTCAP

	// fill Transaction
	abrtTcap.Dtid = TransactionID(abrt.Dtid)

	// fill PAbortCause
	if abrt.PAbortCause != FieldOmissionValue { // checks and fulfills optional field
		abrtTcap.PAbortCause = uint8Ptr(int(abrt.PAbortCause))
	}

	// fill UAbortCause
	if !reflect.ValueOf(abrt.UAbortCause).IsZero() { // checks and fulfills optional field
		abrtTcap.UAbortCause = convertDialoguePortionToDialogueTCAP(&abrt.UAbortCause)
	}

	return &abrtTcap
}

func convertComponentPortionToComponentsTCAP(cmpnt *asn1tcapmodel.ComponentPortion) *ComponentTCAP {
	var components ComponentTCAP

	// CHOICE
	// fulfill one of the: Invoke, ReturnResultLast, ReturnError, Reject, ReturnResultNotLast

	// check Invoke and fill tc
	if !reflect.ValueOf(cmpnt.Invoke).IsZero() { // check if Invoke is not empty
		components.Invoke = convertInvokeToInvokeTCAP(&cmpnt.Invoke)

		// recursive filling
		components.TreatRecursiveMoreComponents(cmpnt)

		return &components
	}

	// check ReturnResultLast and fill tc
	if !reflect.ValueOf(cmpnt.ReturnResultLast).IsZero() { // check if ReturnResultLast is not empty
		components.ReturnResultLast = convertReturnResultToReturnResultTCAP(&cmpnt.ReturnResultLast)

		// recursive filling
		components.TreatRecursiveMoreComponents(cmpnt)

		return &components
	}

	// check ReturnError and fill tc
	if !reflect.ValueOf(cmpnt.ReturnError).IsZero() { // check if ReturnError is not empty
		components.ReturnError = convertReturnErrorToReturnErrorTCAP(&cmpnt.ReturnError)

		// recursive filling
		components.TreatRecursiveMoreComponents(cmpnt)

		return &components
	}

	// check Reject and fill tc
	if !reflect.ValueOf(cmpnt.Reject).IsZero() { // check if Reject is not empty
		components.Reject = convertRejectToRejectTCAP(&cmpnt.Reject)

		// recursive filling
		components.TreatRecursiveMoreComponents(cmpnt)

		return &components
	}

	// check ReturnResultNotLast and fill tc
	if !reflect.ValueOf(cmpnt.ReturnResultNotLast).IsZero() { // check if ReturnResultNotLast is not empty
		components.ReturnResultNotLast = convertReturnResultToReturnResultTCAP(&cmpnt.ReturnResultNotLast)

		// recursive filling
		components.TreatRecursiveMoreComponents(cmpnt)

		return &components
	}

	return nil
}

func (components *ComponentTCAP) TreatRecursiveMoreComponents(cmpnt *asn1tcapmodel.ComponentPortion) {
	if !reflect.ValueOf(cmpnt.MoreComponents).IsZero() { // check if MoreComponents is not empty
		cmpnt2 := &asn1tcapmodel.ComponentPortion{}

		fullBytes := cmpnt.MoreComponents.FullBytes
		// Length have a special treatment, either we can write a custom function to treat short/long forms of length as per ITU-T Q.773 (06/97) page 10
		//or we can make a marshal to let the Go "encoding/binary" package do it, hence encapsulating the input byte to the proper one that can be understood by "encoding/binary"
		newByteString := asn1.RawValue{Tag: asn1.TagSequence, IsCompound: true, Bytes: fullBytes} // Tag = 16 with Constructor = 0x30

		// marshal for proper length (short form or long form) and tag treatment (to add struct tag)
		b, err := asn1.Marshal(newByteString)
		if err != nil {
			return
		}

		_, err = asn1.Unmarshal(b, cmpnt2)
		if err == nil {
			components.MoreComponents = convertComponentPortionToComponentsTCAP(cmpnt2)
		}
	}
}

func convertInvokeToInvokeTCAP(invk *asn1tcapmodel.Invoke) *InvokeTCAP {
	var invokeTcap InvokeTCAP

	invokeTcap.InvokeID = invk.InvokeID

	if invk.LinkedID != FieldOmissionValue { // check optional field if empty
		invokeTcap.LinkedID = &invk.LinkedID
	}

	invokeTcap.OpCode = uint8(invk.OpCode)

	if !reflect.ValueOf(invk.Parameter).IsZero() { // check optional field if empty
		invokeTcap.Parameter = invk.Parameter.FullBytes
	}

	return &invokeTcap
}

func convertReturnResultToReturnResultTCAP(rr *asn1tcapmodel.ReturnResult) *ReturnResultTCAP {
	var rrlTcap ReturnResultTCAP

	rrlTcap.InvokeID = int(rr.InvokeID.Bytes[0])

	if !reflect.ValueOf(rr.ResultRetRes).IsZero() { // check optional field if empty
		rrlTcap.OpCode = uint8Ptr(rr.ResultRetRes.OpCode)
		rrlTcap.Parameter = rr.ResultRetRes.Parameter.FullBytes
	}

	return &rrlTcap
}

func convertReturnErrorToReturnErrorTCAP(re *asn1tcapmodel.ReturnError) *ReturnErrorTCAP {
	var reTcap ReturnErrorTCAP

	reTcap.InvokeID = re.InvokeID

	reTcap.ErrorCode = uint8(re.ErrorCode)

	if !reflect.ValueOf(re.Parameter).IsZero() { // check optional field if empty
		reTcap.Parameter = re.Parameter.FullBytes
	}

	return &reTcap
}

func convertRejectToRejectTCAP(rj *asn1tcapmodel.Reject) *RejectTCAP {
	var rjTcap RejectTCAP

	if rj.DerivableOrNotDerivable.Tag == asn1.TagInteger { // check optional field if empty
		rjTcap.DerivableOrNotDerivable = &rj.DerivableOrNotDerivable.Bytes[len(rj.DerivableOrNotDerivable.Bytes)-1]
	}
	// added for clarity // if rjTcap.DerivableOrNotDerivable is nil, then it is NotDerivable
	if rj.DerivableOrNotDerivable.Tag == asn1.TagNull { // check optional field if empty
		rjTcap.DerivableOrNotDerivable = nil
	}

	if rj.GeneralProblem != FieldOmissionValue { // check optional field if empty
		rjTcap.GeneralProblem = uint8Ptr(int(rj.GeneralProblem))
	}

	if rj.InvokeProblem != FieldOmissionValue { // check optional field if empty
		rjTcap.InvokeProblem = uint8Ptr(int(rj.InvokeProblem))
	}

	if rj.ReturnResultProblem != FieldOmissionValue { // check optional field if empty
		rjTcap.ReturnResultProblem = uint8Ptr(int(rj.ReturnResultProblem))
	}

	if rj.ReturnErrorProblem != FieldOmissionValue { // check optional field if empty
		rjTcap.ReturnErrorProblem = uint8Ptr(int(rj.ReturnErrorProblem))
	}

	return &rjTcap
}

func convertDialoguePortionToDialogueTCAP(dp *asn1tcapmodel.DialoguePortion) *DialogueTCAP {
	var dpTCAP DialogueTCAP

	// create a asn1tcapmodel.DialogAll
	var DiagAll asn1tcapmodel.DialogueAll
	// modify the tag to a SEQUENCE tag with constructor type (the ASN1 library in go works like this, SEQUENCE tag should be for structs with constructor type)
	DiagAllBytes := dp.Data.FullBytes
	DiagAllBytes[0] = SequenceConstructorTag // overwrite the EXTERNAL tag (EXTERNAL Constructor) to the known struct tag for asn1
	_, _ = asn1.Unmarshal(DiagAllBytes, &DiagAll)

	if !reflect.ValueOf(DiagAll.DialogueAsId).IsZero() { // check optional field if empty
		dpTCAP.DialogAsId = []int(DiagAll.DialogueAsId)
	}

	if !reflect.ValueOf(DiagAll.DialoguePDU.DialogueRequest).IsZero() { // check optional field if empty
		dpTCAP.DialogueRequest = convertAARQapduToAARQapduTCAP(&DiagAll.DialoguePDU.DialogueRequest)
	}

	if !reflect.ValueOf(DiagAll.DialoguePDU.DialogueResponse).IsZero() { // check optional field if empty
		dpTCAP.DialogueResponse = convertAAREapduToAAREapduTCAP(&DiagAll.DialoguePDU.DialogueResponse)
	}

	if !reflect.ValueOf(DiagAll.DialoguePDU.DialogueAbort).IsZero() { // check optional field if empty
		dpTCAP.DialogueAbort = convertABRTapduToABRTapduTCAP(&DiagAll.DialoguePDU.DialogueAbort)
	}

	return &dpTCAP
}

func convertAARQapduToAARQapduTCAP(aarq *asn1tcapmodel.AARQapdu) *AARQapduTCAP {
	var aarqTcap AARQapduTCAP

	// protocol version
	if !reflect.ValueOf(aarq.ProtocolVersionPadded).IsZero() { // check optional field if empty
		aarqTcap.ProtocolVersionPadded = &aarq.ProtocolVersionPadded.Bytes[len(aarq.ProtocolVersionPadded.Bytes)-1]
	}

	// ACN version
	aarqTcap.AcnVersion = aarq.ApplicationContextName

	// user information
	if !reflect.ValueOf(aarq.UserInformation).IsZero() { // check optional field if empty
		aarqTcap.UserInformation = aarq.UserInformation.Data.FullBytes
	}

	return &aarqTcap
}

func convertAAREapduToAAREapduTCAP(aare *asn1tcapmodel.AAREapdu) *AAREapduTCAP {
	var aareTcap AAREapduTCAP

	// protocol version
	if !reflect.ValueOf(aare.ProtocolVersionPadded).IsZero() { // check optional field if empty
		aareTcap.ProtocolVersionPadded = &aare.ProtocolVersionPadded.Bytes[len(aare.ProtocolVersionPadded.Bytes)-1]
	}

	// ACN version
	aareTcap.AcnVersion = aare.ApplicationContextName

	// Result
	aareTcap.Result = uint8(aare.Result.Data)

	// ResultSourceDiagnostic // CHOICE
	if aare.ResultSourceDiagnostic.DialogueServiceUser != FieldOmissionValue { // check optional field if empty
		aareTcap.ResultSourceDiagnostic.DialogueServiceUser = uint8Ptr(int(aare.ResultSourceDiagnostic.DialogueServiceUser))
	}
	if aare.ResultSourceDiagnostic.DialogueServiceProvider != FieldOmissionValue { // check optional field if empty
		aareTcap.ResultSourceDiagnostic.DialogueServiceProvider = uint8Ptr(int(aare.ResultSourceDiagnostic.DialogueServiceProvider))
	}

	// user information
	if !reflect.ValueOf(aare.UserInformation).IsZero() { // check optional field if empty
		aareTcap.UserInformation = aare.UserInformation.Data.FullBytes
	}

	return &aareTcap
}

func convertABRTapduToABRTapduTCAP(abrt *asn1tcapmodel.ABRTapdu) *ABRTapduTCAP {
	var abrtTcap ABRTapduTCAP

	// ACN version
	abrtTcap.AbortSource = uint8(abrt.AbortSource)

	// user information
	if !reflect.ValueOf(abrt.UserInformation).IsZero() { // check optional field if empty
		abrtTcap.UserInformation = abrt.UserInformation.Data.FullBytes
	}
	return &abrtTcap
}
