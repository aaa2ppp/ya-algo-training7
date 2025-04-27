package main

import (
	"bytes"
	"errors"
	"log"
	"unsafe"
)

func LZWEncode(input []byte, maxBits int) []uint16 {
	maxSize := uint16(1) << maxBits

	var output []uint16
	lzw := make(map[string]uint16, 0)

	var count uint16 = 256
	var prev uint16

	for l, r := 0, 0; r < len(input); r++ {
		if l == r {
			prev = uint16(input[r])
			if debugEnable {
				log.Printf("%c --- --- %s", input[r], string(input[r]))
			}
			continue
		}

		w := unsafeString(input[l : r+1])
		if cur, ok := lzw[w]; ok {
			if debugEnable {
				log.Printf("%c --- --- %s", input[r], w)
			}
			prev = cur
			continue
		}

		if debugEnable {
			log.Printf("%c %3d %3d %s->%d", input[r], prev, uint16(input[r]), w, count)
		}
		output = append(output, prev, uint16(input[r]))
		if count < uint16(maxSize) {
			lzw[w] = count
			count++
		}
		prev = 0
		l = r + 1
	}

	if prev != 0 {
		if debugEnable {
			log.Printf("- %3d", prev)
		}
		output = append(output, prev)
	}

	return output
}

func LZWDecode(input []uint16) ([]byte, error) {
	var output bytes.Buffer
	lzw := make([][2]int, 0)

	for i := 0; i < len(input); i += 2 {
		l := output.Len()

		if prev := input[i]; prev < 256 {
			output.WriteByte(byte(prev))
		} else {
			prev -= 256
			w := unsafeString(output.Bytes()[lzw[prev][0]:lzw[prev][1]])
			output.WriteString(w)
		}

		if i+1 >= len(input) {
			break
		}

		if cur := input[i+1]; cur < 256 {
			output.WriteByte(byte(cur))
		} else {
			return nil, errors.New("bad sequence")
		}

		lzw = append(lzw, [2]int{l, output.Len()})
	}

	return output.Bytes(), nil
}

func unsafeString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
