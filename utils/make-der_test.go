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
		{
			name:    "Camel-V2 invoke initialDP",
			input:   "6281a94804b70801a16b1e281c060700118605010101a011600f80020780a1090607040000010032016c80a17d020100020100307580010183070313890027821785010a8a088493975617699909bb0580038090a39c01029f320852507017322911f7bf34170201008107919756176999f9a309800752f099d05b37d0bf35038301119f3605f943d000039f3707919756176999f99f3807819830535304f99f390802420122806080020000",
			want:    "6281a74804b70801a16b1e281c060700118605010101a011600f80020780a1090607040000010032016c7fa17d020100020100307580010183070313890027821785010a8a088493975617699909bb0580038090a39c01029f320852507017322911f7bf34170201008107919756176999f9a309800752f099d05b37d0bf35038301119f3605f943d000039f3707919756176999f99f3807819830535304f99f39080242012280608002",
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
