package tcap

import (
	"github.com/gomaja/go-tcap/asn1tcapmodel"
)

// EndOption represents a functional option for configuring End TCAP messages
type EndOption func(*EndTCAP) error

// NewEnd creates an End TCAP message using the options pattern
func NewEnd(dtid []byte, options ...EndOption) (TCAP, error) {
	if err := validateTransactionID(dtid, "dtid"); err != nil {
		return nil, err
	}

	tcEnd := &EndTCAP{
		Dtid: dtid,
	}

	for _, opt := range options {
		if err := opt(tcEnd); err != nil {
			return nil, err
		}
	}

	return tcEnd, nil
}

// WithEndDialogueResponse adds a dialogue to an End TCAP message
func WithEndDialogueResponse(acn, acnVersion int) EndOption {
	return func(end *EndTCAP) error {
		end.Dialogue = newDialogueResponse(acn, acnVersion)
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

func (tcEnd *EndTCAP) Marshal() ([]byte, error) {
	var asn1Tcap asn1tcapmodel.TCMessage

	// Convert based on which field is set in the TCAP struct
	if tcEnd != nil {
		asn1Tcap.End = convertEndTCAPToEnd(tcEnd)
	}

	return marshalAsn1TcapModel(asn1Tcap)
}

func (tcEnd *EndTCAP) MessageType() MessageType {
	return MessageTypeEnd
}
