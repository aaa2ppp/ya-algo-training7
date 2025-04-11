package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"unsafe"
)

// Вася и Маша организовали производство бирдекелей. Они договорились работать день через день,
// при этом Вася любит простую и понятную работу, а Маша — сложную и творческую. В первый день
// работает Вася, во второй — Маша, потом снова Вася и т.д.

// К ним поступило N заказов, для каждого заказа известна его продолжительность в днях и для
// каждого из дней известно, будет ли в этот день работа сложной или простой. Заказы можно
// выполнять в любом порядке, перерывов между заказами нет.

// Определите такой порядок выполнения заказов, чтобы Вася получил как можно больше простых
// задач, а Маша — сложных.

// Выведите одно число — максимальное количество дней с простой работой, которые достанутся Васе.

func countRigthLeft(order string) (right, left int) {
	for i, c := range []byte(order) {
		if c == 'S' {
			right += (i + 1) & 1
			left += i & 1
		}
	}
	return right, left
}

type solveFunc func(order []string) int

func solve(orders []string) (count int) {

	// все что нам надо знать о заказе
	type item struct {
		order string
		right int // профит (Васи), если выполнение заказа начинается с четного дня (0,2,4,...)
		left  int // профит (Васи), если выполнение заказа начинается с нечетного дня (1,3,5,...)
	}

	var res strings.Builder
	if debugEnable {
		defer func() {
			s := res.String()
			log.Println("=", count, s[:len(s)-1])
		}()
	}

	addRight := func(it item) {
		if debugEnable {
			res.WriteString(it.order)
			res.WriteByte('.')
			log.Println("R", it.order, it.right)
		}
		count += it.right
	}

	addLeft := func(it item) {
		if debugEnable {
			res.WriteString(it.order)
			res.WriteByte('.')
			log.Println("L", it.order, it.left)
		}
		count += it.left
	}

	// разделяем заказы по времени выполнения (четный/нечетный) и по четности дня начала
	// NOTE: только заказы нечетной длины, могут менять четность дня начала следующего заказа
	var (
		evenRight []item // четной длины, которые лучше начинать с четного дня
		evenLeft  []item // четной длины, которые лучше начинать с нечетного дня
		oddRight  []item // нечетной дины, которые лучше начинать с четного дня
		oddLeft   []item // нечетной дины, которые лучше начинать с нечетного дня
	)
	for _, order := range orders {
		right, left := countRigthLeft(order)
		item := item{
			order,
			right,
			left,
		}
		if len(order)%2 != 0 {
			if right > left {
				oddRight = append(oddRight, item)
			} else {
				oddLeft = append(oddLeft, item)
			}
		} else {
			if right > left {
				evenRight = append(evenRight, item)
			} else {
				evenLeft = append(evenLeft, item)
			}
		}
	}

	if debugEnable {
		log.Println("evanRight:", evenRight)
		log.Println("evanLeft :", evenLeft)
	}

	// addRight all evenRight
	for _, item := range evenRight {
		addRight(item)
	}

	if len(oddRight)+len(oddLeft) == 0 {
		// addRight all evenLeft
		for _, item := range evenLeft {
			addRight(item)
		}
		return count
	}

	addLeftAllEvenLeft := func() {
		for _, item := range evenLeft {
			addLeft(item)
		}
	}

	slices.SortFunc(oddRight, func(a, b item) int {
		dif := func(it item) int { return it.right - it.left }
		return dif(b) - dif(a)
	})

	slices.SortFunc(oddLeft, func(a, b item) int {
		dif := func(it item) int { return it.left - it.right }
		return dif(b) - dif(a)
	})

	if debugEnable {
		log.Println("oddRight :", oddRight)
		log.Println("oddLeft  :", oddLeft)
	}

	i := 0
	for ; i < len(oddRight) && i < len(oddLeft); i++ {
		addRight(oddRight[i])
		if i == 0 {
			addLeftAllEvenLeft()
		}
		addLeft(oddLeft[i])
	}

	for j := len(oddRight) - 1; i <= j; i, j = i+1, j-1 {
		addRight(oddRight[i])
		if i == 0 {
			addLeftAllEvenLeft()
		}
		if i < j {
			addLeft(oddRight[j])
		}
	}

	for j := len(oddLeft) - 1; i <= j; i, j = i+1, j-1 {
		addRight(oddLeft[j])
		if i == 0 {
			addLeftAllEvenLeft()
		}
		if i < j {
			addLeft(oddLeft[i])
		}
	}

	return count
}

func run(in io.Reader, out io.Writer, solve solveFunc) {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, err := scanInt(sc)
	if err != nil {
		panic(err)
	}

	orders := make([]string, n)
	for i := 0; i < n; i++ {
		s, err := scanWord(sc)
		if err != nil {
			panic(err)
		}
		orders[i] = s
	}

	ans := solve(orders)
	writeInt(bw, ans, defaultWriteOpts())
}

// ----------------------------------------------------------------------------

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	run(os.Stdin, os.Stdout, solve)
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
	~int | ~int64 | ~int32 | ~int16 | ~int8
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

func gcdx(a, b int, x, y *int) int {
	if a == 0 {
		*x = 0
		*y = 1
		return b
	}
	var x1, y1 int
	d := gcdx(b%a, a, &x1, &y1)
	*x = y1 - (b/a)*x1
	*y = x1
	return d
}

func abs[N Number](a N) N {
	if a < 0 {
		return -a
	}
	return a
}

func sign[N Number](a N) N {
	if a < 0 {
		return -1
	} else if a > 0 {
		return 1
	}
	return 0
}

type Ordered interface {
	Number | ~string
}

func max[T Ordered](a, b T) T {
	if a < b {
		return b
	}
	return a
}

func min[T Ordered](a, b T) T {
	if a > b {
		return b
	}
	return a
}

// ----------------------------------------------------------------------------

func makeMatrix[T any](n, m int) [][]T {
	buf := make([]T, n*m)
	matrix := make([][]T, n)
	for i, j := 0, 0; i < n; i, j = i+1, j+m {
		matrix[i] = buf[j : j+m]
	}
	return matrix
}
