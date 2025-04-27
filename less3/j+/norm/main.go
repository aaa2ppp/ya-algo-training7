package main

import (
	"bufio"
	"io"
	"os"
	"unicode"
)

func main() {
	br := bufio.NewReader(os.Stdin)
	bw := bufio.NewWriter(os.Stdout)
	defer bw.Flush()

	for {
		c, err := br.ReadByte()
		if err == io.EOF {
			break
		}
		c = byte(unicode.ToLower(rune(c)))
		if 'a' <= c && c <= 'z' {
			bw.WriteByte(c)
		}
	}
}
