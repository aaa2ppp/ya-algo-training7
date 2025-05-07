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
		writeInt(bw, len(data))
		writeInts(bw, data)

	case "unpack":
		var n int
		if _, err := fmt.Fscanln(br, &n); err != nil {
			panic(err)
		}

		sc := bufio.NewScanner(br)
		sc.Split(bufio.ScanWords)

		data, err := scanInts(sc, make([]byte, n))
		if err != nil {
			panic(err)
		}

		text := unpack(data)
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

type Sign interface {
	~int | ~int64 | ~int32 | ~int16 | ~int8
}

type Unsign interface {
	~uint | ~uint64 | ~uint32 | ~uint16 | ~uint8
}

type Int interface {
	Sign | Unsign
}

type Float interface {
	~float32 | ~float64
}

type Number interface {
	Int | Float
}

// ----------------------------------------------------------------------------

func unsafeString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func scanWord(sc *bufio.Scanner) (string, error) {
	if sc.Scan() {
		return sc.Text(), nil
	}
	if err := sc.Err(); err != nil {
		return "", err
	}
	return "", io.EOF
}

func _parseInt[X Int](b []byte) (X, error) {
	if ^X(0) < 0 {
		v, err := strconv.ParseInt(unsafeString(b), 0, int(unsafe.Sizeof(X(1)))<<3)
		return X(v), err
	} else {
		v, err := strconv.ParseUint(unsafeString(b), 0, int(unsafe.Sizeof(X(1)))<<3)
		return X(v), err
	}
}

func scanIntX[X Int](sc *bufio.Scanner) (X, error) {
	if sc.Scan() {
		return _parseInt[X](sc.Bytes())
	}
	if err := sc.Err(); err != nil {
		return 0, err
	}
	return 0, io.EOF
}

func scanInts[X Int](sc *bufio.Scanner, buf []X) (_ []X, err error) {
	for n := 0; n < len(buf); n++ {
		buf[n], err = scanIntX[X](sc)
		if err != nil {
			return buf[:n], err
		}
	}
	return buf, nil
}

func scanTwoIntX[X Int](sc *bufio.Scanner) (X, X, error) {
	var buf [2]X
	_, err := scanInts(sc, buf[:])
	return buf[0], buf[1], err
}

func scanThreeIntX[X Int](sc *bufio.Scanner) (X, X, X, error) {
	var buf [3]X
	_, err := scanInts(sc, buf[:])
	return buf[0], buf[1], buf[2], err
}

func scanFourIntX[X Int](sc *bufio.Scanner) (X, X, X, X, error) {
	var buf [4]X
	_, err := scanInts(sc, buf[:])
	return buf[0], buf[1], buf[2], buf[3], err
}

func scanFiveIntX[X Int](sc *bufio.Scanner) (X, X, X, X, X, error) {
	var buf [5]X
	_, err := scanInts(sc, buf[:])
	return buf[0], buf[1], buf[2], buf[3], buf[4], err
}

var (
	scanInt      = scanIntX[int]
	scanTwoInt   = scanTwoIntX[int]
	scanThreeInt = scanThreeIntX[int]
	scanFourInt  = scanFourIntX[int]
	scanFiveInt  = scanFiveIntX[int]
)

func _appendInt[T Int](b []byte, v T) []byte {
	if ^T(0) < 0 {
		b = strconv.AppendInt(b, int64(v), 10)
	} else {
		b = strconv.AppendUint(b, uint64(v), 10)
	}
	return b
}

func _writeInt[X Int](bw *bufio.Writer, v X) (int, error) {
	if bw.Available() < 24 {
		bw.Flush()
	}
	return bw.Write(_appendInt(bw.AvailableBuffer(), v))
}

type writeOpts struct {
	sep   string
	begin string
	end   string
}

var defaultWriteOpts = writeOpts{
	sep: " ",
	end: "\n",
}

func writeInt[X Int](bw *bufio.Writer, v X, opts ...writeOpts) error {
	var opt writeOpts
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		opt = defaultWriteOpts
	}

	bw.WriteString(opt.begin)
	_writeInt(bw, v)
	_, err := bw.WriteString(opt.end)
	return err
}

func writeInts[X Int](bw *bufio.Writer, a []X, opts ...writeOpts) error {
	var opt writeOpts
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		opt = defaultWriteOpts
	}

	bw.WriteString(opt.begin)

	if len(a) != 0 {
		_writeInt(bw, a[0])
	}

	for i := 1; i < len(a); i++ {
		bw.WriteString(opt.sep)
		_writeInt(bw, a[i])
	}

	_, err := bw.WriteString(opt.end)
	return err
}

// ----------------------------------------------------------------------------

func gcd[I Int](a, b I) I {
	if a > b {
		a, b = b, a
	}
	for a > 0 {
		a, b = b%a, a
	}
	return b
}

func gcdx[I Int](a, b I, x, y *I) I {
	if a == 0 {
		*x = 0
		*y = 1
		return b
	}
	var x1, y1 I
	d := gcdx(b%a, a, &x1, &y1)
	*x = y1 - (b/a)*x1
	*y = x1
	return d
}

func abs[N Sign | Float](a N) N {
	if a < 0 {
		return -a
	}
	return a
}

func sign[N Sign | Float](a N) N {
	if a < 0 {
		return -1
	} else if a > 0 {
		return 1
	}
	return 0
}

// ----------------------------------------------------------------------------

func makeMatrix[T any](n, m int) [][]T {
	buf := make([]T, n*m)
	mx := make([][]T, n)
	for i, j := 0, 0; i < n; i, j = i+1, j+m {
		mx[i] = buf[j : j+m]
	}
	return mx
}
