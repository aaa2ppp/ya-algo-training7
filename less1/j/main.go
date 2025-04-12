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

// Вася стал аскетом и решил отказаться от всего материального ради духовного.
// Действительно, перспектива переродиться в желтого земляного червяка из-за
// утреннего капучино может здорово напугать.

// Однако отказываться от материального оказалось не так-то просто. Для каждого
// события своей ежедневной материальной жизни Вася определил его материальность
// и обозначил её целым положительным числом mi. Свою духовную силу Вася определил
// как число D. Каждый день он может отказываться от одного события материальной
// жизни и возвращать некоторые события материальной жизни, от которых он
// отказался ранее, чтобы суммарное количество материального снизилось не более
// чем на D. При этом нельзя делать так, чтобы количество материальности в
// какой-либо день выросло — это собьёт Васю с пути аскезы.

// Вася разработал оптимальный план для себя и стал гуру. Теперь его ученики
// отказывались от материального в пользу Васи. Учеников оказалось очень много,
// Вася успешно определяет их события материального мира и духовную силу, но
// теперь ему нужна программа, которая будет разрабатывать план отказа от
// материального. В оптимальном плане нужно отказаться от максимального количества
// событий материальной жизни, а в случае, если это возможно сделать несколькими
// способами, нужно сделать это за наименьшее количество дней.

type event struct {
	name   string
	weight int
}

func solve(maxWeightDif int, events []event) (totalDays int, eventNames []string) {
	slices.SortFunc(events, func(a, b event) int {
		return a.weight - b.weight
	})

	maxWeight := 0
	for i := range events {
		maxWeight += events[i].weight
	}

	// будем записывать в dp минимальное число дней которые необходимы,
	// чтобы отказаться от суммарного веса событий
	dp := make([]int, maxWeight+1)
	for i := range dp {
		dp[i] = -1
	}
	dp[0] = 0

	topWeight := 0
	for _, event := range events {
		minDays := math.MaxInt // чтобы отказаться от возвращенных событий

		for j := max(0, event.weight-maxWeightDif); j <= topWeight; j++ {
			if cnt := dp[j]; cnt != -1 && cnt < minDays {
				minDays = dp[j]
			}
		}

		if debugEnable {
			log.Println(event, "minDays:", minDays)
		}

		if minDays == math.MaxInt {
			// oops!.. событие слишком весомое, чтобы от него отказаться
			break
		}

		minDays++ // отказ от текущего события
		eventNames = append(eventNames, event.name)
		totalDays += minDays

		topWeight = min(topWeight+event.weight, maxWeight)
		for j := topWeight; j >= event.weight; j-- {
			if cnt := dp[j-event.weight]; cnt != -1 {
				if dp[j] == -1 || dp[j] > cnt+minDays {
					dp[j] = cnt + minDays
				}
			}
		}
	}

	slices.Sort(eventNames)
	return totalDays, eventNames
}

func run(in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)
	sc.Split(bufio.ScanWords)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	n, d, err := scanTwoInt(sc)
	if err != nil {
		panic(err)
	}

	events := make([]event, n)
	for i := 0; i < n; i++ {
		name, err := scanWord(sc)
		if err != nil {
			panic(err)
		}
		m, err := scanInt(sc)
		if err != nil {
			panic(err)
		}
		events[i] = event{name, m}
	}

	t, names := solve(d, events)

	writeInts(bw, []int{len(names), t}, defaultWriteOpts())
	for _, s := range names {
		bw.WriteString(s)
		bw.WriteByte('\n')
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
