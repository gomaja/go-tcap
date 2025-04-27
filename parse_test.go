package tcap

import (
	"encoding/hex"
	"testing"

	"github.com/gomaja/go-tcap/utils"
)

// This test should only contain DER encoded asn1 structs
func TestDer(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "1- Two components",
			input:       "646049040086e8976b262824060700118605010101a0196117a109060704000001001403a203020100a305a1030201006c30a220020100301b02012d3016040826611042173454f2a00a810891328490000005f2a10c02010102013f300403020240",
			expectError: false,
		},
		{
			name:        "2- SRI SM",
			input:       "62494804004734a86b1e281c060700118605010101a011600f80020780a1090607040000010014036c21a11f02010002012d3017800891328490507608f38101ff820891328490000005f7",
			expectError: false,
		},
		{
			name:        "3- error SRI SM",
			input:       "643d4904004734a86b262824060700118605010101a0196117a109060704000001001403a203020100a305a1030201006c0da30b02010002010130030a0100",
			expectError: false,
		},
		{
			name:        "4- abort DTID version mismatch",
			input:       "6732490402b0d1c46b2a2828060700118605010101a01d611b80020780a109060704000001001402a203020101a305a103020102",
			expectError: false,
		},
		{
			name:        "5- SRI SM response",
			input:       "6455490402b0d1c66b2a2828060700118605010101a01d611b80020780a109060704000001001402a203020100a305a1030201006c21a21f020100301a02012d3015040806031128951337f4a009810791126316002012",
			expectError: false,
		},
		{
			name:        "6- Begin otid long message",
			input:       "62264804008bd0406b1e281c060700118605010101a011600f80020780a109060704000001001903",
			expectError: false,
		},
		{
			name:        "7- Continue otid dtid long message",
			input:       "653448040419000f4904008bd0406b262824060700118605010101a0196117a109060704000001001903a203020100a305a103020100",
			expectError: false,
		},
		{
			name:        "8- invoke mt-forwardSM fragment 1 of 2",
			input:       "6581d24804008bd04049040419000f6c81c3a181c002010102012c3081b7800826610011829761f6840891328490000005f704819e4009d047f6dbfe06000042217251400000a00500035f020190e53c0b947fd741e8b0bd0c9abfdb6510bcec26a7dd67d09c5e86cf41693728ffaecb41f2f2393da7cbc3f4f4db0d82cbdfe3f27cee0241d9e5f0bc0c32bfd9ecf71d44479741ecb47b0da2bf41e3771bce2ed3cb203abadc0685dd64d09c1e96d341e4323b6d2fcbd3ee33888e96bfeb6734e8c87edbdf2190bc3c96d7d3f476d94d77d5e70500",
			expectError: false,
		},
		{
			name:        "9- returnResultLast",
			input:       "651348040419000f4904008bd0406c05a203020101",
			expectError: false,
		},
		{
			name:        "10- returnResultLast",
			input:       "640d4904008bd0406c05a203020102",
			expectError: false,
		},
		{
			name:        "11- invoke reportSM-DeliveryStatus",
			input:       "6247480403c940ec6b1e281c060700118605010101a011600f80020780a1090607040000010014036c1fa11d02010002012f30150407910201068163280407916427417901f00a0101",
			expectError: false,
		},
		{
			name:        "12- invoke mt-forwardSM",
			input:       "6281b8480403c93f576b1e281c060700118605010101a011600f80020780a1090607040000010019036c818fa1818c02010002012c308183800874020110261338f38407916427417901f0046e040bd0536152e85c0200004221824143220068c1f1f85d77d341582c360693c16c322c168bc5828865719a5e2683ee693a1ad4b44a4136180ce68281de6e900cf78ac95e321a0b449587dd7373592e2f9341d4372838bd06a9c82ca8e99a2689c8a00b34152641cd309b9cb697e7",
			expectError: false,
		},
		{
			name:        "13- abort dtid (for invoke mt-forwardSM)",
			input:       "672d490403c93f576b252823060700118605010101a0186416800100be11280f060704000001010101a004a4028000",
			expectError: false,
		},
		{
			name:        "14- invoke alertServiceCentre",
			input:       "6240480400d199b06b1a2818060700118605010101a00d600ba1090607040000010017026c1ca11a0201010201403012040891881088775859f70406915418730536",
			expectError: false,
		},
		{
			name:        "15- invoke alertServiceCentreWithoutResult",
			input:       "622448047c0801f86c1ca11a02010102013130120407917933192122f30407916427417960f1",
			expectError: false,
		},
		{
			name:        "16- invoke forwardSM",
			input:       "62818a48048c150d066c8181a17f02010002012e3077800832140080803138f684069169318488880463040b916971101174f40000422182612464805bd2e2b1252d467ff6de6c47efd96eb6a1d056cb0d69b49a10269c098537586e96931965b260d15613da72c29b91261bde72c6a1ad2623d682b5996d58331271375a0d1733eee4bd98ec768bd966b41c0d",
			expectError: false,
		},
		{
			name:        "17- invoke sendRoutingInfo",
			input:       "6259480403ed2d126b1a2818060700118605010101a00d600ba1090607040000010005036c35a1330201c5020116302b80049152828883010086079152629610103287050583370000aa0a0a0104040504038090a3ab04030205e0",
			expectError: false,
		},
		{
			name:        "18- invoke mt-forwardSM (Short Message fragment 1 of 2)",
			input:       "6581d24804008bd04049040419000f6c81c3a181c002010102012c3081b7800826610011829761f6840891328490000005f704819e4009d047f6dbfe06000042217251400000a00500035f020190e53c0b947fd741e8b0bd0c9abfdb6510bcec26a7dd67d09c5e86cf41693728ffaecb41f2f2393da7cbc3f4f4db0d82cbdfe3f27cee0241d9e5f0bc0c32bfd9ecf71d44479741ecb47b0da2bf41e3771bce2ed3cb203abadc0685dd64d09c1e96d341e4323b6d2fcbd3ee33888e96bfeb6734e8c87edbdf2190bc3c96d7d3f476d94d77d5e70500",
			expectError: false,
		},
		{
			name:        "19- invoke forwardSM (Short Message fragment 2 of 2)",
			input:       "655a4804008bd04049040419000f6c4ca14a02010202012c3042800826610011829761f6840891328490000005f7042c4409d047f6dbfe060000422172514000001d0500035f0202cae8ba5c9e2ecb5de377fb157ea9d1b0d93b1e06",
			expectError: false,
		},
		{
			name:        "20- returnResultLast, a response for an FSM",
			input:       "64354904000000016b262824060700118605010101a0196117a109060704000001001903a203020100a305a1030201006c05a203020100",
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tcapBytes, err := hex.DecodeString(tc.input)
			if err != nil {
				t.Errorf("failed to decode input: %v", err)
			}
			tc, err := ParseDER(tcapBytes)
			if err != nil {
				t.Errorf("failed to parse input: %v", err)
			}
			marshalled, err := tc.Marshal()
			if err != nil {
				t.Errorf("failed to marshal input: %v", err)
			}
			if string(marshalled) != string(tcapBytes) {
				t.Errorf("marshalled bytes don't match")
			}
		})
	}
}

// This test should only contain non-DER encoded asn1 structs (contains encoding with indefinite length)
func TestNonDERToDER(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "Tcap stream containing indefinite length (Non-DER). (returnError)",
			input:       "6443490400519a286b2a2828060700118605010101a01d611b80020780a109060704000001001903a203020100a305a1030201006c80a30b02010002010630030201010000",
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bytesNonDER, err := hex.DecodeString(tc.input)
			if err != nil {
				t.Errorf("Error decoding hex: %v", err)
			}
			derBytes, err := utils.MakeDER(bytesNonDER)
			if err != nil {
				t.Errorf("Error converting to DER: %v", err)
			}

			tc, err := ParseDER(derBytes)
			if err != nil {
				t.Errorf("failed to parse input: %v", err)
			}
			marshalled, err := tc.Marshal()
			if err != nil {
				t.Errorf("failed to marshal input: %v", err)
			}
			if string(marshalled) != string(derBytes) {
				t.Errorf("marshalled bytes don't match")
			}

		})
	}
}
