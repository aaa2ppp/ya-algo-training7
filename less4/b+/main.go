package main

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"unsafe"
)

type Stack[T any] []T

func (s Stack[T]) Len() int {
	return len(s)
}

func (s Stack[T]) Empty() bool {
	return s.Len() == 0
}

func (s *Stack[T]) Push(v T) {
	*s = append(*s, v)
}

func (s Stack[T]) Back() T {
	if s.Empty() {
		panic("stack is empty")
	}
	return s[len(s)-1]
}

func (s *Stack[T]) Pop() T {
	v := s.Back()
	old := *s
	*s = old[:len(old)-1]
	return v
}

func (s *Stack[T]) Clear() {
	old := *s
	*s = old[:0]
}

type Queue[T any] struct {
	in  Stack[T]
	out Stack[T]
}

func (q Queue[T]) Len() int {
	return q.in.Len() + q.out.Len()
}

func (q Queue[T]) Empty() bool {
	return q.Len() == 0
}

func (q *Queue[T]) Push(v T) {
	q.in.Push(v)
}

func (q *Queue[T]) Front() T {
	if q.Empty() {
		panic("queue is empty")
	}
	if q.out.Empty() {
		for !q.in.Empty() {
			q.out.Push(q.in.Pop())
		}
	}
	return q.out.Back()
}

func (q *Queue[T]) Pop() T {
	v := q.Front()
	q.out.Pop()
	return v
}

func (q *Queue[T]) Clear() {
	q.in.Clear()
	q.out.Clear()
}

func run(in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	var queue Queue[int]
	for {
		cmd, err := scanWord(sc)
		if err != nil {
			panic(err)
		}

		switch cmd {
		case "push":
			// Добавить в очередь число n (значение n задается после команды).
			// Программа должна вывести ok.
			v, err := scanInt(sc)
			if err != nil {
				panic(err)
			}
			queue.Push(v)
			bw.WriteString("ok\n")

		case "pop":
			// Удалить из очереди первый элемент. Программа должна вывести его значение.
			if queue.Empty() {
				bw.WriteString("error\n")
			} else {
				writeInt(bw, queue.Pop())
			}

		case "front":
			// Программа должна вывести значение первого элемента, не удаляя его из очереди.
			if queue.Empty() {
				bw.WriteString("error\n")
			} else {
				writeInt(bw, queue.Front())
			}

		case "size":
			// Программа должна вывести количество элементов в очереди.
			writeInt(bw, queue.Len())

		case "clear":
			// Программа должна очистить очередь и вывести ok.
			queue.Clear()
			bw.WriteString("ok\n")

		case "exit":
			// Программа должна вывести bye и завершить работу.
			bw.WriteString("bye\n")
			return
		}
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
