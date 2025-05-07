package bitio

import (
	"strings"
	"testing"
)

func TestWriter_WriteRead(t *testing.T) {
	var sb strings.Builder
	w := NewWriter(&sb)
	w.WriteBits('A', 7)
	w.Close()
	if got := sb.String(); got != "A" {
		t.Errorf("got = %s, want A", got)
	}

	r := NewReader(strings.NewReader("B"))
	got, err := r.ReadBits(7)
	if err != nil {
		t.Fatal(err)
	}
	if got != 'B' {
		t.Errorf("got = %c, want B", got)
	}
}
