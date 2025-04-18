package main

import (
	"bytes"
	"io"
	"math/rand"
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
328 0 0 0 0
5
QUERY 1 3
UPDATE 2 832
QUERY 3 3
QUERY 2 3
UPDATE 2 0
`)},
			`2
1
1
`,
			true,
		},
		// {
		// 	"2",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	true,
		// },
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

func countСonsecutiveZeros(nums []int32, l, r int) int {
	var maximum, cnt int
	for i := l; i < r; i++ { // r is open
		if nums[i] == 0 {
			cnt++
			maximum = max(maximum, cnt)
		} else {
			cnt = 0
		}
	}
	return maximum
}

func test_solve(t *testing.T, nums []int32, tree []item, l, r int) {
	want := countСonsecutiveZeros(nums, l, r)

	got := func() int {
		defer func() {
			if p := recover(); p != nil {
				t.Logf("nums: %v", nums)
				t.Logf("query: [%d, %d)", l, r)
				// t.Fatalf("panic: %v", p)
				panic(p)
			}
		}()
		return query(tree, l, r)
	}()

	if got != want {
		t.Logf("nums: %v", nums)
		t.Logf("query: %d %d %d", l, r, 2)
		t.Fatalf("got = %d, want %d", got, want)
	}
}

func Fuzz_solve10(f *testing.F) {
	const (
		maxN   = 10
		maxVal = 10
	)
	for i := 0; i < 10; i++ {
		f.Add(int64(i))
	}
	f.Fuzz(func(t *testing.T, seed int64) {
		rand := rand.New(rand.NewSource(seed))

		n := rand.Intn(maxN) + 1
		nums := make([]int32, n)

		for i := 0; i < n; i++ {
			nums[i] = rand.Int31n(maxVal*3/2) + 1
			if nums[i] > maxVal {
				nums[i] = 0
			}
		}

		tree := precalc(nums)

		for i := range nums {
			for l := 0; l < len(nums); l++ {
				for r := l + 1; r <= len(nums); r++ { // r is open
					test_solve(t, nums, tree, l, r)
				}
			}

			nums[i] = rand.Int31n(15) + 1
			nums[i] = rand.Int31n(maxVal*3/2) + 1
			if nums[i] > maxVal {
				nums[i] = 0
			}
			update(tree, i, nums[i])
		}
	})
}

func Fuzz_solve100(f *testing.F) {
	const (
		maxN   = 100
		maxVal = 100
	)
	for i := 0; i < 10; i++ {
		f.Add(int64(i))
	}
	f.Fuzz(func(t *testing.T, seed int64) {
		rand := rand.New(rand.NewSource(seed))

		n := rand.Intn(maxN) + 1
		nums := make([]int32, n)

		for i := 0; i < n; i++ {
			nums[i] = rand.Int31n(maxVal*3/2) + 1
			if nums[i] > maxVal {
				nums[i] = 0
			}
		}

		tree := precalc(nums)

		for i := range nums {
			for l := 0; l < len(nums); l++ {
				for r := l + 1; r <= len(nums); r++ { // r is open
					test_solve(t, nums, tree, l, r)
				}
			}

			nums[i] = rand.Int31n(15) + 1
			nums[i] = rand.Int31n(maxVal*3/2) + 1
			if nums[i] > maxVal {
				nums[i] = 0
			}
			update(tree, i, nums[i])
		}
	})
}

var bench_ans int

func Benchmark_solve(b *testing.B) {
	rand := rand.New(rand.NewSource(1))

	const (
		N      = 500000
		M      = 50000
		maxVal = 1000
	)

	nums := make([]int32, N)
	for i := 0; i < N; i++ {
		nums[i] = rand.Int31n(maxVal) + 1
		if nums[i] > 100 {
			nums[i] = 0
		}
	}

	for i := 0; i < b.N; i++ {
		tree := precalc(nums)

		for i := 0; i < M; i++ {
			l := rand.Intn(len(nums))
			r := rand.Intn(len(nums))
			if l > r {
				l, r = r, l
			}
			r++
			bench_ans = query(tree, l, r)

			i := rand.Intn(len(nums))
			nums[i] = rand.Int31n(150) + 1
			if nums[i] > 100 {
				nums[i] = 0
			}
			update(tree, i, nums[i])
		}
	}
}
