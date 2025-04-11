package main

import (
	"bytes"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"testing"
)

func test_run(t *testing.T, solve solveFunc) {
	type args struct {
		in io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		debug   bool
	}{
		{
			"1",
			args{strings.NewReader(`4
DSD
SS
DD
SDD
`)},
			`3
`,
			true,
		},
		{
			"2",
			args{strings.NewReader(`5
DDDD
SD
SSS
DD
SSSDS`)},
			`5`,
			true,
		},
		{
			"3",
			args{strings.NewReader(`5
SSDS
SD
DSDSD
S
DSD`)},
			`6`,
			true,
		},
		{
			"4",
			args{strings.NewReader(`5
SD
SD
D
DSSSD
DS`)},
			`5`,
			true,
		},
		{
			"5",
			args{strings.NewReader(`7
SS
SSD
DD
D
SD
SSSDD
DDS`)},
			`6`,
			true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func(v bool) { debugEnable = v }(debugEnable)
			debugEnable = tt.debug
			out := &bytes.Buffer{}
			run(tt.args.in, out, solve)
			if gotOut := out.String(); trimLines(gotOut) != trimLines(tt.wantOut) {
				t.Errorf("run() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func trimLines(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t\r\n")
	}
	for n := len(lines); n > 0 && lines[n-1] == ""; n-- {
		lines = lines[:n-1]
	}
	return strings.Join(lines, "\n")
}

func Test_run_solve(t *testing.T) {
	test_run(t, solve)
}

func Test_run_bruteforce(t *testing.T) {
	test_run(t, bruteforce)
}

const maxOrderLen = 5

func generateOrders(seed int64, n int, mm ...int) []string {
	m := maxOrderLen
	if len(mm) > 0 {
		m = mm[0]
	}

	orders := make([]string, n)
	rand := rand.New(rand.NewSource(seed))
	var sb strings.Builder
	for i := range orders {
		m := rand.Intn(m) + 1
		sb.Reset()
		sb.Grow(m)
		x := rand.Int()
		for i := 0; i < m; i++ {
			if x&1 == 1 {
				sb.WriteByte('S')
			} else {
				sb.WriteByte('D')
			}
			x >>= 1
		}
		orders[i] = sb.String()
	}
	return orders
}

func Fuzz_solve2(f *testing.F) {
	for i := 1; i <= 10; i++ {
		f.Add(int64(i))
	}
	f.Fuzz(func(t *testing.T, seed int64) {
		orders := generateOrders(seed, 2)
		got := solve(orders)
		want := bruteforce(orders)
		if got != want {
			t.Log(strconv.Itoa(len(orders)) + "\n" + strings.Join(orders, "\n"))
			t.Errorf("solve()=%d, want %d", got, want)
		}
	})
}

func Fuzz_solve3(f *testing.F) {
	for i := 1; i <= 10; i++ {
		f.Add(int64(i))
	}
	f.Fuzz(func(t *testing.T, seed int64) {
		orders := generateOrders(seed, 3)
		got := solve(orders)
		want := bruteforce(orders)
		if got != want {
			t.Log(strconv.Itoa(len(orders)) + "\n" + strings.Join(orders, "\n"))
			t.Errorf("solve()=%d, want %d", got, want)
		}
	})
}

func Fuzz_solve5(f *testing.F) {
	for i := 1; i <= 10; i++ {
		f.Add(int64(i))
	}
	f.Fuzz(func(t *testing.T, seed int64) {
		orders := generateOrders(seed, 5)
		got := solve(orders)
		want := bruteforce(orders)
		if got != want {
			t.Log(strconv.Itoa(len(orders)) + "\n" + strings.Join(orders, "\n"))
			t.Errorf("solve()=%d, want %d", got, want)
		}
	})
}

func Fuzz_solve7(f *testing.F) {
	for i := 1; i <= 10; i++ {
		f.Add(int64(i))
	}
	f.Fuzz(func(t *testing.T, seed int64) {
		orders := generateOrders(seed, 7)
		got := solve(orders)
		want := bruteforce(orders)
		if got != want {
			t.Log(strconv.Itoa(len(orders)) + "\n" + strings.Join(orders, "\n"))
			t.Errorf("solve()=%d, want %d", got, want)
		}
	})
}

func Fuzz_solve10(f *testing.F) {
	for i := 1; i <= 10; i++ {
		f.Add(int64(i))
	}
	f.Fuzz(func(t *testing.T, seed int64) {
		orders := generateOrders(seed, 10)
		got := solve(orders)
		want := bruteforce(orders)
		if got != want {
			t.Log(strconv.Itoa(len(orders)) + "\n" + strings.Join(orders, "\n"))
			t.Errorf("solve()=%d, want %d", got, want)
		}
	})
}

func Fuzz_bruteforce5(f *testing.F) {
	for i := 1; i <= 10; i++ {
		f.Add(int64(i))
	}
	f.Fuzz(func(t *testing.T, seed int64) {
		orders := generateOrders(seed, 5)
		v1 := bruteforce(orders)
		v2 := bruteforce2(orders)
		if v1 != v2 {
			t.Log(strconv.Itoa(len(orders)) + "\n" + strings.Join(orders, "\n"))
			t.Errorf("bruteforce()=%d, bruteforce2()=%d", v1, v1)
		}
	})
}

func Fuzz_bruteforce10(f *testing.F) {
	for i := 1; i <= 10; i++ {
		f.Add(int64(i))
	}
	f.Fuzz(func(t *testing.T, seed int64) {
		orders := generateOrders(seed, 10)
		v1 := bruteforce(orders)
		v2 := bruteforce2(orders)
		if v1 != v2 {
			t.Log(strconv.Itoa(len(orders)) + "\n" + strings.Join(orders, "\n"))
			t.Errorf("bruteforce()=%d, bruteforce2()=%d", v1, v1)
		}
	})
}

func Benchmark_solve(b *testing.B) {
	orders := generateOrders(1, 100000, 100)
	for i := 0; i < b.N; i++ {
		solve(orders)
	}
}
