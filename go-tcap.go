package tcap

import "fmt"

// TCAP will have only one field fulfilled and the others will be nil
type TCAP struct { // choice
	Unidirectional *UnidirectionalTCAP
	Begin          *BeginTCAP
	End            *EndTCAP
	Continue       *ContinueTCAP
	Abort          *AbortTCAP
}

type UnidirectionalTCAP struct {
	Dialogue   *DialogueTCAP
	Components ComponentTCAP
}
type BeginTCAP struct {
	Otid       TransactionID
	Dialogue   *DialogueTCAP
	Components *ComponentTCAP
}
type EndTCAP struct {
	Dtid       TransactionID
	Dialogue   *DialogueTCAP
	Components *ComponentTCAP
}
type ContinueTCAP struct {
	Otid       TransactionID
	Dtid       TransactionID
	Dialogue   *DialogueTCAP
	Components *ComponentTCAP
}
type AbortTCAP struct {
	Dtid        TransactionID
	PAbortCause *uint8
	UAbortCause *DialogueTCAP
}

// TransactionID size from 1 up to 4 bytes in BigEndian format
type TransactionID []byte

type DialogueTCAP struct {
	DialogAsId []int

	DialogueRequest  *AARQapduTCAP
	DialogueResponse *AAREapduTCAP
	DialogueAbort    *ABRTapduTCAP
}

type AARQapduTCAP struct {
	ProtocolVersionPadded *uint8
	AcnVersion            []int
	UserInformation       []byte
}
type AAREapduTCAP struct {
	ProtocolVersionPadded  *uint8
	AcnVersion             []int
	Result                 uint8
	ResultSourceDiagnostic AssociateSourceDiagnostic
	UserInformation        []byte
}

type AssociateSourceDiagnostic struct {
	DialogueServiceUser     *uint8
	DialogueServiceProvider *uint8
}

type ABRTapduTCAP struct {
	AbortSource     uint8
	UserInformation []byte
}

// ComponentTCAP will have only one field fulfilled and the others will be nil, except MoreComponents may exist additionally to any other field
type ComponentTCAP struct { // choice Invoke, ReturnResultLast, ReturnError, Reject, ReturnResultNotLast
	Invoke              *InvokeTCAP
	ReturnResultLast    *ReturnResultTCAP
	ReturnError         *ReturnErrorTCAP
	Reject              *RejectTCAP
	ReturnResultNotLast *ReturnResultTCAP

	// Linked list here to include presence of more than one component
	MoreComponents *ComponentTCAP
}

type InvokeTCAP struct {
	InvokeID  int // integer value range -128 to 127
	LinkedID  *int
	OpCode    uint8
	Parameter []byte
}
type ReturnResultTCAP struct {
	InvokeID  int // integer value range -128 to 127
	OpCode    *uint8
	Parameter []byte
}
type ReturnErrorTCAP struct {
	InvokeID  int // integer value range -128 to 127
	ErrorCode uint8
	Parameter []byte
}
type RejectTCAP struct {
	DerivableOrNotDerivable *uint8 // if value is nil means it is NotDerivable
	GeneralProblem          *uint8
	InvokeProblem           *uint8
	ReturnResultProblem     *uint8
	ReturnErrorProblem      *uint8
}

// validateTransactionID validates that a transaction ID meets ITU-T Q.773 requirements
// Transaction ID must be 1 to 4 bytes in length
func validateTransactionID(tid []byte, fieldName string) error {
	if len(tid) < MinTransactionIDLength || len(tid) > MaxTransactionIDLength {
		return newValidationError(fieldName, len(tid),
			fmt.Errorf("must be %d to %d bytes in length, got %d bytes",
				MinTransactionIDLength, MaxTransactionIDLength, len(tid)))
	}
	return nil
}

// validateInvokeID validates that an invoke ID is within the valid range
// Invoke ID must be in range -128 to 127 (signed 8-bit integer)
func validateInvokeID(invID int, fieldName string) error {
	if invID < MinInvokeID || invID > MaxInvokeID {
		return newValidationError(fieldName, invID,
			fmt.Errorf("must be in range %d to %d, got %d",
				MinInvokeID, MaxInvokeID, invID))
	}
	return nil
}
