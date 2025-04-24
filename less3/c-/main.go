package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"unsafe"
)

func solve(nums []int64) []int64 {
	n := len(nums)

	ans := make([]int64, n)
	sizes := make([]int, n) // размеры чисел

	var (
		totalCount   int // общее число еденичек во всех числах
		totalMaxSize int // максимальны размер числа
		xorSum       int64
	)

	var format string
	showNumsAndAns := func() {
		format = "%18d %0" + strconv.Itoa(totalMaxSize) + "b"
		format = format + " " + format
		for i := range nums {
			log.Printf(format, nums[i], nums[i], ans[i], ans[i])
		}
		log.Println("----------")
	}

	if debugEnable {
		defer showNumsAndAns()
	}

	for i, v := range nums {
		count, size := countOnes(v)
		totalCount += count
		totalMaxSize = max(totalMaxSize, size)

		// переносим все единички в нижние разряды и суммируем
		ans[i] = 1<<count - 1
		sizes[i] = count
		xorSum ^= ans[i]
	}

	if debugEnable {
		showNumsAndAns()
	}

	if totalCount&1 != 0 {
		// невозможно распределить нечетное число бит
		return nil
	}

	srcCol := 0 // источник битов, для висячей строки
	bitIdx := 0

	for ; bitIdx < totalMaxSize && xorSum != 0; bitIdx++ {
		if xorSum&(1<<bitIdx) == 0 {
			// если текущий бит нулевой, ничего делать не надо
			continue
		}

		// иначе попробуем перенести одну единичку из текущей колонки вперед,
		// чтобы обнулить текущий бит суммы.
		// Будем двигать бит числа минимальной длины размер которого > i

		var (
			minSize      = totalMaxSize + 1
			maxSize      = 0
			maxSizeCount = 0
			minSizeIdx   = -1
		)

		for j, size := range sizes {
			if size <= bitIdx {
				// В этой позиции у числа только лидируещие 0
				continue
			}
			if size > maxSize {
				maxSize = size
				maxSizeCount = 1
			} else if size == maxSize {
				maxSizeCount++
			}
			if size < minSize {
				minSize = size
				minSizeIdx = j
			}
		}

		if minSizeIdx == -1 {
			// oops?!.. нечего двигать
			return nil
		}

		if minSize == maxSize && maxSizeCount == 1 {
			// висячая строка

			if bitIdx+1 >= totalMaxSize {
				// чтобы поправить текущую позицию, нам нужна
				// покрайней мере одна позиция впереди
				return nil
			}

			// найдем две единички в одной из предыдущих колонок
			row1, row2 := -1, -1
			for ; srcCol < maxSize; srcCol++ {

				for row := 0; row < len(ans); row++ { // XXX неоптимальный старт
					if ans[row]&(int64(1)<<srcCol) != 0 {
						if row1 == -1 {
							row1 = row
							continue
						}
						if row2 == -1 {
							row2 = row
							break
						}
					}
				}

				if row1 != -1 && row2 != -1 {
					// bingo!
					break
				}

				row1, row2 = -1, -1
			}

			if row1 == -1 || row2 == -1 {
				// oops!.. ничего не нашли
				return nil
			}

			xorSum ^= ans[row1]
			xorSum ^= ans[row2]

			ans[row1] &^= int64(1) << srcCol
			ans[row1] |= int64(1) << bitIdx

			ans[row2] &^= int64(1) << srcCol
			ans[row2] |= int64(1) << (bitIdx + 1)

			xorSum ^= ans[row1]
			xorSum ^= ans[row2]
			continue
		}

		if minSize == totalMaxSize {
			// некуда двигать, уперлись в размер
			return nil
		}

		// извлекаем число из суммы
		xorSum ^= ans[minSizeIdx]

		// переносим текущий бит вперед, перед первым значащим битом числа
		ans[minSizeIdx] &^= int64(1) << bitIdx
		ans[minSizeIdx] |= int64(1) << sizes[minSizeIdx]
		sizes[minSizeIdx]++

		// возвращаем число в сумму
		xorSum ^= ans[minSizeIdx]
	}

	return ans
}

func countOnes(x int64) (int, int) {
	var size, count int
	for x > 0 {
		size++
		count += int(x & 1)
		x >>= 1
	}
	return count, size
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

	nums := make([]int64, n)
	if err := scanInts(sc, nums); err != nil {
		panic(err)
	}

	ans := solve(nums)

	if ans != nil {
		writeInts(bw, ans, defaultWriteOpts())
	} else {
		bw.WriteString("impossible\n")
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
