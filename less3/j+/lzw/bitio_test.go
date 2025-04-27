package main

import (
	"strings"
	"testing"
)

func TestBitWriter_WriteRead(t *testing.T) {
	var sb strings.Builder
	w := NewBitWriter(&sb)
	w.WriteBits('A', 8)
	t.Log(sb.String())
	r := NewBitReader(strings.NewReader("B"))
	c, err := r.ReadBits(8)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%c", byte(c))
}
