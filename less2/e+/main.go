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
	left  int
	total int
}

func precalc(nums []int32) []item {
	n := 1
	for n < len(nums) {
		n *= 2
	}

	tree := make([]item, n*2-1)
	for i, j := 0, n-1; i < len(nums); i, j = i+1, j+1 {
		if nums[i] == 0 {
			tree[j] = item{total: 1}
		}
	}

	// for i := n - 1 + len(nums); i < len(tree); i++ {
	// 	tree[i] = item{}
	// }

	for i := n - 2; i >= 0; i-- {
		a, b := tree[i*2+1], tree[i*2+2]
		tree[i] = item{
			left:  a.total,
			total: a.total + b.total,
		}
	}

	return tree
}

// query ищет k ноль на *открытом* 0-base интервале [l, r)
func query(tree []item, ql, qr, k int) int {

	// первые три числа определяют узел: индекс и контролируемый интервал
	var dfs func(i, l, r, ql, qr, k int) (int, int)

	dfs = func(i, l, r, ql, qr, k int) (int, int) {
		if debugEnable {
			log.Println("dfs:", i, l, r, ql, qr, k)
		}

		// если интервал не пересекаются с моим
		if qr <= l || r <= ql {
			// у меня нет k-го нуля и у меня нет ни одного нуля на этом интервале
			return -1, 0
		}

		// если интервал полностью накрывает мой
		if ql <= l && r <= qr {

			// если у меня не хватает нулей
			if tree[i].total < k {
				// у меня нет k-го нуля, но есть total нулей на этом интервале
				return -1, tree[i].total
			}

			// k-й ноль *точно* есть

			// если я состоитою только из нулей
			if tree[i].total == r-l {
				// bingo!
				return l + k - 1, 0
			}

			// k-й ноль или в левом или в правом ребенке

			m := (l + r) / 2
			if tree[i].left >= k {
				// *точно* в левом
				return dfs(i*2+1, l, m, ql, qr, k)
			} else {
				// *точно* в правом
				return dfs(i*2+2, m, r, ql, qr, k-tree[i].left)
			}
		}

		// есть неполное пересечение интервала с моим.
		// ничего определенного сказать не могу, спрошу у детей...
		m := (l + r) / 2

		// спрошу у левого ребенка
		idx, lCnt := dfs(i*2+1, l, m, ql, qr, k)
		if idx != -1 {
			// bingo!
			return idx, 0
		}

		// спросим у правого ребенка
		idx, rCnt := dfs(i*2+2, m, r, ql, qr, k-lCnt)
		if idx != -1 {
			// bingo!
			return idx, 0
		}

		// у меня нет k-го нуля, но я у меня есть несколько нулей на этом интервале
		return -1, lCnt + rCnt
	}

	n := (len(tree) + 1) / 2
	idx, _ := dfs(0, 0, n, ql, qr, k)
	return idx
}

func update(tree []item, i int, val int32) {
	n := (len(tree) + 1) / 2
	i += n - 1
	if val == 0 {
		tree[i].total = 1
	} else {
		tree[i].total = 0
	}
	for i > 0 {
		i = (i - 1) / 2 // переходим к родителю
		a, b := tree[i*2+1], tree[i*2+2]
		tree[i] = item{
			left:  a.total,
			total: a.total + b.total,
		}
	}
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

	if debugEnable {
		log.Println(nums)
	}

	tree := precalc(nums)

	if debugEnable {
		log.Println(tree)
	}

	m, err := scanInt(sc)
	if err != nil {
		panic(err)
	}

	for i := 0; i < m; i++ {
		op, err := scanWord(sc)
		if err != nil {
			panic(err)
		}
		switch op {
		case "s":
			l, r, k, err := scanThreeInt(sc)
			if err != nil {
				panic(err)
			}
			// делаем из закрытого 1-base интервала
			// открытый 0-base интервал [l, r)
			ans := query(tree, l-1, r, k)
			if ans != -1 {
				ans++
			}
			if debugEnable {
				log.Printf("q %d %d %d -> %d", l, r, k, ans)
			}
			writeInt(bw, ans, writeOpts{end: ' '})
		case "u":
			i, v, err := scanTwoInt(sc)
			if err != nil {
				panic(err)
			}
			// преходим из 1-base в 0-base индексацию
			update(tree, i-1, int32(v))
			if debugEnable {
				log.Printf("u %d %d -> %v", i, v, tree)
			}
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

// func max[T Ordered](a, b T) T {
// 	if a < b {
// 		return b
// 	}
// 	return a
// }

// func min[T Ordered](a, b T) T {
// 	if a > b {
// 		return b
// 	}
// 	return a
// }

// ----------------------------------------------------------------------------

func makeMatrix[T any](n, m int) [][]T {
	buf := make([]T, n*m)
	matrix := make([][]T, n)
	for i, j := 0, 0; i < n; i, j = i+1, j+m {
		matrix[i] = buf[j : j+m]
	}
	return matrix
}
