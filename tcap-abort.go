package tcap

import (
	"github.com/gomaja/go-tcap/asn1tcapmodel"
)

func (tcAbort *AbortTCAP) Marshal() ([]byte, error) {
	var asn1Tcap asn1tcapmodel.TCMessage

	// Convert based on which field is set in the TCAP struct
	if tcAbort != nil {
		asn1Tcap.Abort = convertAbortTCAPToAbort(tcAbort)
	}

	return marshalAsn1TcapModel(asn1Tcap)
}

func (tcAbort *AbortTCAP) MessageType() MessageType {
	return MessageTypeAbort
}
