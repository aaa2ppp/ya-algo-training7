package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
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
		w := string(input[l : r+1])
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

func run(in io.Reader, out io.Writer, pack func([]byte) []byte, unpack func([]byte) []byte) {
	br := bufio.NewReader(in)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	var cmd string
	if _, err := fmt.Fscanln(br, &cmd); err != nil {
		panic(err)
	}

	switch cmd {
	case "pack":
		text, err := br.ReadBytes('\n')
		if err != nil && err != io.EOF {
			panic(err)
		}
		text = bytes.TrimRight(text, " \t\r\n")

		data := pack(text)
		writeInt(bw, len(data), defaultWriteOpts())
		writeInts(bw, data, defaultWriteOpts())

	case "unpack":
		var n int
		if _, err := fmt.Fscanln(br, &n); err != nil {
			panic(err)
		}

		// scanInt заточен под чтение знаковых интов. Прямо в байты читать можно,
		// но если значения < 128 :(
		data := make([]int16, n)

		sc := bufio.NewScanner(br)
		sc.Split(bufio.ScanWords)
		if err := scanInts(sc, data); err != nil {
			panic(err)
		}

		// перекладываем в байты
		dataBytes := make([]byte, n)
		for i := range data {
			dataBytes[i] = byte(data[i])
		}

		text := unpack(dataBytes)
		bw.Write(text)
		bw.WriteByte('\n')

	default:
		panic("unknown command " + cmd)
	}
}

// ----------------------------------------------------------------------------

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	run(os.Stdin, os.Stdout, encode, decode)
}

// ----------------------------------------------------------------------------

func unsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func scanWord(sc *bufio.Scanner) (string, error) {
	if sc.Scan() {
		return sc.Text(), nil
	}
	return "", io.EOF
}

func scanInt(sc *bufio.Scanner) (int, error)                  { return scanIntX[int](sc) }
func scanTwoInt(sc *bufio.Scanner) (_, _ int, _ error)        { return scanTwoIntX[int](sc) }
func scanThreeInt(sc *bufio.Scanner) (_, _, _ int, _ error)   { return scanThreeIntX[int](sc) }
func scanFourInt(sc *bufio.Scanner) (_, _, _, _ int, _ error) { return scanFourIntX[int](sc) }

func scanIntX[T Int](sc *bufio.Scanner) (res T, err error) {
	sc.Scan()
	v, err := strconv.ParseInt(unsafeString(sc.Bytes()), 0, int(unsafe.Sizeof(res))<<3)
	return T(v), err
}

func scanTwoIntX[T Int](sc *bufio.Scanner) (v1, v2 T, err error) {
	v1, err = scanIntX[T](sc)
	if err == nil {
		v2, err = scanIntX[T](sc)
	}
	return v1, v2, err
}

func scanThreeIntX[T Int](sc *bufio.Scanner) (v1, v2, v3 T, err error) {
	v1, err = scanIntX[T](sc)
	if err == nil {
		v2, err = scanIntX[T](sc)
	}
	if err == nil {
		v3, err = scanIntX[T](sc)
	}
	return v1, v2, v3, err
}

func scanFourIntX[T Int](sc *bufio.Scanner) (v1, v2, v3, v4 T, err error) {
	v1, err = scanIntX[T](sc)
	if err == nil {
		v2, err = scanIntX[T](sc)
	}
	if err == nil {
		v3, err = scanIntX[T](sc)
	}
	if err == nil {
		v4, err = scanIntX[T](sc)
	}
	return v1, v2, v3, v4, err
}

func scanInts[T Int](sc *bufio.Scanner, a []T) error {
	for i := range a {
		v, err := scanIntX[T](sc)
		if err != nil {
			return err
		}
		a[i] = v
	}
	return nil
}

type Int interface {
	~int | ~int64 | ~int32 | ~int16 | ~int8 | byte
}

type Number interface {
	Int | ~float32 | ~float64
}

type writeOpts struct {
	sep   byte
	begin byte
	end   byte
}

func defaultWriteOpts() writeOpts {
	return writeOpts{sep: ' ', end: '\n'}
}

func writeInt[I Int](bw *bufio.Writer, v I, opts writeOpts) error {
	var buf [32]byte

	var err error
	if opts.begin != 0 {
		err = bw.WriteByte(opts.begin)
	}

	if err == nil {
		_, err = bw.Write(strconv.AppendInt(buf[:0], int64(v), 10))
	}

	if err == nil && opts.end != 0 {
		err = bw.WriteByte(opts.end)
	}

	return err
}

func writeInts[I Int](bw *bufio.Writer, a []I, opts writeOpts) error {
	var err error
	if opts.begin != 0 {
		err = bw.WriteByte(opts.begin)
	}

	if len(a) != 0 {
		var buf [32]byte

		if opts.sep == 0 {
			opts.sep = ' '
		}

		_, err = bw.Write(strconv.AppendInt(buf[:0], int64(a[0]), 10))

		for i := 1; err == nil && i < len(a); i++ {
			err = bw.WriteByte(opts.sep)
			if err == nil {
				_, err = bw.Write(strconv.AppendInt(buf[:0], int64(a[i]), 10))
			}
		}
	}

	if err == nil && opts.end != 0 {
		err = bw.WriteByte(opts.end)
	}

	return err
}
