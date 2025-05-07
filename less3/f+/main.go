package main

import (
	"bufio"
	"io"
	"iter"
	"log"
	"math"
	"math/bits"
	"os"
	"strconv"
	"unsafe"
)

type points iter.Seq2[int, [3]int]

type solveFunc func(n, k int, points points) (bool, [3]int)

func slowSolve(n, k int, points points) (bool, [3]int) {
	xy := makeMatrix[byte](n, n)
	xz := makeMatrix[byte](n, n)
	yz := makeMatrix[byte](n, n)

	var xyN, xzN, yzN int

	for _, p := range points {
		x, y, z := p[0]-1, p[1]-1, p[2]-1 // to 0-indexing

		if xy[x][y] == 0 {
			xyN++
			xy[x][y] = 1
		}

		if xz[x][z] == 0 {
			xzN++
			xz[x][z] = 1
		}

		if yz[y][z] == 0 {
			yzN++
			yz[y][z] = 1
		}
	}

	if debugEnable {
		log.Println("xy:", xyN)
		// for _, row := range xy {
		// 	log.Println(row)
		// }
		log.Println("xz:", xzN)
		// for _, row := range xz {
		// 	log.Println(row)
		// }
		log.Println("yz:", yzN)
		// for _, row := range yz {
		// 	log.Println(row)
		// }
	}

	// TODO: убрать копипаст
	if xyN >= xzN && xyN >= yzN {
		for x := 0; x < n; x++ {
			for y := 0; y < n; y++ {
				if xy[x][y] != 0 {
					continue
				}
				for z := 0; z < n; z++ {
					if xz[x][z] != 0 || yz[y][z] != 0 {
						continue
					}
					return false, [3]int{x + 1, y + 1, z + 1} // to 1-indexing
				}
			}
		}
	} else if xzN >= xyN && xzN >= yzN {
		for x := 0; x < n; x++ {
			for z := 0; z < n; z++ {
				if xz[x][z] != 0 {
					continue
				}
				for y := 0; y < n; y++ {
					if xy[x][y] != 0 || yz[y][z] != 0 {
						continue
					}
					return false, [3]int{x + 1, y + 1, z + 1} // to 1-indexing
				}
			}
		}
	} else { // yzN >= xyN && yzN >= xzN
		for y := 0; y < n; y++ {
			for z := 0; z < n; z++ {
				if yz[y][z] != 0 {
					continue
				}
				for x := 0; x < n; x++ {
					if xy[x][y] != 0 || xz[x][z] != 0 {
						continue
					}
					return false, [3]int{x + 1, y + 1, z + 1} // to 1-indexing
				}
			}
		}
	}

	return true, [3]int{}
}

type bitArray []uint64

func makeBitArray(n int) bitArray {
	arr := make([]uint64, (n+63)>>6)

	// Если есть хвостик, то зальем его единичками.
	// Пригодиться, когда будем делать OR и искать нули
	if k := n & 63; k != 0 {
		arr[len(arr)-1] = math.MaxInt64 << k
	}

	return arr
}

func (arr bitArray) isSet(i int) bool {
	return arr[i>>6]&(1<<(i&63)) != 0
}

func (arr bitArray) set(i int) {
	arr[i>>6] |= 1 << (i & 63)
}

// zeros возвращает индексы всех нулей в массиве
func (arr bitArray) zeros() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i, v := range arr {
			if v == math.MaxUint64 {
				// все биты единички
				continue
			}

			// // ищем нулевые биты
			// for j := 0; j < 64; j++ {
			// 	if v&1 == 0 {
			// 		if !yield(i*64 + j) {
			// 			return
			// 		}
			// 	}
			// 	v >>= 1
			// }

			u := ^v
			j := 0
			for u > 0 {
				cnt := bits.TrailingZeros64(u)
				j += cnt
				if (v>>j)&1 != 0 {
					panic("I'm a sucker")
				}
				if !yield(i*64 + j) {
					return
				}
				u >>= cnt + 1
				j++
			}
		}
	}
}

// orAndZeros делает побитовое OR с другим массивом. Возвращает индесы всех нулей,
// в результирующем массиве.
func (arr bitArray) orAndZeros(other bitArray) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := range arr {
			v := arr[i] | other[i]
			if v == math.MaxUint64 {
				// все биты единички
				continue
			}

			// ищем нулевые биты

			// for j := 0; j < 64; j++ {
			// 	if v&1 == 0 {
			// 		if !yield(i*64 + j) {
			// 			return
			// 		}
			// 	}
			// 	v >>= 1
			// }

			u := ^v
			j := 0
			for u > 0 {
				cnt := bits.TrailingZeros64(u)
				j += cnt
				if (v>>j)&1 != 0 {
					panic("I'm a sucker")
				}
				if !yield(i*64 + j) {
					return
				}
				u >>= cnt + 1
				j++
			}
		}
	}
}

type bitMatrix []bitArray

func makeBitMatrix(n, m int) bitMatrix {
	mx := make([]bitArray, n)
	for i := 0; i < n; i++ {
		mx[i] = makeBitArray(m)
	}
	return mx
}

func (mx bitMatrix) isSet(i, j int) bool {
	return mx[i].isSet(j)
}

func (mx bitMatrix) set(i, j int) {
	mx[i].set(j)
}

func solve(n, k int, points points) (bool, [3]int) {
	xy := makeBitMatrix(n, n)
	yx := makeBitMatrix(n, n)

	xz := makeBitMatrix(n, n)
	zx := makeBitMatrix(n, n)

	yz := makeBitMatrix(n, n)
	zy := makeBitMatrix(n, n)

	var xyN, xzN, yzN int

	for _, p := range points {
		x, y, z := p[0]-1, p[1]-1, p[2]-1 // to 0-indexing
		if !xy.isSet(x, y) {
			xyN++
			xy.set(x, y)
			yx.set(y, x)
		}
		if !xz.isSet(x, z) {
			xzN++
			xz.set(x, z)
			zx.set(z, x)
		}
		if !yz.isSet(y, z) {
			yzN++
			yz.set(y, z)
			zy.set(z, y)
		}
	}

	if debugEnable {
		log.Println("xy:", xyN)
		// for _, row := range xy {
		// 	log.Println(row)
		// }
		log.Println("xz:", xzN)
		// for _, row := range xz {
		// 	log.Println(row)
		// }
		log.Println("yz:", yzN)
		// for _, row := range yz {
		// 	log.Println(row)
		// }
	}

	// TODO: убрать копипаст
	if xyN >= xzN && xyN >= yzN {
		if debugEnable {
			log.Println("use xy")
		}
		for x := 0; x < n; x++ {
			for y := range xy[x].zeros() {
				for z := range xz[x].orAndZeros(yz[y]) {
					return false, [3]int{x + 1, y + 1, z + 1} // to 1-indexing
				}
			}
		}
	} else if xzN >= xyN && xzN >= yzN {
		if debugEnable {
			log.Println("use xz")
		}
		for x := 0; x < n; x++ {
			for z := range xz[x].zeros() {
				for y := range xy[x].orAndZeros(zy[z]) {
					return false, [3]int{x + 1, y + 1, z + 1} // to 1-indexing
				}
			}
		}
	} else { // yzN >= xyN && yzN >= xzN
		if debugEnable {
			log.Println("use yz")
		}
		for y := 0; y < n; y++ {
			for z := range yz[y].zeros() {
				for x := range yx[y].orAndZeros(zx[z]) {
					return false, [3]int{x + 1, y + 1, z + 1} // to 1-indexing
				}
			}
		}
	}
	return true, [3]int{}
}

func run(in io.Reader, out io.Writer, solve solveFunc) {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, k, err := scanTwoInt(sc)
	if err != nil {
		panic(err)
	}

	points := func(yield func(int, [3]int) bool) {
		for i := 0; i < k; i++ {
			x, y, z, err := scanThreeInt(sc)
			if err != nil {
				panic(err)
			}
			if !yield(i, [3]int{x, y, z}) {
				break
			}
		}
	}

	ok, p := solve(n, k, points)

	if ok {
		bw.WriteString("YES\n")
	} else {
		bw.WriteString("NO\n")
		writeInts(bw, p[:], defaultWriteOpts())
	}
}

// ----------------------------------------------------------------------------

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	if _, ok := os.LookupEnv("SLOW"); ok {
		run(os.Stdin, os.Stdout, slowSolve)
	} else {
		run(os.Stdin, os.Stdout, solve)
	}
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
