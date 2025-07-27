package tcap

import (
	"github.com/gomaja/go-tcap/asn1tcapmodel"
)

// BeginOption represents a functional option for configuring Begin TCAP messages
type BeginOption func(*BeginTCAP) error

// NewBegin creates a Begin TCAP message using the options pattern
// Parameters:
//   - otid: Originating Transaction ID, size from 1 to 4 bytes in BigEndian format.
func NewBegin(otid []byte, options ...BeginOption) (TCAP, error) {
	if err := validateTransactionID(otid, "otid"); err != nil {
		return nil, err
	}

	tcBegin := &BeginTCAP{
		Otid: otid,
	}

	for _, opt := range options {
		if err := opt(tcBegin); err != nil {
			return nil, err
		}
	}

	return tcBegin, nil
}

// WithBeginDialogueRequest adds a dialogue to a Begin TCAP message
func WithBeginDialogueRequest(acn, acnVersion int) BeginOption {
	return func(begin *BeginTCAP) error {
		begin.Dialogue = newDialogueRequest(acn, acnVersion)
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

func (tcBegin *BeginTCAP) Marshal() ([]byte, error) {
	var asn1Tcap asn1tcapmodel.TCMessage

	// Convert based on which field is set in the TCAP struct
	if tcBegin != nil {
		asn1Tcap.Begin = convertBeginTCAPToBegin(tcBegin)
	}

	return marshalAsn1TcapModel(asn1Tcap)
}

func (tcBegin *BeginTCAP) MessageType() MessageType {
	return MessageTypeBegin
}
