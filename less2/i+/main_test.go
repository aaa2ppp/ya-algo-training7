package main

import (
	"bytes"
	"io"
	"math"
	"math/rand"
	"os"
	"slices"
	"strings"
	"testing"
)

func Test_run(t *testing.T) {
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
			args{strings.NewReader(`5
2 4 3 1 5
5
m 1 3
a 2 4 100
m 1 3
a 5 5 10
m 1 5
`)},
			`4 104 104 
`,
			true,
		},
		{
			"2",
			args{strings.NewReader(`7
1 2 2 2 0 0 2
2
a 3 6 1
m 5 5
`)},
			`1`,
			true,
		},
		// {
		// 	"3",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	true,
		// },
		// {
		// 	"4",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	true,
		// },
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func(v bool) { debugEnable = v }(debugEnable)
			debugEnable = tt.debug
			out := &bytes.Buffer{}
			run(tt.args.in, out)
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

func updateRange(nums []int64, l, r int, add int64) {
	for i := l; i < r; i++ {
		nums[i] += add
	}
}

func findMaximum(nums []int64, l, r int) int64 {
	maximum := int64(math.MinInt64)
	for i := l; i < r; i++ {
		maximum = max(maximum, nums[i])
	}
	return maximum
}

func Fuzz_solve10(f *testing.F) {
	for i := 0; i < 10; i++ {
		f.Add(int64(i))
	}
	const (
		maxN = 10
		maxV = 10
	)
	f.Fuzz(func(t *testing.T, a int64) {
		rand := rand.New(rand.NewSource(a))
		n := rand.Intn(maxN) + 1
		nums := make([]int64, n)
		for i := 0; i < n; i++ {
			nums[i] = rand.Int63n(maxV)
		}

		tree := precalc(nums)

		al := rand.Intn(len(nums))
		ar := rand.Intn(len(nums))
		if al > ar {
			al, ar = ar, al
		}
		ar++ // to open

		add := rand.Int63n(maxV)
		numsCopy := slices.Clone(nums)
		updateRange(numsCopy, al, ar, add)
		update(tree, al, ar, add)

		for ql := 0; ql < len(nums); ql++ {
			for qr := ql + 1; qr <= len(numsCopy); qr++ { // r is open
				want := findMaximum(numsCopy, ql, qr)
				got := query(tree, ql, qr)
				if got != want {
					t.Logf("nums: %v", nums)
					t.Logf("a: [%d, %d) %d", al, ar, add)
					t.Logf("nums: %v", numsCopy)
					t.Logf("m: [%d, %d)", ql, qr)
					t.Fatalf("got = %d, want %d", got, want)
				}
			}
		}

	})
}

func Benchmark_solve(b *testing.B) {
	buf, err := os.ReadFile("test_data/02")
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		run(bytes.NewReader(buf), io.Discard)
	}
}
