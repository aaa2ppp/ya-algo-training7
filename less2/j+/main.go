package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"unsafe"
)

type item struct {
	val     uint64
	len     int
	promise int32
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
const (
	multiplier = uint64(100069)
)

var (
	power    []uint64
	powerSum []uint64
)

func precalc(nums []int32) []item {
	n := 1
	for n < len(nums) {
		n *= 2
	}

	power = make([]uint64, n)
	power[0] = 1
	for i := 1; i < len(power); i++ {
		power[i] = power[i-1] * multiplier
	}

	powerSum = make([]uint64, n+1)
	powerSum[0] = 0
	for i := 0; i < len(power); i++ {
		powerSum[i+1] += powerSum[i] + power[i]
	}

	tree := make([]item, n*2-1)
	for i, j := 0, n-1; i < len(nums); i, j = i+1, j+1 {
		tree[j] = item{
			val: uint64(nums[i]),
			len: 1,
		}
	}

	// for i := (n - 1) + len(nums); i < len(tree); i++ {
	// 	tree[i] = item{val: 0, len: 0}
	// }

	for i := n - 2; i >= 0; i-- {
		a, b := tree[i*2+1], tree[i*2+2]
		tree[i] = calcNode(a, b)
	}

	return tree
}

func calcNode(a, b item) item {
	return item{
		val: a.val*power[b.len] + b.val,
		len: a.len + b.len,
	}
}

func updateNode(node *item, val int32) {
	node.val = uint64(val) * powerSum[node.len]
	node.promise = val
}

func fulfillPromise(tree []item, i int) {
	if promise := tree[i].promise; promise != 0 {
		updateNode(&tree[i*2+1], promise)
		updateNode(&tree[i*2+2], promise)
		tree[i].promise = 0
	}
}

func update(tree []item, ql, qr int, val int32) {
	n := (len(tree) + 1) / 2

	var dfs func(i, l, r int) bool

	dfs = func(i, l, r int) bool {
		if debugEnable {
			log.Println("dfs:", i, l, r, ql, qr)
		}

		// если интервал не пересекается с интервалом узла
		if qr <= l || r <= ql {
			return false
		}

		// если интервал полностью накрывает интервал узла
		if ql <= l && r <= qr {
			// обновим узел и пообещаем обновить детей
			updateNode(&tree[i], val)
			return true
		}

		// иначе частичное пересечение

		// выполним ранее данное обещание
		fulfillPromise(tree, i)

		// и обновим детей
		m := (l + r) / 2
		aUpd := dfs(i*2+1, l, m)
		bUpd := dfs(i*2+2, m, r)
		if aUpd || bUpd {
			a, b := tree[i*2+1], tree[i*2+2]
			tree[i] = calcNode(a, b)
		}

		return true
	}

	dfs(0, 0, n)
}

func query(tree []item, ql, qr int) uint64 {
	n := (len(tree) + 1) / 2

	var dfs func(i, l, r int) item

	dfs = func(i, l, r int) item {
		if debugEnable {
			log.Println("dfs:", i, l, r, ql, qr)
		}

		// если интервал не пересекается с интервалом узла
		if qr <= l || r <= ql {
			return item{}
		}

		// если интервал полностью накрывает интервал узла
		if ql <= l && r <= qr {
			return tree[i]
		}

		// выполним обещание
		fulfillPromise(tree, i)

		// и спросим детей
		m := (l + r) / 2
		a := dfs(i*2+1, l, m)
		b := dfs(i*2+2, m, r)
		return calcNode(a, b)
	}

	ans := dfs(0, 0, n)
	return ans.val
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
			if debugEnable {
				log.Println("== u", l, r, k)
			}
			// преходим из 1-base [l, r] в 0-base индексацию [l, r)
			update(tree, l-1, r, int32(k))

		case 1: // compare
			if debugEnable {
				log.Println("== c", l, r, k)
			}
			// преходим из 1-base в 0-base индексацию
			l1, l2 := l-1, r-1
			r1, r2 := l1+k, l2+k
			a := query(tree, l1, r1)
			b := query(tree, l2, r2)
			if a == b {
				bw.WriteByte('+')
			} else {
				bw.WriteByte('-')
			}
		default:
			panic("unknown operation " + strconv.Itoa(t))
		}
	}
	bw.WriteByte('\n')
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
