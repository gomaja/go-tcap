# go-tcap

[![Go CI](https://github.com/gomaja/go-tcap/actions/workflows/ci.yml/badge.svg)](https://github.com/gomaja/go-tcap/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/gomaja/go-map.svg)](https://pkg.go.dev/github.com/gomaja/go-map)
[![Go Report Card](https://goreportcard.com/badge/github.com/gomaja/go-map)](https://goreportcard.com/report/github.com/gomaja/go-map)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A robust TCAP (Transaction Capabilities Application Part) implementation in Go.

## Overview

Package tcap provides simple and painless handling of TCAP (Transaction Capabilities Application Part) in SS7/SIGTRAN protocol stack, implemented in the Go Programming Language.

TCAP is an ASN.1-based protocol used in telecommunications networks for SS7-SIGTRAN information exchange between applications. It's commonly used in mobile networks for operations like SMS delivery, subscriber information retrieval, and authentication.

The TCAP structures in this library are directly defined based on the ASN.1 definition from ITU-T Q.773 (06/97). This implementation does not use any ASN.1 files or ASN.1 parsers, making it lightweight and easy to integrate.

## Installation

```bash
go get github.com/gomaja/go-tcap
```

## Usage

### Parsing TCAP Messages

```go
package main

import (
    "encoding/hex"
    "fmt"

    "github.com/gomaja/go-tcap"
)

func main() {
    // Example TCAP message in hex format
    hexMsg := "62494804004734a86b1e281c060700118605010101a011600f80020780a1090607040000010014036c21a11f02010002012d3017800891328490507608f38101ff820891328490123009f1"

    // Decode hex to bytes
    tcapBytes, _ := hex.DecodeString(hexMsg)

    // Parse TCAP message
    tcapMsg, err := tcap.ParseDER(tcapBytes)
    if err != nil {
        fmt.Printf("Error parsing TCAP message: %v\n", err)
        return
    }

    // Access TCAP message fields
    if tcapMsg.Begin != nil {
        fmt.Printf("Begin message with OTID: %x\n", tcapMsg.Begin.Otid)

        if tcapMsg.Begin.Components != nil && tcapMsg.Begin.Components.Invoke != nil {
            fmt.Printf("Invoke with ID: %d, OpCode: %d\n", 
                tcapMsg.Begin.Components.Invoke.InvokeID,
                tcapMsg.Begin.Components.Invoke.OpCode)
        }
    }
}
```

### Creating TCAP Messages

```go
package main

import (
    "encoding/hex"
    "fmt"

    "github.com/gomaja/go-tcap"
)

func main() {
    // Create a Begin message with Invoke component
    otid := []byte{0x01, 0x02, 0x03, 0x04}
    invokeID := 1
    opCode := uint8(45) // Example operation code
    payload := []byte{0x01, 0x02, 0x03} // Example payload

    tcapMsg := tcap.NewBeginInvoke(otid, invokeID, opCode, payload)

    // Marshal to binary
    tcapBytes, err := tcapMsg.Marshal()
    if err != nil {
        fmt.Printf("Error marshaling TCAP message: %v\n", err)
        return
    }

    // Print as hex
    fmt.Printf("TCAP message: %x\n", tcapBytes)
}
```

## Features

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

## API Documentation

### Main Types

- `TCAP`: The main structure representing a TCAP message
- `BeginTCAP`, `EndTCAP`, `ContinueTCAP`, `AbortTCAP`, `UnidirectionalTCAP`: Different TCAP message types
- `ComponentTCAP`: Represents the component portion of a TCAP message
- `DialogueTCAP`: Represents the dialogue portion of a TCAP message

### Key Functions

- `ParseAny([]byte) (*TCAP, error)`: Parse any TCAP message (DER or non-DER encoded)
- `ParseDER([]byte) (*TCAP, error)`: Parse a DER-encoded TCAP message
- `(TCAP).Marshal() ([]byte, error)`: Marshal a TCAP message to binary
- `NewBegin([]byte) *TCAP`: Create a new Begin message
- `NewBeginInvoke([]byte, int, uint8, []byte) *TCAP`: Create a new Begin message with Invoke component
- `NewEndReturnResultLast([]byte, int, *uint8, []byte) *TCAP`: Create a new End message with ReturnResultLast component
- `NewContinueInvoke([]byte, []byte, int, uint8, []byte) *TCAP`: Create a new Continue message with Invoke component

## Common Use Cases

This library can be used to implement various SS7/SIGTRAN protocols that use TCAP, such as:

- MAP (Mobile Application Part)
- CAP (CAMEL Application Part)
- INAP (Intelligent Network Application Part)

Common operations include:
- SMS delivery and routing
- Subscriber information retrieval
- Authentication and authorization
- Call handling and routing


## Author

Marwan Jadid

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/gomaja/go-tcap/blob/main/LICENSE) file for details.
