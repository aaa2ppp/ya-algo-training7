package main

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"unsafe"
)

type item struct {
	val uint64
	len int
}

// Простые числа от 100000 - выбирай не хочу
// 100003 100019 100043 100049 100057 100069 100103 100109 100129 100151
// 100153 100169 100183 100189 100193 100207 100213 100237 100267 100271
// 100279 100291 100297 100313 100333 100343 100357 100361 100363 100379
// 100391 100393 100403 100411 100417 100447 100459 100469 100483 100493
// 100501 100511 100517 100519 100523 100537 100547 100549 100559 100591
// 100609 100613 100621 100649 100669 100673 100693 100699 100703 100733
// 100741 100747 100769 100787 100799 100801 100811 100823 100829 100847
// 100853 100907 100913 100927 100931 100937 100943 100957 100981 100987
// 100999 101009 101021 101027 101051 101063 101081 101089 101107 101111
// 101113 101117 101119 101141 101149 101159 101161 101173 101183 101197
// 101203 101207 101209 101221 101267 101273 101279 101281 101287 101293
// 101323 101333 101341 101347 101359 101363 101377 101383 101399 101411
// 101419 101429 101449 101467 101477 101483 101489 101501 101503 101513
// ...
const multiplier = uint64(100069)

var power []uint64

func precalc(nums []int32) []item {
	n := 1
	for n <= len(nums) { // нам нужен +1 символ в начале под 1
		n *= 2
	}

	power = make([]uint64, n+1)
	power[0] = 1
	for i := 1; i < len(power); i++ {
		power[i] = power[i-1] * multiplier
	}

	tree := make([]item, n*2)
	tree[n-1] = item{val: 1, len: 1}
	for i, j := 0, n; i < len(nums); i, j = i+1, j+1 {
		tree[j] = item{val: uint64(nums[i]), len: 1}
	}

	for i := n + len(nums); i < len(tree); i++ {
		tree[i] = item{val: 1, len: 1}
	}

	for i := n - 2; i >= 0; i-- {
		a, b := tree[i*2+1], tree[i*2+2]
		tree[i] = item{
			val: a.val*power[b.len] + b.val,
			len: a.len + b.len,
		}
	}

	return tree
}

func update(tree []item, ul, ur int, val int32) {
	n := (len(tree) + 1) / 2
	_ = n
	// tree[i] = val

	// for i > 0 {
	// 	i = (i - 1) / 2 // переходим к родителю
	// 	a, b := tree[i*2+1], tree[i*2+2]
	// 	tree[i] = max(a, b)
	// }
}

func compare(tree []item, l1, l2, size int) bool {

	return false // TODO
}

// query возвращает хешсумму на префиксе
func query(tree []item, i int) uint64 {
	n := (len(tree) + 1) / 2
	_ = n

	// var dfs func(i, l, r int) int

	// dfs = func(i, l, r int) int {
	// 	if tree[i] < target {
	// 		// мой максимум меньше цели
	// 		return -1
	// 	}

	// 	if r-l == 1 {
	// 		// я лист
	// 		return l
	// 	}

	// 	// иначе спрошу у детей
	// 	ans := -1
	// 	m := (l + r) / 2

	// 	if ans == -1 && idx < m {
	// 		ans = dfs(i*2+1, l, m)
	// 	}

	// 	if ans == -1 && idx < r {
	// 		ans = dfs(i*2+2, m, r)
	// 	}

	// 	return ans
	// }

	// return dfs(0, 0, n)
	return 0 // TODO
}

func run(in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, err := scanInt(sc)
	if err != nil {
		panic(err)
	}

	nums := make([]int32, n)
	if err := scanInts(sc, nums); err != nil {
		panic(err)
	}

	tree := precalc(nums)

	q, err := scanInt(sc)
	if err != nil {
		panic(err)
	}

	for i := 0; i < q; i++ {
		t, l, r, k, err := scanFourInt(sc)
		if err != nil {
			panic(err)
		}
		switch t {
		case 0: // update
			// преходим из 1-base [l, r] в 0-base индексацию [l, r)
			update(tree, l-1, r, int32(k))

		case 1: // compare
			l1, l2 := l, r
			// преходим из 1-base в 0-base индексацию
			if compare(tree, l1-1, l2-1, k) {
				bw.WriteByte('+')
			} else {
				bw.WriteByte('-')
			}
		default:
			panic("unknown operation " + strconv.Itoa(t))
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
