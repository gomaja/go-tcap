# go-tcap

A robust ASN.1 TCAP (Transaction Capabilities Application Part) implementation from ITU-T Q.773 (06/97).

[![Go CI](https://github.com/gomaja/go-tcap/actions/workflows/ci.yml/badge.svg)](https://github.com/gomaja/go-tcap/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/gomaja/go-tcap.svg)](https://pkg.go.dev/github.com/gomaja/go-tcap)
[![Go Report Card](https://goreportcard.com/badge/github.com/gomaja/go-tcap)](https://goreportcard.com/report/github.com/gomaja/go-tcap)
[![Go Version](https://img.shields.io/github/go-mod/go-version/gomaja/go-tcap)](https://github.com/gomaja/go-tcap/blob/main/go.mod)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=gomaja_go-tcap&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=gomaja_go-tcap)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=gomaja_go-tcap&metric=coverage)](https://sonarcloud.io/summary/new_code?id=gomaja_go-tcap)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=gomaja_go-tcap&metric=bugs)](https://sonarcloud.io/summary/new_code?id=gomaja_go-tcap)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=gomaja_go-tcap&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=gomaja_go-tcap)

## Overview

Package tcap provides handling of TCAP (Transaction Capabilities Application Part) in SS7/SIGTRAN protocol stack.

TCAP is an ASN.1-based protocol used in telecommunications networks for SS7-SIGTRAN information exchange between applications. It's commonly used in mobile networks for operations like SMS delivery, subscriber information retrieval, and authentication.

The TCAP structures in this library are directly defined based on the ASN.1 definition from ITU-T Q.773 (06/97). This implementation does not use any ASN.1 files or ASN.1 parsers, making it easy to use and integrate.

## Installation

```bash
go get github.com/gomaja/go-tcap
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
