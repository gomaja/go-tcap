package tcap

import (
	"strings"
	"testing"
)

func TestValidateTransactionID(t *testing.T) {
	tests := []struct {
		name        string
		tid         []byte
		fieldName   string
		expectError bool
		errorText   string
	}{
		{
			name:        "Valid 1 byte",
			tid:         []byte{0x01},
			fieldName:   "otid",
			expectError: false,
		},
		{
			name:        "Valid 2 bytes",
			tid:         []byte{0x01, 0x02},
			fieldName:   "otid",
			expectError: false,
		},
		{
			name:        "Valid 3 bytes",
			tid:         []byte{0x01, 0x02, 0x03},
			fieldName:   "otid",
			expectError: false,
		},
		{
			name:        "Valid 4 bytes",
			tid:         []byte{0x01, 0x02, 0x03, 0x04},
			fieldName:   "otid",
			expectError: false,
		},
		{
			name:        "Invalid empty",
			tid:         []byte{},
			fieldName:   "otid",
			expectError: true,
			errorText:   "must be 1 to 4 bytes in length, got 0 bytes",
		},
		{
			name:        "Invalid 5 bytes",
			tid:         []byte{0x01, 0x02, 0x03, 0x04, 0x05},
			fieldName:   "dtid",
			expectError: true,
			errorText:   "must be 1 to 4 bytes in length, got 5 bytes",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateTransactionID(tc.tid, tc.fieldName)
			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if !strings.Contains(err.Error(), tc.errorText) {
					t.Errorf("expected error containing '%s', got '%s'", tc.errorText, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateInvokeID(t *testing.T) {
	tests := []struct {
		name        string
		invID       int
		fieldName   string
		expectError bool
		errorText   string
	}{
		{
			name:        "Valid minimum",
			invID:       -128,
			fieldName:   "invID",
			expectError: false,
		},
		{
			name:        "Valid maximum",
			invID:       127,
			fieldName:   "invID",
			expectError: false,
		},
		{
			name:        "Valid zero",
			invID:       0,
			fieldName:   "invID",
			expectError: false,
		},
		{
			name:        "Valid positive",
			invID:       42,
			fieldName:   "invID",
			expectError: false,
		},
		{
			name:        "Valid negative",
			invID:       -42,
			fieldName:   "invID",
			expectError: false,
		},
		{
			name:        "Invalid too low",
			invID:       -129,
			fieldName:   "invID",
			expectError: true,
			errorText:   "must be in range -128 to 127, got -129",
		},
		{
			name:        "Invalid too high",
			invID:       128,
			fieldName:   "linkedID",
			expectError: true,
			errorText:   "must be in range -128 to 127, got 128",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateInvokeID(tc.invID, tc.fieldName)
			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if !strings.Contains(err.Error(), tc.errorText) {
					t.Errorf("expected error containing '%s', got '%s'", tc.errorText, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestNewBegin(t *testing.T) {
	tests := []struct {
		name        string
		otid        []byte
		expectError bool
	}{
		{
			name:        "Valid otid",
			otid:        []byte{0x01, 0x02},
			expectError: false,
		},
		{
			name:        "Invalid empty otid",
			otid:        []byte{},
			expectError: true,
		},
		{
			name:        "Invalid long otid",
			otid:        []byte{0x01, 0x02, 0x03, 0x04, 0x05},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tcap, err := NewBegin(tc.otid)
			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				if tcap != nil {
					t.Errorf("expected nil TCAP on error")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if tcap == nil {
					t.Errorf("expected valid TCAP")
					return
				}
				if tcap.Begin == nil {
					t.Errorf("expected Begin field to be set")
				}
			}
		})
	}
}

func TestNewBeginInvoke(t *testing.T) {
	tests := []struct {
		name        string
		otid        []byte
		invID       int
		opCode      uint8
		payload     []byte
		expectError bool
		errorText   string
	}{
		{
			name:        "Valid parameters",
			otid:        []byte{0x01, 0x02},
			invID:       42,
			opCode:      0x2D,
			payload:     []byte{0x01, 0x02, 0x03},
			expectError: false,
		},
		{
			name:        "Invalid otid",
			otid:        []byte{},
			invID:       42,
			opCode:      0x2D,
			payload:     []byte{0x01, 0x02, 0x03},
			expectError: true,
			errorText:   "must be 1 to 4 bytes in length",
		},
		{
			name:        "Invalid invID",
			otid:        []byte{0x01, 0x02},
			invID:       128,
			opCode:      0x2D,
			payload:     []byte{0x01, 0x02, 0x03},
			expectError: true,
			errorText:   "must be in range -128 to 127",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tcap, err := NewBegin(tc.otid, WithBeginInvoke(tc.invID, tc.opCode, tc.payload))
			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if tc.errorText != "" && !strings.Contains(err.Error(), tc.errorText) {
					t.Errorf("expected error containing '%s', got '%s'", tc.errorText, err.Error())
				}
				if tcap != nil {
					t.Errorf("expected nil TCAP on error")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if tcap == nil {
					t.Errorf("expected valid TCAP")
					return
				}
				if tcap.Begin == nil || tcap.Begin.Components == nil || tcap.Begin.Components.Invoke == nil {
					t.Errorf("expected Begin Invoke to be properly set")
				}
			}
		})
	}
}

func TestNewContinue(t *testing.T) {
	tests := []struct {
		name        string
		otid        []byte
		dtid        []byte
		expectError bool
		errorText   string
	}{
		{
			name:        "Valid parameters",
			otid:        []byte{0x01, 0x02},
			dtid:        []byte{0x03, 0x04},
			expectError: false,
		},
		{
			name:        "Invalid otid",
			otid:        []byte{},
			dtid:        []byte{0x03, 0x04},
			expectError: true,
			errorText:   "must be 1 to 4 bytes in length",
		},
		{
			name:        "Invalid dtid",
			otid:        []byte{0x01, 0x02},
			dtid:        []byte{0x01, 0x02, 0x03, 0x04, 0x05},
			expectError: true,
			errorText:   "must be 1 to 4 bytes in length",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tcap, err := NewContinue(tc.otid, tc.dtid)
			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if tc.errorText != "" && !strings.Contains(err.Error(), tc.errorText) {
					t.Errorf("expected error containing '%s', got '%s'", tc.errorText, err.Error())
				}
				if tcap != nil {
					t.Errorf("expected nil TCAP on error")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if tcap == nil {
					t.Errorf("expected valid TCAP")
					return
				}
				if tcap.Continue == nil {
					t.Errorf("expected Continue field to be set")
				}
			}
		})
	}
}
