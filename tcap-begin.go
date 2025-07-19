package tcap

// BeginOption represents a functional option for configuring Begin TCAP messages
type BeginOption func(*BeginTCAP) error

// NewBegin creates a Begin TCAP message using the options pattern
// Parameters:
//   - otid: Originating Transaction ID, size from 1 to 4 bytes in BigEndian format.
func NewBegin(otid []byte, options ...BeginOption) (*TCAP, error) {
	if err := validateTransactionID(otid, "otid"); err != nil {
		return nil, err
	}

	tcap := &TCAP{
		Begin: &BeginTCAP{
			Otid: otid,
		},
	}

	for _, opt := range options {
		if err := opt(tcap.Begin); err != nil {
			return nil, err
		}
	}

	return tcap, nil
}

// WithBeginDialogue adds a dialogue to a Begin TCAP message
func WithBeginDialogue(acn, acnVersion int) BeginOption {
	return func(begin *BeginTCAP) error {
		begin.Dialogue = &DialogueTCAP{}
		begin.Dialogue.DialogueRequest = &AARQapduTCAP{}
		begin.Dialogue.DialogAsId = DefaultDialogueAsId
		begin.Dialogue.DialogueRequest.ProtocolVersionPadded = uint8Ptr(DefaultProtocolVersion)
		begin.Dialogue.DialogueRequest.AcnVersion = []int{0, 4, 0, 0, 1, 0, acn, acnVersion}
		return nil
	}
}

// WithBeginDialogueObject adds a dialogue object to a Begin TCAP message
func WithBeginDialogueObject(dialogue *DialogueTCAP) BeginOption {
	return func(begin *BeginTCAP) error {
		begin.Dialogue = dialogue
		return nil
	}
}

// WithBeginInvoke adds an Invoke component to a Begin TCAP message
func WithBeginInvoke(invID int, opCode uint8, payload []byte) BeginOption {
	return func(begin *BeginTCAP) error {
		if err := validateInvokeID(invID, "invID"); err != nil {
			return err
		}
		if begin.Components == nil {
			begin.Components = &ComponentTCAP{}
		}
		begin.Components.Invoke = &InvokeTCAP{
			InvokeID:  invID,
			OpCode:    opCode,
			Parameter: payload,
		}
		return nil
	}
}
