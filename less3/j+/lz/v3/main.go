package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"unsafe"
)

const maxCount = (1 << 14) - 1

func encode(input []byte) []byte {
	index := map[string]int{}

	var output []byte
	var outByte byte
	var bitCnt int

	writeBit := func(v byte) {
		outByte |= (v & 1) << (bitCnt & 7)
		bitCnt++
		if bitCnt&7 == 0 {
			output = append(output, outByte)
			outByte = 0
		}
	}

	writeBits := func(v int, n int) {
		for i := 0; i < n; i++ {
			writeBit(byte(v))
			v >>= 1
		}
	}

	l, r, p := 0, 0, 0

	prev, count := 0, 0
	for r < len(input) {
		w := unsafeString(input[l : r+1])
		cur, ok := index[w]
		if ok {
			prev = cur
			r++
			continue
		}

		if count >= (1 << p) {
			p++
		}
		// log.Printf("%d (%d) '%c'", prev, p, input[r])
		// writeBits(prev, p)
		writeBits(int(input[r]), 8)
		// log.Println("bitCnt:", bitCnt)

		prev = 0
		r++
		l = r

		if count == maxCount {
			continue
		}
		count++
		index[w] = count
		// log.Println(w, "->", count)
	}

	log.Println("words count:", count, "p:", p)

	if prev != 0 {
		// log.Printf("%d (%d)", prev, p)
		// writeBits(prev, p)
		// log.Println("bitCnt:", bitCnt)
	}

	if bitCnt&7 != 0 {
		output = append(output, outByte)
	}

	return output // todo
}

func decode(input []byte) []byte {
	var bitCnt int
	readBit := func() int {
		if bitCnt>>3 >= len(input) {
			return -1
		}
		x := (int(input[bitCnt>>3]) >> (bitCnt & 7)) & 1
		bitCnt++
		return x
	}

	readBits := func(n int) int {
		var x int
		for i := 0; i < n; i++ {
			x |= readBit() << i
		}
		return x
	}

	dict := map[int]string{}
	dict[0] = ""

	var output bytes.Buffer
	p := 0
	count := 0
	for {
		if count >= (1 << p) {
			p++
		}

		prev := readBits(p)
		if prev < 0 {
			break
		}
		w := dict[prev]
		output.WriteString(w)

		c := readBits(5)
		if c <= 0 {
			break
		}
		c += 'a' - 1
		output.WriteByte(byte(c))

		if count == maxCount {
			continue
		}

		count++
		w += string(byte(c))
		// log.Println(w)
		dict[count] = w
	}

	return output.Bytes()
}

func main() {

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 && os.Args[1] == "encode" {
		os.Stdout.Write(encode(input))
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "decode" {
		os.Stdout.Write(decode(input))
		return
	}

	log.Println("input len:", len(input))

	encoded := encode(input)
	log.Printf("encoded len: %v %0.2f%%", len(encoded),
		float64(len(encoded))/float64(len(input))*100)

	decoded := decode(encoded)
	log.Println("decoded len:", len(decoded))

	if bytes.Equal(input, decoded) {
		log.Println("ok")
	} else {
		log.Println("decoded != input")
	}

	// if _, err := os.Stdout.Write(decoded); err != nil {
	// 	log.Fatal(err)
	// }
}

func unsafeString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
