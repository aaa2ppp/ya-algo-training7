package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"unsafe"
)

const maxCount = (1 << 14)

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

	// резервируем три бита под кол-во значащих битов в последнем байте
	writeBits(0, 3)

	l, r, p := 0, 0, 0
	prev, count := 0, 1
	for r < len(input) {
		w := unsafeString(input[l : r+1])
		cur, ok := index[w]
		if ok {
			prev = cur
			r++
			continue
		}

		if count > (1 << p) {
			p++
		}

		writeBits(prev, p)

		for _, bit := range huffmanDict[input[r]] {
			writeBit(bit)
		}

		prev = 0
		r++
		l = r

		if count < maxCount {
			index[w] = count
			count++
		}
	}

	if prev != 0 {
		writeBits(prev, p)
	}

	if bitCnt&7 != 0 {
		output = append(output, outByte)
	}

	// запоминаем кол-во значащих битов в последнем байте.
	// если 0, то все биты значащие
	output[0] |= byte(bitCnt & 7)

	return output
}

func decode(input []byte) []byte {
	totalBits := len(input) * 8
	if input[0]&7 != 0 {
		totalBits -= 8 - int(input[0]&7)
	}

	var bitCnt int

	readBit := func() int {
		if bitCnt == totalBits {
			return -1
		}
		v := (int(input[bitCnt>>3]) >> (bitCnt & 7)) & 1
		bitCnt++
		return v
	}

	readBits := func(n int) int {
		var v int
		for i := 0; i < n; i++ {
			v |= readBit() << i
		}
		return v
	}

	readBits(3)

	dict := map[int]string{}
	dict[0] = ""
	tree := newHuffmanTree(huffmanDict)

	var output bytes.Buffer
	p := 0
	count := 1
mainLoop:
	for {
		if count > (1 << p) {
			p++
		}

		prev := readBits(p)
		if prev < 0 {
			break
		}
		w := dict[prev]
		output.WriteString(w)

		// c := readBits(5)
		// if c <= 0 {
		// 	break
		// }
		// c += 'a' - 1

		node := tree
		for !node.isLeaf() {
			bit := readBit()
			if bit < 0 {
				break mainLoop
			}
			switch bit {
			case 0:
				node = node.Left
			case 1:
				node = node.Right
			}
		}
		c := node.Let

		output.WriteByte(byte(c))

		if count < maxCount {
			dict[count] = w + string(byte(c))
			count++
		}
	}

	return output.Bytes()
}

// XXX заточено под "тексты на английском языке, записанные только маленькими
// английскими буквами, без пробелов, знаков препинания и вообще каких-либо символов,
// отличных от маленьких английских букв. Эти тексты являются художественными
// произведениями (естественным текстом)."
var huffmanDict = [][]byte{
	'e': {0, 0, 1},
	'l': {0, 0, 0, 0},
	'n': {0, 1, 0, 0},
	'r': {0, 1, 0, 1},
	'h': {1, 0, 0, 0},
	's': {1, 0, 0, 1},
	'i': {1, 0, 1, 1},
	'o': {1, 1, 0, 0},
	't': {1, 1, 0, 1},
	'a': {1, 1, 1, 1},
	'y': {0, 0, 0, 1, 0},
	'b': {0, 0, 0, 1, 1},
	'c': {0, 1, 1, 0, 0},
	'f': {0, 1, 1, 0, 1},
	'u': {0, 1, 1, 1, 0},
	'm': {1, 0, 1, 0, 0},
	'w': {1, 0, 1, 0, 1},
	'd': {1, 1, 1, 0, 0},
	'v': {0, 1, 1, 1, 1, 0},
	'p': {1, 1, 1, 0, 1, 0},
	'g': {1, 1, 1, 0, 1, 1},
	'k': {0, 1, 1, 1, 1, 1, 1},
	'x': {0, 1, 1, 1, 1, 1, 0, 0, 0},
	'z': {0, 1, 1, 1, 1, 1, 0, 0, 1},
	'q': {0, 1, 1, 1, 1, 1, 0, 1, 0},
	'j': {0, 1, 1, 1, 1, 1, 0, 1, 1},
}

type huffmanNode struct {
	Left  *huffmanNode
	Right *huffmanNode
	Let   byte
}

func (node *huffmanNode) isLeaf() bool {
	// return node.Left == nil && node.Right == nil
	return node.Let != 0
}

func newHuffmanTree(dict [][]byte) *huffmanNode {
	root := &huffmanNode{}
	for c := 'a'; c <= 'z'; c++ {
		node := root
		for _, dir := range dict[c] {
			switch dir {
			case 0:
				if node.Left == nil {
					node.Left = &huffmanNode{}
				}
				node = node.Left
			case 1:
				if node.Right == nil {
					node.Right = &huffmanNode{}
				}
				node = node.Right
			}
		}
		node.Let = byte(c)
	}
	return root
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
}

func unsafeString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
