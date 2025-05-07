package lzw

import (
	"bytes"
	"testing"
)

func Test_EncodeDecode(t *testing.T) {
	tests := []struct {
		name  string
		text  []byte
		debug bool
	}{
		{
			"1",
			[]byte("ABABCBABABCAD"),
			true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func(v bool) { debugEnable = v }(debugEnable)
			debugEnable = tt.debug

			encoded := Encode(tt.text, 16)
			decoded := Decode(encoded)
			if !bytes.Equal(tt.text, decoded) {
				t.Errorf("\ndecode() = %s, \nwant %s", decoded, tt.text)
			}
		})
	}
}
