package tcap

const (
	// ASN.1 field omission indicator
	// Used throughout the codebase to indicate optional fields should be omitted
	FieldOmissionValue = 255

	// Transaction ID constraints per ITU-T Q.773
	MinTransactionIDLength = 1
	MaxTransactionIDLength = 4

	// Invoke ID constraints (signed 8-bit integer range)
	MinInvokeID = -128
	MaxInvokeID = 127

	// ASN.1 tag values
	ExternalConstructorTag = 0x28 // EXTERNAL with constructor
	SequenceConstructorTag = 0x30 // SEQUENCE with constructor

	// Protocol version for TCAP dialogue
	DefaultProtocolVersion = 0x80 // 128 decimal
)

// Dialogue ASN.1 Object Identifier constants
var (
	// DefaultDialogueAsId represents the standard TCAP dialogue object identifier
	DefaultDialogueAsId = []int{0, 0, 17, 773, 1, 1, 1}

	// DefaultAcnPrefix represents the prefix for the Application Context Name (ACN) in TCAP dialogue identifiers.
	DefaultAcnPrefix = []int{0, 4, 0, 0, 1, 0}
)
