package main

import (
	"bufio"
	"io"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"unsafe"
)

// Один из главных недостатков ровера-доставщика — ограниченный по объёму отсек,
// в котором иногда чуть-чуть не хватает места. Экспериментальная модель ровера
// обладает эластичным отсеком для перевозки заказов.

// Базовый объём отсека составляет S литров. Пока отсек не заполнен, товары в нём
// не испытывают дополнительного давления. Однако, поскольку отсек эластичный,
// в него можно положить дополнительные товары сверх базового объёма. Если объём
// положенных в отсек товаров U превышает S, то все товары в отсеке будут испытывать
// давление P=U−S.

// Каждый товар обладает тремя характеристиками: объёмом vi, стоимостью ci и
// давлением, которое он выдерживает pi.

// Всего необходимо доставить N товаров, однако в первую поездку ровера необходимо
// отправить товары с максимальной суммарной стоимостью — это обрадует заказчика.
// Определите максимальную стоимость товаров, которые можно разместить в ровере
// так, чтобы все они выдерживали давление.

// Ограничения:
// 	1 <= N <= 100
// 	0 <= S <= 10^9
// 	1 <= vi <= 1000
// 	0 <= ci <= 10^6
// 	0 <= pi <= 10^9

type item struct {
	id       int
	volume   int
	cost     int
	pressure int
}

func solve(maxVolume int, items []item) (int, []int) {
	n := len(items)

	// проверим сначала может и так все влезет
	totalVolume := 0
	maxPresscure := 0
	minPresscure := math.MaxInt
	for i := range items {
		totalVolume += items[i].volume
		minPresscure = min(minPresscure, items[i].pressure)
		maxPresscure = max(maxPresscure, items[i].pressure)
	}

	if totalVolume <= maxVolume+minPresscure {
		totalCost := 0
		allIDs := make([]int, n)
		for i := 0; i < n; i++ {
			totalCost += items[i].cost
			allIDs[i] = items[i].id
		}
		return totalCost, allIDs
	}

	// сортируем по убыванию допустимого давления
	slices.SortFunc(items, func(a, b item) int {
		return b.pressure - a.pressure
	})

	type dpItem struct {
		idx  int // prev dp.i (of sorted items)
		prev int // prev dp.j
		cost int
	}

	dp := makeMatrix[dpItem](n+1, maxVolume+maxPresscure+1)
	for i := 1; i < len(dp[0]); i++ {
		dp[0][i] = dpItem{-1, -1, -1}
	}

	var (
		maxCost       = -1
		maxCostIdx    = 0 // dp.i
		maxCostVolume = 0
	)

	topVolume := 0
	for i, it := range items {
		copy(dp[i+1], dp[i])
		dpRow := dp[i+1]

		topVolume = min(topVolume+it.volume, maxVolume+it.pressure)

		for j := topVolume; j >= it.volume; j-- {
			if prevCost := dpRow[j-it.volume].cost; prevCost != -1 {
				curCost := prevCost + it.cost
				if curCost > dpRow[j].cost {
					dpRow[j] = dpItem{
						idx:  i,
						prev: j - it.volume,
						cost: curCost,
					}
					if curCost > maxCost {
						maxCost = curCost
						maxCostIdx = i + 1
						maxCostVolume = j
					}
				}
			}
		}

		if debugEnable {
			log.Printf("%d: p=%d %v\n", i+1, it.pressure, dpRow)
		}
	}

	if debugEnable {
		log.Printf("maxCost=%d at (%d,%d)\n", maxCost, maxCostIdx, maxCostVolume)
	}

	if maxCost == -1 {
		// unpossible
		return 0, nil
	}

	var ans []int
	for i, j := maxCostIdx, maxCostVolume; j > 0; i, j = dp[i][j].idx, dp[i][j].prev {
		ans = append(ans, items[dp[i][j].idx].id)
	}

	return maxCost, ans
}

func run(in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, maxVolume, err := scanTwoInt(sc)
	if err != nil {
		panic(err)
	}

	items := make([]item, n)
	for i := 0; i < n; i++ {
		v, c, p, err := scanThreeInt(sc)
		if err != nil {
			panic(err)
		}
		items[i] = item{
			id:       i + 1,
			volume:   v,
			cost:     c,
			pressure: p,
		}
	}

	maxCost, ans := solve(maxVolume, items)
	writeInts(bw, []int{len(ans), maxCost}, defaultWriteOpts())
	if len(ans) > 0 {
		slices.Sort(ans) // на всякий случай
		writeInts(bw, ans, defaultWriteOpts())
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
