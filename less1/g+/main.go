package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"unsafe"
)

// У Пети есть набор из N кирпичиков. Каждый кирпичик полностью окрашен в один из K цветов,
// i-й кирпичик имеет размер 1×1×L. Петя знает, что он может построить из кирпичиков
// прямоугольную стену толщиной 1 и высотой K, причем первый горизонтальный слой кирпичиков
// в стене будет первого цвета, второй — второго и т. д. Теперь Петя хочет узнать, может ли
// он из своего набора построить две прямоугольные стены, обладающие тем же свойством.
// Помогите ему выяснить это.

type brick struct {
	len int // длина кирпича
	col int // цвет кирпича (1-indexing)
}

// solve возвращает true, если можно построить две стены и список (1-indexing) кирпичей,
// из которых следует построить первую стену
func solve(k int, bricks []brick) (bool, []int) {

	// будем решать задачу о рюкзаке для кажого цвета кирпичей

	// базовая длина стены (считаем по первому цвету)
	baseLen := 0
	for _, b := range bricks {
		if b.col == 1 {
			baseLen += b.len
		}
	}

	dp := makeMatrix[int](k, baseLen+1)
	for c := 0; c < k; c++ {
		for x := 1; x <= baseLen; x++ {
			dp[c][x] = -1
		}
	}

	cnt := make([]int, baseLen+1) // количество цветов для которых достижима длина
	endLen := make([]int, k)      // последняя достижимая длина для цвета
	needLen := baseLen            // ищем первую длину, которая достижима для всех цветов

mainLoop:
	for idx, b := range bricks {
		c := b.col - 1 // c to 0-indexing
		endLen[c] = min(endLen[c]+b.len, baseLen)

		if debugEnable {
			log.Println("b:", b, "endLen:", endLen[c])
		}

		for x := endLen[c]; x > 0; x-- {
			if x-b.len >= 0 && dp[c][x] == -1 && dp[c][x-b.len] != -1 {
				dp[c][x] = idx
				cnt[x]++
				if cnt[x] == k {
					// bingo!
					needLen = x
					break mainLoop
				}
			}
		}
	}

	if debugEnable {
		for _, row := range dp {
			log.Println(row)
		}
	}

	if needLen == baseLen {
		return false, nil
	}

	var ans []int
	for c := 0; c < k; c++ {
		for x := needLen; x > 0; {
			idx := dp[c][x]
			ans = append(ans, idx+1) // idx to 1-indexing
			x -= bricks[idx].len
		}
	}

	return true, ans
}

func run(in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, k, err := scanTwoInt(sc)
	if err != nil {
		panic(err)
	}

	bricks := make([]brick, 0, n)
	for i := 0; i < n; i++ {
		l, c, err := scanTwoInt(sc)
		if err != nil {
			panic(err)
		}
		bricks = append(bricks, brick{l, c})
	}

	if debugEnable {
		log.Println("bricks:", bricks)
	}

	ok, ans := solve(k, bricks)
	if ok {
		bw.WriteString("YES\n")
		writeInts(bw, ans, defaultWriteOpts())
	} else {
		bw.WriteString("NO\n")
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
