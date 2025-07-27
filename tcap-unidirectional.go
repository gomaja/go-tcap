package tcap

import (
	"github.com/gomaja/go-tcap/asn1tcapmodel"
)

func (tcUnidirectional *UnidirectionalTCAP) Marshal() ([]byte, error) {
	var asn1Tcap asn1tcapmodel.TCMessage

	// Convert based on which field is set in the TCAP struct
	if tcUnidirectional != nil {
		asn1Tcap.Unidirectional = convertUnidirectionalTCAPToUnidirectional(tcUnidirectional)
	}

	return marshalAsn1TcapModel(asn1Tcap)
}

func (tcUnidirectional *UnidirectionalTCAP) MessageType() MessageType {
	return MessageTypeUnidirectional
}
