package utils

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestMakeDER(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "Boolean Type included (invoke sendRoutingInfoForSM)",
			input:   "3019800a915282051447720982f9810101820891328490001015f8",
			want:    "3019800a915282051447720982f98101ff820891328490001015f8",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Decode hex string to bytes
			originalBytes, err := hex.DecodeString(tt.input)
			if err != nil {
				t.Fatalf("Failed to decode hex string: %v", err)
			}

			got, err := MakeDER(originalBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeDER() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			wantBytes, err := hex.DecodeString(tt.want)
			if err != nil {
				t.Fatalf("Failed to decode expected hex string: %v", err)
			}

			if !bytes.Equal(got, wantBytes) {
				t.Errorf("MakeDER() = %x, want %x", got, wantBytes)
			}

		})
	}
}
