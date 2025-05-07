package main

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"unsafe"
)

const dequeMinLen = 8

type Deque[T any] struct {
	data  []T
	size  int
	first int
}

func (q Deque[T]) Len() int {
	return q.size
}

func (q Deque[T]) Empty() bool {
	return q.Len() == 0
}

func (q *Deque[T]) grow(n int) {
	old := q.data
	n = max(dequeMinLen, 2*len(old)+n)
	q.data = append(append(make([]T, 0, n), old[q.first:]...), old[0:q.first]...)[:n]
	q.first = 0
}

func (q *Deque[T]) Grow(n int) {
	if n < 0 {
		panic("Deque.Grow: negative count")
	}
	if len(q.data)-q.size < n {
		q.grow(n)
	}
}

func (q *Deque[T]) PushFront(v T) {
	if q.size == len(q.data) {
		q.grow(0)
	}
	q.first--
	if q.first == -1 {
		q.first = len(q.data) - 1
	}
	q.data[q.first] = v
	q.size++
}

func (q *Deque[T]) PushBack(v T) {
	if q.size == len(q.data) {
		q.grow(0)
	}
	i := q.first + q.size
	if i >= len(q.data) {
		i -= len(q.data)
	}
	q.data[i] = v
	q.size++
}

func (q Deque[T]) Front() T {
	if q.Empty() {
		panic("queue is empty")
	}
	return q.data[q.first]
}

func (q Deque[T]) Back() T {
	if q.Empty() {
		panic("queue is empty")
	}
	i := q.first + q.size - 1
	if i >= len(q.data) {
		i -= len(q.data)
	}
	return q.data[i]
}

func (q *Deque[T]) PopFront() T {
	v := q.Front()
	q.first++
	if q.first == len(q.data) {
		q.first = 0
	}
	q.size--
	return v
}

func (q *Deque[T]) PopBack() T {
	v := q.Back()
	q.size--
	return v
}

func (q *Deque[T]) Clear() {
	q.size = 0
	clear(q.data) // for case when T contains pointers
}

func run(in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	var deque Deque[int]
	for {
		cmd, err := scanWord(sc)
		if err != nil {
			panic(err)
		}

		switch cmd {
		case "push_front":
			// Добавить (положить) в начало дека новый элемент. Программа должна вывести ok.
			v, err := scanInt(sc)
			if err != nil {
				panic(err)
			}
			deque.PushFront(v)
			bw.WriteString("ok\n")

		case "push_back":
			// Добавить (положить) в конец дека новый элемент. Программа должна вывести ok.
			v, err := scanInt(sc)
			if err != nil {
				panic(err)
			}
			deque.PushBack(v)
			bw.WriteString("ok\n")

		case "pop_front":
			// Извлечь из дека первый элемент. Программа должна вывести его значение.
			if deque.Empty() {
				bw.WriteString("error\n")
			} else {
				writeInt(bw, deque.PopFront())
			}

		case "pop_back":
			// Извлечь из дека последний элемент. Программа должна вывести его значение.
			if deque.Empty() {
				bw.WriteString("error\n")
			} else {
				writeInt(bw, deque.PopBack())
			}

		case "front":
			// Узнать значение первого элемента (не удаляя его). Программа должна вывести его значение.
			if deque.Empty() {
				bw.WriteString("error\n")
			} else {
				writeInt(bw, deque.Front())
			}

		case "back":
			// Узнать значение последнего элемента (не удаляя его). Программа должна вывести его значение.
			if deque.Empty() {
				bw.WriteString("error\n")
			} else {
				writeInt(bw, deque.Back())
			}

		case "size":
			// Вывести количество элементов в деке.
			writeInt(bw, deque.Len())

		case "clear":
			// Очистить дек (удалить из него все элементы) и вывести ok.
			deque.Clear()
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
