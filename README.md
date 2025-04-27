# go-tcap

[![Go CI](https://github.com/gomaja/go-tcap/actions/workflows/ci.yml/badge.svg)](https://github.com/gomaja/go-tcap/actions/workflows/ci.yml)

Robust TCAP implementation in Golang

Package tcap provides simple and painless handling of TCAP(Transaction Capabilities Application Part) in SS7/SIGTRAN protocol stack, implemented in the Go Programming Language.

TCAP is an ASN.1-based protocol. The TCAP structures are directly defined based on the ASN.1 definition. ITU-T Q.773 (06/97)
This implementation does not use any ASN.1 files or ASN.1 parsers.

## Supported Features

### Transaction Portion

#### Message Types

| Message type   | Supported? |
|----------------|------------|
| Unidirectional | Yes        |
| Begin          | Yes        |
| End            | Yes        |
| Continue       | Yes        |
| Abort          | Yes        |

#### Fields

| Tag                        | Supported? |
|----------------------------|------------|
| Originating Transaction ID | Yes        |
| Destination Transaction ID | Yes        |
| P-Abort Cause              | Yes        |

### Component Portion

#### Component types

| Component type           | Supported? |
|--------------------------|------------|
| Invoke                   | Yes        |
| Return Result (Last)     | Yes        |
| Return Result (Not Last) | Yes        |
| Return Error             | Yes        |
| Reject                   | Yes        |


### Dialogue Portion

#### Dialogue types

| Dialogue type                       | Supported? |
|-------------------------------------|------------|
| Dialogue Request (AARQ-apdu)        | Yes        |
| Dialogue Response (AARE-apdu)       | Yes        |
| Dialogue Abort (ABRT-apdu)          | Yes        |
| Unidirectional Dialogue (AUDT-apdu) | Yes        |

#### Elements

| Tag                         | Type         | Supported? |
|-----------------------------|--------------|------------|
| Object Identifier           | Structured   | Yes        |
| Dialogue PDU                | Structured   | Yes        |
| Object Identifier           | Unstructured | Yes        |
| Unidirectional Dialogue PDU | Unstructured | Yes        |


## Author

Marwan Jadid

## License

[MIT](https://github.com/gomaja/go-tcap/blob/main/LICENSE)
