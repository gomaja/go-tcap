package asn1tcapmodel

import "encoding/asn1"

// DialogueAll contains all structured dialogue structs
type DialogueAll struct {
	//https://oid-base.com/get/0.0.17.773.1.1.1
	//ASN.1 notation: {itu-t(0) recommendation(0) q(17) 773 as(1) dialogue-as(1) version1(1)}
	//dot notation: 0.0.17.773.1.1.1
	//use value: asn1.ObjectIdentifier([]int{0, 0, 17, 773, 1, 1, 1})
	DialogueAsId asn1.ObjectIdentifier `asn1:"optional"`

	DialoguePDU DialoguePDU `asn1:"tag:0"`
}

type DialoguePDU struct { // choice
	DialogueRequest  AARQapdu `asn1:"application,tag:0,optional"`
	DialogueResponse AAREapdu `asn1:"application,tag:1,optional"`
	DialogueAbort    ABRTapdu `asn1:"application,tag:4,optional"`
}

type AARQapdu struct {
	// Protocol Version: 0x03020780
	//[]byte{3, 2, 7, 128}
	//Use value: asn1.BitString{Bytes: []byte{128}, BitLength: 1}
	ProtocolVersionPadded asn1.BitString `asn1:"tag:0,optional"`

	// use value: asn1.ObjectIdentifier([]int{0, 4, 0, 0, 1, 0, ApplicationContextName, AcnVersion})
	ApplicationContextName asn1.ObjectIdentifier `asn1:"tag:1,explicit"` // or alternatively, the asn1.ObjectIdentifier can be wrapped in a struct and explicit can be removed
	UserInformation        UserInformation       `asn1:"tag:30,optional"`
}

type UserInformation struct {
	// this is an EXTERNAL (tag of EXTERNAL is 8), make sure to add this tag with constructor type when marshalling (final tag EXTERNAL + Constructor = 0x28)
	Data asn1.RawValue
}

type AAREapdu struct {
	// Protocol Version: 0x03020780
	//[]byte{3, 2, 7, 128}
	//Use value: asn1.BitString{Bytes: []byte{128}, BitLength: 1}
	ProtocolVersionPadded asn1.BitString `asn1:"tag:0,optional"`

	// use value: asn1.ObjectIdentifier([]int{0, 4, 0, 0, 1, 0, ApplicationContextName, AcnVersion})
	ApplicationContextName asn1.ObjectIdentifier     `asn1:"tag:1,explicit"`
	Result                 AssociateResult           `asn1:"tag:2"`
	ResultSourceDiagnostic AssociateSourceDiagnostic `asn1:"tag:3"`
	UserInformation        UserInformation           `asn1:"tag:30,optional"`
}

type ABRTapdu struct {
	// ABRTsourceInt
	AbortSource int `asn1:"tag:0"`

	UserInformation UserInformation `asn1:"tag:30,optional"`
}

type ABRTsourceInt int

const (
	DialogueServiceUser     ABRTsourceInt = 0
	DialogueServiceProvider ABRTsourceInt = 1
)

type AssociateResult struct {
	// AssociateResultInt
	Data int
}
type AssociateResultInt int

const (
	AssociateResultAccepted        AssociateResultInt = 0
	AssociateResultRejectPermanent AssociateResultInt = 1
)

type AssociateSourceDiagnostic struct { // choice
	// AssociateSourceDiagUserInt // default:255 will change omitting the optional field to value 255 instead of 0, the 255 (FieldOmissionValue) is used to check if the field arrived was empty
	DialogueServiceUser int `asn1:"default:255,tag:1,explicit,optional"`

	// AssociateSourceDiagProviderInt // default:255 will change omitting the optional field to value 255 instead of 0, the 255 (FieldOmissionValue) is used to check if the field arrived was empty
	DialogueServiceProvider int `asn1:"default:255,tag:2,explicit,optional"`
}

type AssociateSourceDiagUserInt int

const (
	ASDUserNull                           AssociateSourceDiagUserInt = 0
	ASDUserNoReasonGiven                  AssociateSourceDiagUserInt = 1
	ASDUserApplicationContextNotSupported AssociateSourceDiagUserInt = 2
)

type AssociateSourceDiagProviderInt int

const (
	ASDProviderNull                    AssociateSourceDiagProviderInt = 0
	ASDProviderNoReasonGiven           AssociateSourceDiagProviderInt = 1
	ASDProviderNoCommonDialoguePortion AssociateSourceDiagProviderInt = 2
)
