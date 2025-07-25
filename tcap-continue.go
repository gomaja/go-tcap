package tcap

// ContinueOption represents a functional option for configuring Continue TCAP messages
type ContinueOption func(*ContinueTCAP) error

// NewContinue creates a Continue TCAP message using the options pattern
// Parameters:
//   - otid: Originating Transaction ID, size from 1 to 4 bytes in BigEndian format.
//   - dtid: Destination Transaction ID, size from 1 to 4 bytes in BigEndian format.
func NewContinue(otid []byte, dtid []byte, options ...ContinueOption) (*TCAP, error) {
	if err := validateTransactionID(otid, "otid"); err != nil {
		return nil, err
	}
	if err := validateTransactionID(dtid, "dtid"); err != nil {
		return nil, err
	}

	tcap := &TCAP{
		Continue: &ContinueTCAP{
			Otid: otid,
			Dtid: dtid,
		},
	}

	for _, opt := range options {
		if err := opt(tcap.Continue); err != nil {
			return nil, err
		}
	}

	return tcap, nil
}

// WithContinueDialogue adds a dialogue to a Continue TCAP message
func WithContinueDialogue(acn, acnVersion int) ContinueOption {
	return func(cont *ContinueTCAP) error {
		cont.Dialogue = &DialogueTCAP{}
		cont.Dialogue.DialogueRequest = &AARQapduTCAP{}
		cont.Dialogue.DialogAsId = DefaultDialogueAsId
		cont.Dialogue.DialogueRequest.ProtocolVersionPadded = uint8Ptr(DefaultProtocolVersion)
		cont.Dialogue.DialogueRequest.AcnVersion = []int{0, 4, 0, 0, 1, 0, acn, acnVersion}
		return nil
	}
}

// WithContinueDialogueObject adds a dialogue object to a Continue TCAP message
func WithContinueDialogueObject(dialogue *DialogueTCAP) ContinueOption {
	return func(cont *ContinueTCAP) error {
		cont.Dialogue = dialogue
		return nil
	}
}

// WithContinueInvoke adds an Invoke component to a Continue TCAP message
func WithContinueInvoke(invID int, opCode uint8, payload []byte) ContinueOption {
	return func(cont *ContinueTCAP) error {
		if err := validateInvokeID(invID, "invID"); err != nil {
			return err
		}
		if cont.Components == nil {
			cont.Components = &ComponentTCAP{}
		}
		cont.Components.Invoke = &InvokeTCAP{
			InvokeID:  invID,
			OpCode:    opCode,
			Parameter: payload,
		}
		return nil
	}
}

// WithContinueReturnResultLast adds a ReturnResultLast component to a Continue TCAP message
func WithContinueReturnResultLast(invID int, opCode *uint8, payload []byte) ContinueOption {
	return func(cont *ContinueTCAP) error {
		if err := validateInvokeID(invID, "invID"); err != nil {
			return err
		}
		if cont.Components == nil {
			cont.Components = &ComponentTCAP{}
		}
		cont.Components.ReturnResultLast = &ReturnResultTCAP{
			InvokeID:  invID,
			OpCode:    opCode,
			Parameter: payload,
		}
		return nil
	}
}
