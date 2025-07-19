package tcap

// EndOption represents a functional option for configuring End TCAP messages
type EndOption func(*EndTCAP) error

// NewEnd creates an End TCAP message using the options pattern
func NewEnd(dtid []byte, options ...EndOption) (*TCAP, error) {
	if err := validateTransactionID(dtid, "dtid"); err != nil {
		return nil, err
	}

	tcap := &TCAP{
		End: &EndTCAP{
			Dtid: dtid,
		},
	}

	for _, opt := range options {
		if err := opt(tcap.End); err != nil {
			return nil, err
		}
	}

	return tcap, nil
}

// WithEndDialogue adds a dialogue to an End TCAP message
func WithEndDialogue(acn, acnVersion int) EndOption {
	return func(end *EndTCAP) error {
		end.Dialogue = &DialogueTCAP{}
		end.Dialogue.DialogueRequest = &AARQapduTCAP{}
		end.Dialogue.DialogAsId = DefaultDialogueAsId
		end.Dialogue.DialogueRequest.ProtocolVersionPadded = uint8Ptr(DefaultProtocolVersion)
		end.Dialogue.DialogueRequest.AcnVersion = []int{0, 4, 0, 0, 1, 0, acn, acnVersion}
		return nil
	}
}

// WithEndDialogueObject adds a dialogue object to an End TCAP message
func WithEndDialogueObject(dialogue *DialogueTCAP) EndOption {
	return func(end *EndTCAP) error {
		end.Dialogue = dialogue
		return nil
	}
}

// WithEndReturnResultLast adds a ReturnResultLast component to an End TCAP message
func WithEndReturnResultLast(invID int, opCode *uint8, payload []byte) EndOption {
	return func(end *EndTCAP) error {
		if err := validateInvokeID(invID, "invID"); err != nil {
			return err
		}
		if end.Components == nil {
			end.Components = &ComponentTCAP{}
		}
		end.Components.ReturnResultLast = &ReturnResultTCAP{
			InvokeID:  invID,
			OpCode:    opCode,
			Parameter: payload,
		}
		return nil
	}
}
