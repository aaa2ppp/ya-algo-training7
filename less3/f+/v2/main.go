package main

import (
	"bufio"
	"io"
	"iter"
	"log"
	"math"
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
		arr[len(arr)-1] = math.MaxUint64 << k
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

			// ищем нулевые биты
			for j := 0; j < 64; j++ {
				if v&1 == 0 {
					if !yield(i*64 + j) {
						return
					}
				}
				v >>= 1
			}
		}
	}
}

// orAndZeros делает побитовое OR с другим массивом. Возвращает индесы первого нуля,
// в результирующем массиве.
func (arr bitArray) orAndFindZero(other bitArray) int {
	if len(arr) != len(other) {
		panic("length arr and other must be equal")
	}
	for i := range arr {
		v := arr[i] | other[i]
		if v == math.MaxUint64 {
			// все биты единички
			continue
		}

		// ищем первый нулевой бит
		for j := 0; j < 64; j++ {
			if v&1 == 0 {
				return i*64 + j
			}
			v >>= 1
		}
	}
	return -1
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
		// 	log.Printf("%064b", row)
		// }
		log.Println("xz:", xzN)
		// for _, row := range xz {
		// 	log.Printf("%064b", row)
		// }
		log.Println("yz:", yzN)
		// for _, row := range yz {
		// 	log.Printf("%064b", row)
		// }
	}

	// TODO: убрать копипаст
	if xyN >= xzN && xyN >= yzN {
		if debugEnable {
			log.Println("use xy")
		}
		for x := range n {
			for y := range xy[x].zeros() {
				z := xz[x].orAndFindZero(yz[y])
				if z != -1 {
					return false, [3]int{x + 1, y + 1, z + 1} // to 1-indexing
				}
			}
		}
	} else if xzN >= xyN && xzN >= yzN {
		if debugEnable {
			log.Println("use xz")
		}
		for x := range n {
			for z := range xz[x].zeros() {
				y := xy[x].orAndFindZero(zy[z])
				if y != -1 {
					return false, [3]int{x + 1, y + 1, z + 1} // to 1-indexing
				}
			}
		}
	} else { // yzN >= xyN && yzN >= xzN
		if debugEnable {
			log.Println("use yz")
		}
		for y := range n {
			for z := range yz[y].zeros() {
				x := yx[y].orAndFindZero(zx[z])
				if x != -1 {
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
		writeInts(bw, p[:])
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

func unsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
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

func _parseInt[T Int](b []byte) (T, error) {
	if ^T(0) < 0 {
		v, err := strconv.ParseInt(unsafeString(b), 0, int(unsafe.Sizeof(T(1)))<<3)
		return T(v), err
	} else {
		v, err := strconv.ParseUint(unsafeString(b), 0, int(unsafe.Sizeof(T(1)))<<3)
		return T(v), err
	}
}

func scanIntX[T Int](sc *bufio.Scanner) (T, error) {
	if sc.Scan() {
		return _parseInt[T](sc.Bytes())
	}
	if err := sc.Err(); err != nil {
		return 0, err
	}
	return 0, io.EOF
}

func scanInts[X Int](sc *bufio.Scanner, buf []X) (_ []X, err error) {
	for n := 0; n < len(buf); n++ {
		buf[n], err = scanIntX[X](sc)
		if err != nil {
			return buf[:n], err
		}
	}
	return buf, nil
}

func scanTwoIntX[T Int](sc *bufio.Scanner) (T, T, error) {
	var buf [2]T
	_, err := scanInts(sc, buf[:])
	return buf[0], buf[1], err
}

func scanThreeIntX[T Int](sc *bufio.Scanner) (T, T, T, error) {
	var buf [3]T
	_, err := scanInts(sc, buf[:])
	return buf[0], buf[1], buf[2], err
}

func scanFourIntX[T Int](sc *bufio.Scanner) (T, T, T, T, error) {
	var buf [4]T
	_, err := scanInts(sc, buf[:])
	return buf[0], buf[1], buf[2], buf[3], err
}

func scanFiveIntX[T Int](sc *bufio.Scanner) (T, T, T, T, T, error) {
	var buf [5]T
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

func _writeInt[T Int](bw *bufio.Writer, v T) (int, error) {
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

func writeInt[T Int](bw *bufio.Writer, v T, opts ...writeOpts) error {
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

func writeInts[T Int](bw *bufio.Writer, a []T, opts ...writeOpts) error {
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
