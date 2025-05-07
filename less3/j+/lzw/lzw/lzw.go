package lzw

import (
	"bytes"
	"log"
	"os"
	"unsafe"
)

var _, debugEnable = os.LookupEnv("DEBUG")

func Encode(input []byte, maxBits int) []uint16 {
	if debugEnable {
		log.Println("== LZWEncode")
	}

	if len(input) == 0 {
		return nil
	}

	var output []uint16

	maxCode := (uint16(1) << maxBits) - 1
	lzw := make(map[string]uint16, maxCode-256+1)

	var (
		newCode  uint16 = 256
		prevCode uint16 = uint16(input[0])
	)

	if debugEnable {
		log.Printf("%c:", prevCode)
	}

	for l, r := 0, 1; r < len(input); r++ {
		word := unsafeString(input[l : r+1])

		if curCode, ok := lzw[word]; ok {
			if debugEnable {
				log.Printf("%c: found %s=%d", input[r], word, curCode)
			}
			prevCode = curCode
			continue
		}

		if newCode <= uint16(maxCode) {
			if debugEnable {
				log.Printf("%c: add %s->%d out %d", input[r], word, newCode, prevCode)
			}

			lzw[word] = newCode
			newCode++
		}

		output = append(output, prevCode)

		prevCode = uint16(input[r])
		l = r
	}

	if debugEnable {
		log.Printf("   out %d", prevCode)
	}
	output = append(output, prevCode)

	return output
}

func Decode(input []uint16) []byte {
	if debugEnable {
		log.Println("== LZWDecode")
	}

	if len(input) == 0 {
		return nil
	}

	var output bytes.Buffer

	type word struct {
		l, r int
	}
	lzw := make([]word, 0)

	output.WriteByte(byte(input[0]))
	if debugEnable {
		log.Printf("%d: out %c", input[0], input[0])
	}

	l := 0
	for _, code := range input[1:] {
		r := output.Len()

		if code < 256 {
			output.WriteByte(byte(code))
		} else {
			idx := int(code - 256)
			if idx == len(lzw) {
				output.Write(output.Bytes()[l:r])
				output.WriteByte(output.Bytes()[l])
			} else {
				it := lzw[idx]
				output.Write(output.Bytes()[it.l:it.r])
			}
		}

		if debugEnable {
			newCode := len(lzw) + 256
			log.Printf("%d: out %s add %d->%s", code, output.Bytes()[r:], newCode, output.Bytes()[l:r+1])
		}

		lzw = append(lzw, word{l, r + 1})
		l = r
	}

	return output.Bytes()
}

func unsafeString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
