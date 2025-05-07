package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"unsafe"
)

type node struct {
	prev *node
	next *node
	val  int
}

type river struct {
	cur    *node
	idx    int
	square int
}

func (r *river) append(val int) {
	r.square += val * val
	node := &node{val: val}
	if r.cur == nil {
		r.cur = node
		r.idx = 0
		return
	}

	r.cur.next = node
	node.prev = r.cur
	r.cur = node
	r.idx++
}

func (r *river) goIdx(i int) {
	if r.idx < i {
		for r.idx != i {
			r.cur = r.cur.next
			r.idx++
		}
	} else if r.idx > i {
		for r.idx != i {
			r.cur = r.cur.prev
			r.idx--
		}
	}
}

func (r *river) split(i int) {
	r.goIdx(i)

	r.square -= r.cur.val * r.cur.val
	val1 := r.cur.val / 2
	val2 := r.cur.val - val1
	r.square += val1*val1 + val2*val2

	newNode := &node{
		prev: r.cur,
		next: r.cur.next,
		val:  val2,
	}
	if r.cur.next != nil {
		r.cur.next.prev = newNode
	}
	r.cur.next = newNode
	r.cur.val = val1
}

func (r *river) remove(i int) {
	r.goIdx(i)

	r.square -= r.cur.val * r.cur.val
	val1 := r.cur.val / 2
	val2 := r.cur.val - val1

	if prev := r.cur.prev; prev != nil {
		r.square -= prev.val * prev.val
		prev.val += val1
		if r.cur.next == nil {
			prev.val += val2
		}
		r.square += prev.val * prev.val
		prev.next = r.cur.next
	}

	if next := r.cur.next; next != nil {
		r.square -= next.val * next.val
		next.val += val2
		if r.cur.prev == nil {
			next.val += val1
		}
		r.square += next.val * next.val
		next.prev = r.cur.prev
	}

	if r.cur.next != nil {
		r.cur = r.cur.next
	} else {
		r.cur = r.cur.prev
		r.idx--
	}
}

func run(in io.Reader, out io.Writer) {
	log.SetFlags(0)
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, err := scanInt(sc)
	if err != nil {
		panic(err)
	}

	var r river
	for range n {
		v, err := scanInt(sc)
		if err != nil {
			panic(err)
		}
		r.append(v)
	}

	k, err := scanInt(sc)
	if err != nil {
		panic(err)
	}

	if debugEnable {
		log.Println("r:", r)
	}
	writeInt(bw, r.square)

	for range k {
		e, i, err := scanTwoInt(sc)
		if err != nil {
			panic(err)
		}

		switch e {
		case 1:
			r.remove(i - 1) // to 0-indexing
		case 2:
			r.split(i - 1) // to 0-indexing
		default:
			panic("unknown evant " + strconv.Itoa(e))
		}

		if debugEnable {
			log.Println("r:", r)
		}
		writeInt(bw, r.square)
	}
}

// ----------------------------------------------------------------------------

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	run(os.Stdin, os.Stdout)
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
	n := 0
	for ; n < len(buf) && err == nil; n++ {
		buf[n], err = scanIntX[X](sc)
	}
	return buf[:n], err
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
