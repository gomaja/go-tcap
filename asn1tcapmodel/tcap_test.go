package asn1tcapmodel

import (
	"encoding/asn1"
	"encoding/hex"
	"testing"
)

func TestEnd(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "Test End Tcap Message (returnError to invoke sendRoutingInfoForSM)",
			input:       "303d4904004734a86b262824060700118605010101a0196117a109060704000001001403a203020100a305a1030201006c0da30b02010002010130030a0100",
			expectError: false,
		},
		{
			name:        "Test End Tcap Messsage (returnResultLast sendRoutingInfoForSM)",
			input:       "3055490402b0d1c66b2a2828060700118605010101a01d611b80020780a109060704000001001402a203020100a305a1030201006c21a21f020100301a02012d3015040806031128951337f4a009810791126316002012",
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tcapBytes, err := hex.DecodeString(tc.input)
			if err != nil {
				t.Errorf("failed to decode input: %v", err)
			}
			transaction := &End{}
			_, err = asn1.Unmarshal(tcapBytes, transaction)
			if (err != nil) != tc.expectError {
				t.Errorf("expected error: %v, got: %v", tc.expectError, err)
			}
		})
	}
}

func TestBegin(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "Test Begin Tcap Message (invoke sendRoutingInfoForSM)",
			input:       "30494804004734a86b1e281c060700118605010101a011600f80020780a1090607040000010014036c21a11f02010002012d3017800891328490507608f38101ff820891328490000005f7",
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tcapBytes, err := hex.DecodeString(tc.input)
			if err != nil {
				t.Errorf("failed to decode input: %v", err)
			}
			transaction := &Begin{}
			_, err = asn1.Unmarshal(tcapBytes, transaction)
			if (err != nil) != tc.expectError {
				t.Errorf("expected error: %v, got: %v", tc.expectError, err)
			}
		})
	}
}

func TestAbort(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "Abort dtid due to version mismatch (response to invoke sendRoutingInfoForSM)",
			input:       "3032490402b0d1c46b2a2828060700118605010101a01d611b80020780a109060704000001001402a203020101a305a103020102",
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tcapBytes, err := hex.DecodeString(tc.input)
			if err != nil {
				t.Errorf("failed to decode input: %v", err)
			}
			transaction := &Abort{}
			_, err = asn1.Unmarshal(tcapBytes, transaction)
			if (err != nil) != tc.expectError {
				t.Errorf("expected error: %v, got: %v", tc.expectError, err)
			}
		})
	}
}

func TestContinue(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "Continue due to long SMS (invoke mt-forwardSM (Short Message fragment 1 of 2))",
			input:       "3081d24804008bd04049040419000f6c81c3a181c002010102012c3081b7800826610011829761f6840891328490000005f704819e4009d047f6dbfe06000042217251400000a00500035f020190e53c0b947fd741e8b0bd0c9abfdb6510bcec26a7dd67d09c5e86cf41693728ffaecb41f2f2393da7cbc3f4f4db0d82cbdfe3f27cee0241d9e5f0bc0c32bfd9ecf71d44479741ecb47b0da2bf41e3771bce2ed3cb203abadc0685dd64d09c1e96d341e4323b6d2fcbd3ee33888e96bfeb6734e8c87edbdf2190bc3c96d7d3f476d94d77d5e70500",
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tcapBytes, err := hex.DecodeString(tc.input)
			if err != nil {
				t.Errorf("failed to decode input: %v", err)
			}
			transaction := &Continue{}
			_, err = asn1.Unmarshal(tcapBytes, transaction)
			if (err != nil) != tc.expectError {
				t.Errorf("expected error: %v, got: %v", tc.expectError, err)
			}
		})
	}
}
