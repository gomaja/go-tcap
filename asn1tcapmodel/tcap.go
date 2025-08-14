package asn1tcapmodel

import (
	"encoding/asn1"
	"fmt"
)

type TCMessage struct { // CHOICE
	Unidirectional Unidirectional `asn1:"application,tag:1,optional"`
	Begin          Begin          `asn1:"application,tag:2,optional"`
	End            End            `asn1:"application,tag:4,optional"`
	Continue       Continue       `asn1:"application,tag:5,optional"`
	Abort          Abort          `asn1:"application,tag:7,optional"`
}

type Unidirectional struct {
	DialoguePortion DialoguePortion `asn1:"application,tag:11,optional"`
	Components      ComponentPortion
}

type Begin struct {
	Otid            TransactionID    `asn1:"application,tag:8"`
	DialoguePortion DialoguePortion  `asn1:"application,tag:11,optional"`
	Components      ComponentPortion `asn1:"application,tag:12,optional"`
}

type End struct {
	Dtid            TransactionID    `asn1:"application,tag:9"`
	DialoguePortion DialoguePortion  `asn1:"application,tag:11,optional"`
	Components      ComponentPortion `asn1:"application,tag:12,optional"`
}

type Continue struct {
	Otid            TransactionID    `asn1:"application,tag:8"`
	Dtid            TransactionID    `asn1:"application,tag:9"`
	DialoguePortion DialoguePortion  `asn1:"application,tag:11,optional"`
	Components      ComponentPortion `asn1:"application,tag:12,optional"`
}

type Abort struct {
	Dtid TransactionID `asn1:"application,tag:9"`

	// PAbortCauseInt // default:255 will change omitting the optional field to the value 255 instead of 0, the 255 (FieldOmissionValue) is used to check if the field arrived was empty
	PAbortCause int `asn1:"default:255,application,tag:10,optional"` // [APPLICATION 10] = 0x2a = 42 dec

	UAbortCause DialoguePortion `asn1:"application,tag:11,optional"` // [APPLICATION 11]
}

type DialoguePortion struct {
	// this is an EXTERNAL (tag of EXTERNAL is 8), make sure to add this tag with constructor type when marshalling (final tag EXTERNAL + Constructor = 0x28)
	Data asn1.RawValue
}

type TransactionID []byte // OCTET STRING // size from 1 up to 4 bytes in BigEndian format

type PAbortCauseInt int

const (
	UnrecognizedMessageType          PAbortCauseInt = 0
	UnrecognizedTransactionID        PAbortCauseInt = 1
	BadlyFormattedTransactionPortion PAbortCauseInt = 2
	IncorrectTransactionPortion      PAbortCauseInt = 3
	ResourceLimitation               PAbortCauseInt = 4
)

// String() helper for debugging
func (p PAbortCauseInt) String() string {
	switch p {
	case UnrecognizedMessageType:
		return "unrecognizedMessageType"
	case UnrecognizedTransactionID:
		return "unrecognizedTransactionID"
	case BadlyFormattedTransactionPortion:
		return "badlyFormattedTransactionPortion"
	case IncorrectTransactionPortion:
		return "incorrectTransactionPortion"
	case ResourceLimitation:
		return "resourceLimitation"
	default:
		return fmt.Sprintf("unknown(%d)", p)
	}
}

type ComponentPortion struct { // CHOICE
	Invoke              Invoke       `asn1:"tag:1,optional"`
	ReturnResultLast    ReturnResult `asn1:"tag:2,optional"`
	ReturnError         ReturnError  `asn1:"tag:3,optional"`
	Reject              Reject       `asn1:"tag:4,optional"`
	ReturnResultNotLast ReturnResult `asn1:"tag:7,optional"`

	// This is added here to treat cases showing more than one component, component number 2 and above are filled in
	//the below struct, component number 1 is already filled in one of the previous fields
	MoreComponents asn1.RawValue `asn1:"optional"`
}

type Invoke struct {
	// InvokeIDInt // integer value range -128 to 127
	InvokeID int // InvokeID InvokeIDInt starts at 0, and may increment to 1, 2, ... for long SMS

	// LinkedIDInt // default:255 will change omitting the optional field to the value 255 instead of 0, the 255 (FieldOmissionValue) is used to check if the field arrived was empty
	LinkedID int `asn1:"default:255,tag:0,optional"` // [ContextSpecific 0]

	// OpCodeInt
	OpCode int

	// Parameter is ANY; represent as raw bytes or asn1.RawValue in Go.
	Parameter asn1.RawValue `asn1:"optional"`
}

type ReturnResult struct {
	// InvokeIDInt // this is an int // initialize using: var b []byte ; b = append(b, byte(rr.InvokeID)) ; asn1.RawValue{Tag: asn1.TagInteger, Bytes: b}
	InvokeID asn1.RawValue // RawValue is used to treat zero value avoiding that it provides zero value which indicates omitted

	ResultRetRes ResultRetRes `asn1:"optional"`
}

type ResultRetRes struct {
	// OpCodeInt
	OpCode int

	Parameter asn1.RawValue
}

type ReturnError struct {
	// InvokeIDInt
	InvokeID int

	// ErrorCodeInt
	ErrorCode int

	Parameter asn1.RawValue `asn1:"optional"`
}

type Reject struct {
	//RejectInvokeIDRej // CHOICE
	//value must be: asn1.RawValue{Class: asn1.ClassUniversal, Tag: asn1.TagInteger, Bytes: []byte{byte(DerivableInt)}}
	//Derivable asn1.RawValue `asn1:"optional"` // the actual InvokeIdType
	//value must be: asn1.RawValue{Class: asn1.ClassUniversal, Tag: asn1.TagNull, Bytes: []byte{byte(DerivableInt)}}
	//NotDerivable asn1.RawValue `asn1:"optional"` // represents NULL // tag = 5 dec
	DerivableOrNotDerivable asn1.RawValue // fill it here and verify if Derivable or NotDerivable after

	//RejectProblem // CHOICE
	// GeneralProblemInt // default:255 will change omitting the optional field to value 255 instead of 0, the 255 (FieldOmissionValue) is used to check if the field arrived was empty
	GeneralProblem int `asn1:"default:255,tag:0,optional"`
	// InvokeProblemInt // default:255 will change omitting the optional field to value 255 instead of 0, the 255 (FieldOmissionValue) is used to check if the field arrived was empty
	InvokeProblem int `asn1:"default:255,tag:1,optional"`
	// ReturnResultProblemInt // default:255 will change omitting the optional field to value 255 instead of 0, the 255 (FieldOmissionValue) is used to check if the field arrived was empty
	ReturnResultProblem int `asn1:"default:255,tag:2,optional"`
	// ReturnErrorProblemInt // default:255 will change omitting the optional field to value 255 instead of 0, the 255 (FieldOmissionValue) is used to check if the field arrived was empty
	ReturnErrorProblem int `asn1:"default:255,tag:3,optional"`
}

type GeneralProblemInt int

const (
	UnrecognizedComponent    GeneralProblemInt = 0
	MistypedComponent        GeneralProblemInt = 1
	BadlyStructuredComponent GeneralProblemInt = 2
)

type InvokeProblemInt int

const (
	DuplicateInvokeID         InvokeProblemInt = 0
	UnrecognizedOperation     InvokeProblemInt = 1
	MistypedParameterInvoke   InvokeProblemInt = 2
	ResourceLimitationInvoke  InvokeProblemInt = 3
	InitiatingRelease         InvokeProblemInt = 4
	UnrecognizedLinkedID      InvokeProblemInt = 5
	LinkedResponseUnexpected  InvokeProblemInt = 6
	UnexpectedLinkedOperation InvokeProblemInt = 7
)

type ReturnResultProblemInt int

const (
	UnrecognizedInvokeIDResult ReturnResultProblemInt = 0
	ReturnResultUnexpected     ReturnResultProblemInt = 1
	MistypedParameterResult    ReturnResultProblemInt = 2
)

type ReturnErrorProblemInt int

const (
	UnrecognizedInvokeIDError ReturnErrorProblemInt = 0
	ReturnErrorUnexpected     ReturnErrorProblemInt = 1
	UnrecognizedError         ReturnErrorProblemInt = 2
	UnexpectedError           ReturnErrorProblemInt = 3
	MistypedParameterError    ReturnErrorProblemInt = 4
)
