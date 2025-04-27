package asn1tcapmodel

import "encoding/asn1"

// UniDialogueAll contains all unstructured dialogue structs
type UniDialogueAll struct {
	//https://oid-base.com/get/0.0.17.773.1.2.1
	//ASN.1 notation: {itu-t(0) recommendation(0) q(17) 773 as(1) unidialogue-as(2) version1(1)}
	//dot notation: 0.0.17.773.1.2.1
	//use value: asn1.ObjectIdentifier([]int{0, 0, 17, 773, 1, 2, 1})
	UniDialogueAsId asn1.ObjectIdentifier `asn1:"optional"`

	UniDialoguePDU UniDialoguePDU `asn1:"tag:0"` // TODO: verify tag
}

type UniDialoguePDU struct { // choice
	UniDialoguePDU AUDTApdu `asn1:"application,tag:0,optional"`
}

type AUDTApdu struct {
	// Protocol Version: 0x03020780
	//[]byte{3, 2, 7, 128}
	//Use value: asn1.BitString{Bytes: []byte{128}, BitLength: 1}
	ProtocolVersionPadded asn1.BitString `asn1:"tag:0,optional"`

	// use value: asn1.ObjectIdentifier([]int{0, 4, 0, 0, 1, 0, ApplicationContextName, AcnVersion})
	ApplicationContextName asn1.ObjectIdentifier `asn1:"tag:1,explicit"` // or alternatively, the asn1.ObjectIdentifier can be wrapped in a struct and explicit can be removed
	UserInformation        UserInformation       `asn1:"tag:30,optional"`
}
