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

// NewBegin create a Begin tcap message
// otid size from 1 to 4 bytes in BigEndian format
func NewBegin(otid []byte) (*TCAP, error) {
	if err := validateTransactionID(otid, "otid"); err != nil {
		return nil, err
	}
	return NewBeginWithDialogue(otid, nil, nil), nil
}

// NewBeginWithDialogue create a Begin tcap message with a dialogue
// otid size from 1 to 4 bytes in BigEndian format
func NewBeginWithDialogue(otid []byte, acn *int, acnVersion *int) *TCAP {
	tcTcap := &TCAP{}
	tcTcap.Begin = &BeginTCAP{}

	tcTcap.Begin.Otid = otid

	if acn != nil && acnVersion != nil {
		tcTcap.Begin.Dialogue = &DialogueTCAP{}
		tcTcap.Begin.Dialogue.DialogueRequest = &AARQapduTCAP{}

		tcTcap.Begin.Dialogue.DialogAsId = []int{0, 0, 17, 773, 1, 1, 1}
		tcTcap.Begin.Dialogue.DialogueRequest.ProtocolVersionPadded = uint8Ptr(128) // protocol version 128 = 0x80
		tcTcap.Begin.Dialogue.DialogueRequest.AcnVersion = []int{0, 4, 0, 0, 1, 0, *acn, *acnVersion}
	}

	return tcTcap
}

// NewBeginInvoke create a Begin Invoke tcap message
// otid size from 1 to 4 bytes in BigEndian format
func NewBeginInvoke(otid []byte, invID int, opCode uint8, payload []byte) (*TCAP, error) {
	if err := validateTransactionID(otid, "otid"); err != nil {
		return nil, err
	}
	if err := validateInvokeID(invID, "invID"); err != nil {
		return nil, err
	}
	return NewBeginInvokeWithDialogue(otid, invID, opCode, payload, nil, nil), nil
}

// NewBeginInvokeWithDialogue create a Begin Invoke tcap message with a dialogue
// otid size from 1 to 4 bytes in BigEndian format
func NewBeginInvokeWithDialogue(otid []byte, invID int, opCode uint8, payload []byte, acn *int, acnVersion *int) *TCAP {
	tcTcap := &TCAP{}
	tcTcap.Begin = &BeginTCAP{}
	tcTcap.Begin.Components = &ComponentTCAP{}
	tcTcap.Begin.Components.Invoke = &InvokeTCAP{}

	tcTcap.Begin.Otid = otid
	tcTcap.Begin.Components.Invoke.InvokeID = invID
	tcTcap.Begin.Components.Invoke.OpCode = opCode
	tcTcap.Begin.Components.Invoke.Parameter = payload

	if acn != nil && acnVersion != nil {
		tcTcap.Begin.Dialogue = &DialogueTCAP{}
		tcTcap.Begin.Dialogue.DialogueRequest = &AARQapduTCAP{}

		tcTcap.Begin.Dialogue.DialogAsId = []int{0, 0, 17, 773, 1, 1, 1}
		tcTcap.Begin.Dialogue.DialogueRequest.ProtocolVersionPadded = uint8Ptr(128) // protocol version 128 = 0x80
		tcTcap.Begin.Dialogue.DialogueRequest.AcnVersion = []int{0, 4, 0, 0, 1, 0, *acn, *acnVersion}
	}

	return tcTcap
}

// NewEndReturnResultLast create an End ReturnResultLast tcap message
// dtid size from 1 to 4 bytes in BigEndian format
func NewEndReturnResultLast(dtid []byte, invID int, opCode *uint8, payload []byte) (*TCAP, error) {
	if err := validateTransactionID(dtid, "dtid"); err != nil {
		return nil, err
	}
	if err := validateInvokeID(invID, "invID"); err != nil {
		return nil, err
	}
	return NewEndReturnResultLastWithDialogue(dtid, invID, opCode, payload, nil, nil), nil
}

// NewEndReturnResultLastWithDialogue create an End ReturnResultLast tcap message with a dialogue
// dtid size from 1 to 4 bytes in BigEndian format
func NewEndReturnResultLastWithDialogue(dtid []byte, invID int, opCode *uint8, payload []byte, acn *int, acnVersion *int) *TCAP {
	tcTcap := &TCAP{}
	tcTcap.End = &EndTCAP{}
	tcTcap.End.Components = &ComponentTCAP{}
	tcTcap.End.Components.ReturnResultLast = &ReturnResultTCAP{}

	tcTcap.End.Dtid = dtid
	tcTcap.End.Components.ReturnResultLast.InvokeID = invID
	tcTcap.End.Components.ReturnResultLast.OpCode = opCode
	tcTcap.End.Components.ReturnResultLast.Parameter = payload

	if acn != nil && acnVersion != nil {
		tcTcap.End.Dialogue = &DialogueTCAP{}
		tcTcap.End.Dialogue.DialogueRequest = &AARQapduTCAP{}

		tcTcap.End.Dialogue.DialogAsId = []int{0, 0, 17, 773, 1, 1, 1}
		tcTcap.End.Dialogue.DialogueRequest.ProtocolVersionPadded = uint8Ptr(128) // protocol version 128 = 0x80
		tcTcap.End.Dialogue.DialogueRequest.AcnVersion = []int{0, 4, 0, 0, 1, 0, *acn, *acnVersion}
	}

	return tcTcap
}

// NewEndReturnResultLastWithDialogueObject creates an End ReturnResultLast tcap message with a dialogue object
// Parameters:
//   - dtid: Destination Transaction ID, size from 1 to 4 bytes in BigEndian format.
//   - dialogueObject: A pointer to a DialogueTCAP object, representing the dialogue to include in the message.
//     If nil, no dialogue will be included in the message.
func NewEndReturnResultLastWithDialogueObject(dtid []byte, invID int, opCode *uint8, payload []byte, dialogueObject *DialogueTCAP) *TCAP {
	tcTcap := &TCAP{}
	tcTcap.End = &EndTCAP{}
	tcTcap.End.Components = &ComponentTCAP{}
	tcTcap.End.Components.ReturnResultLast = &ReturnResultTCAP{}

	tcTcap.End.Dtid = dtid
	tcTcap.End.Components.ReturnResultLast.InvokeID = invID
	tcTcap.End.Components.ReturnResultLast.OpCode = opCode
	tcTcap.End.Components.ReturnResultLast.Parameter = payload
	// Assign the dialogueObject to the Dialogue field. If dialogueObject is nil, no dialogue will be included.
	tcTcap.End.Dialogue = dialogueObject

	return tcTcap
}

// NewContinue create a Continue tcap message
// otid and dtid size from 1 to 4 bytes in BigEndian format
func NewContinue(otid []byte, dtid []byte) (*TCAP, error) {
	if err := validateTransactionID(otid, "otid"); err != nil {
		return nil, err
	}
	if err := validateTransactionID(dtid, "dtid"); err != nil {
		return nil, err
	}
	return NewContinueWithDialogue(otid, dtid, nil, nil), nil
}

// NewContinueWithDialogue create a Continue tcap message with a dialogue
// otid size from 1 to 4 bytes in BigEndian format
func NewContinueWithDialogue(otid []byte, dtid []byte, acn *int, acnVersion *int) *TCAP {
	tcTcap := &TCAP{}
	tcTcap.Continue = &ContinueTCAP{}

	tcTcap.Continue.Otid = otid
	tcTcap.Continue.Dtid = dtid

	if acn != nil && acnVersion != nil {
		tcTcap.Continue.Dialogue = &DialogueTCAP{}
		tcTcap.Continue.Dialogue.DialogueRequest = &AARQapduTCAP{}

		tcTcap.Continue.Dialogue.DialogAsId = []int{0, 0, 17, 773, 1, 1, 1}
		tcTcap.Continue.Dialogue.DialogueRequest.ProtocolVersionPadded = uint8Ptr(128) // protocol version 128 = 0x80
		tcTcap.Continue.Dialogue.DialogueRequest.AcnVersion = []int{0, 4, 0, 0, 1, 0, *acn, *acnVersion}
	}

	return tcTcap
}

// NewContinueWithDialogueObject creates a Continue TCAP message with a dialogue object.
// Parameters:
//   - otid: Originating Transaction ID, size from 1 to 4 bytes in BigEndian format.
//   - dtid: Destination Transaction ID, size from 1 to 4 bytes in BigEndian format.
//   - dialogueObject: A pointer to a DialogueTCAP object, representing the dialogue to include in the message.
//     If nil, no dialogue will be included in the message.
func NewContinueWithDialogueObject(otid []byte, dtid []byte, dialogueObject *DialogueTCAP) *TCAP {
	tcTcap := &TCAP{}
	tcTcap.Continue = &ContinueTCAP{}

	tcTcap.Continue.Otid = otid
	tcTcap.Continue.Dtid = dtid
	// Assign the dialogueObject to the Dialogue field. If dialogueObject is nil, no dialogue will be included.
	tcTcap.Continue.Dialogue = dialogueObject

	return tcTcap
}

// NewContinueInvoke create a Continue Invoke tcap message
// otid and dtid size from 1 to 4 bytes in BigEndian format
func NewContinueInvoke(otid []byte, dtid []byte, invID int, opCode uint8, payload []byte) (*TCAP, error) {
	if err := validateTransactionID(otid, "otid"); err != nil {
		return nil, err
	}
	if err := validateTransactionID(dtid, "dtid"); err != nil {
		return nil, err
	}
	if err := validateInvokeID(invID, "invID"); err != nil {
		return nil, err
	}
	return NewContinueInvokeWithDialogue(otid, dtid, invID, opCode, payload, nil, nil), nil
}

// NewContinueInvokeWithDialogue create a Continue Invoke tcap message with a dialogue
// otid and dtid size from 1 to 4 bytes in BigEndian format
func NewContinueInvokeWithDialogue(otid, dtid []byte, invID int, opCode uint8, payload []byte, acn *int, acnVersion *int) *TCAP {
	tcTcap := &TCAP{}
	tcTcap.Continue = &ContinueTCAP{}
	tcTcap.Continue.Components = &ComponentTCAP{}
	tcTcap.Continue.Components.Invoke = &InvokeTCAP{}

	tcTcap.Continue.Otid = otid
	tcTcap.Continue.Dtid = dtid
	tcTcap.Continue.Components.Invoke.InvokeID = invID
	tcTcap.Continue.Components.Invoke.OpCode = opCode
	tcTcap.Continue.Components.Invoke.Parameter = payload

	if acn != nil && acnVersion != nil {
		tcTcap.Continue.Dialogue = &DialogueTCAP{}
		tcTcap.Continue.Dialogue.DialogueRequest = &AARQapduTCAP{}

		tcTcap.Continue.Dialogue.DialogAsId = []int{0, 0, 17, 773, 1, 1, 1}
		tcTcap.Continue.Dialogue.DialogueRequest.ProtocolVersionPadded = uint8Ptr(128) // protocol version 128 = 0x80
		tcTcap.Continue.Dialogue.DialogueRequest.AcnVersion = []int{0, 4, 0, 0, 1, 0, *acn, *acnVersion}
	}

	return tcTcap
}

// NewContinueReturnResultLast creates a Continue ReturnResultLast TCAP message.
// otid and dtid are transaction IDs, each with a size of 1 to 4 bytes in BigEndian format.
func NewContinueReturnResultLast(otid, dtid []byte, invID int, opCode *uint8, payload []byte) *TCAP {
	return NewContinueReturnResultLastWithDialogue(otid, dtid, invID, opCode, payload, nil, nil)
}

// NewContinueReturnResultLastWithDialogue creates a Continue ReturnResultLast TCAP message with a dialogue.
// Parameters:
//   - otid: Originating Transaction ID, size from 1 to 4 bytes in BigEndian format.
//   - dtid: Destination Transaction ID, size from 1 to 4 bytes in BigEndian format.
func NewContinueReturnResultLastWithDialogue(otid, dtid []byte, invID int, opCode *uint8, payload []byte, acn *int, acnVersion *int) *TCAP {
	tcTcap := &TCAP{}
	tcTcap.Continue = &ContinueTCAP{}
	tcTcap.Continue.Components = &ComponentTCAP{}
	tcTcap.Continue.Components.ReturnResultLast = &ReturnResultTCAP{}

	tcTcap.Continue.Otid = otid
	tcTcap.Continue.Dtid = dtid
	tcTcap.Continue.Components.ReturnResultLast.InvokeID = invID
	tcTcap.Continue.Components.ReturnResultLast.OpCode = opCode
	tcTcap.Continue.Components.ReturnResultLast.Parameter = payload

	if acn != nil && acnVersion != nil {
		tcTcap.Continue.Dialogue = &DialogueTCAP{}
		tcTcap.Continue.Dialogue.DialogueRequest = &AARQapduTCAP{}

		tcTcap.Continue.Dialogue.DialogAsId = []int{0, 0, 17, 773, 1, 1, 1}
		tcTcap.Continue.Dialogue.DialogueRequest.ProtocolVersionPadded = uint8Ptr(128) // protocol version 128 = 0x80
		tcTcap.Continue.Dialogue.DialogueRequest.AcnVersion = []int{0, 4, 0, 0, 1, 0, *acn, *acnVersion}
	}

	return tcTcap
}

// NewContinueReturnResultLastWithDialogueObject creates a Continue ReturnResultLast TCAP message with a dialogue object.
// Parameters:
//   - otid: Originating Transaction ID, size from 1 to 4 bytes in BigEndian format.
//   - dtid: Destination Transaction ID, size from 1 to 4 bytes in BigEndian format.
//   - dialogueObject: A pointer to a DialogueTCAP object, representing the dialogue to include in the message.
//     If nil, no dialogue will be included in the message.
func NewContinueReturnResultLastWithDialogueObject(otid, dtid []byte, invID int, opCode *uint8, payload []byte, dialogueObject *DialogueTCAP) *TCAP {
	tcTcap := &TCAP{}
	tcTcap.Continue = &ContinueTCAP{}
	tcTcap.Continue.Components = &ComponentTCAP{}
	tcTcap.Continue.Components.ReturnResultLast = &ReturnResultTCAP{}

	tcTcap.Continue.Otid = otid
	tcTcap.Continue.Dtid = dtid
	tcTcap.Continue.Components.ReturnResultLast.InvokeID = invID
	tcTcap.Continue.Components.ReturnResultLast.OpCode = opCode
	tcTcap.Continue.Components.ReturnResultLast.Parameter = payload
	// Assign the dialogueObject to the Dialogue field. If dialogueObject is nil, no dialogue will be included.
	tcTcap.Continue.Dialogue = dialogueObject

	return tcTcap
}

// validateTransactionID validates that a transaction ID meets ITU-T Q.773 requirements
// Transaction ID must be 1 to 4 bytes in length
func validateTransactionID(tid []byte, fieldName string) error {
	if len(tid) < 1 || len(tid) > 4 {
		return fmt.Errorf("%s must be 1 to 4 bytes in length, got %d bytes", fieldName, len(tid))
	}
	return nil
}

// validateInvokeID validates that an invoke ID is within the valid range
// Invoke ID must be in range -128 to 127 (signed 8-bit integer)
func validateInvokeID(invID int, fieldName string) error {
	if invID < -128 || invID > 127 {
		return fmt.Errorf("%s must be in range -128 to 127, got %d", fieldName, invID)
	}
	return nil
}
