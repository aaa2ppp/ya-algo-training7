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
2 2 2 1 5
4
1 1
2 3
2 5
1 5
`)},
			`1
3
5
5
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

func findMaximum(nums []int32, l, r int) int32 {
	var maximum int32 = -1
	for i := l; i <= r; i++ {
		maximum = max(maximum, nums[i])
	}
	return maximum
}

func Fuzz_solve10(f *testing.F) {
	for i := 0; i < 10; i++ {
		f.Add(int64(i))
	}
	f.Fuzz(func(t *testing.T, seed int64) {
		rand := rand.New(rand.NewSource(seed))

		n := rand.Intn(10) + 1
		nums := make([]int32, n)

		for i := 0; i < n; i++ {
			nums[i] = rand.Int31n(10) + 1
		}

		tab, pow := precalc(nums)

		for l := 0; l < len(nums); l++ {
			for r := l; r < len(nums); r++ {
				wantVal := findMaximum(nums, l, r)

				got := func() int {
					defer func() {
						if p := recover(); p != nil {
							t.Logf("nums: %v", nums)
							t.Logf("query: [%d, %d]", l, r)
							// t.Fatalf("panic: %v", p)
							panic(p)
						}
					}()
					return query(tab, pow, l, r)
				}()

				if !(l <= got && got <= r) {
					t.Logf("nums: %v", nums)
					t.Logf("query: [%d, %d]", l, r)
					t.Fatalf("got = %d out of range [%d, %d]", got, l, r)
				}
				if nums[got] != wantVal {
					t.Logf("nums: %v", nums)
					t.Logf("query: [%d, %d]", l, r)
					t.Fatalf("nums[got] = %d, want %d", nums[got], wantVal)
				}
			}
		}
	})
}

func Fuzz_solve100(f *testing.F) {
	for i := 0; i < 10; i++ {
		f.Add(int64(i))
	}
	f.Fuzz(func(t *testing.T, seed int64) {
		rand := rand.New(rand.NewSource(seed))

		n := rand.Intn(100) + 1
		nums := make([]int32, n)

		for i := 0; i < n; i++ {
			nums[i] = rand.Int31n(100) + 1
		}

		tab, pow := precalc(nums)

		for l := 0; l < len(nums); l++ {
			for r := l; r < len(nums); r++ {
				wantVal := findMaximum(nums, l, r)

				got := func() int {
					defer func() {
						if p := recover(); p != nil {
							t.Logf("nums: %v", nums)
							t.Logf("query: [%d, %d]", l, r)
							// t.Fatalf("panic: %v", p)
							panic(p)
						}
					}()
					return query(tab, pow, l, r)
				}()

				if !(l <= got && got <= r) {
					t.Logf("nums: %v", nums)
					t.Logf("query: [%d, %d]", l, r)
					t.Fatalf("got = %d out of range [%d, %d]", got, l, r)
				}
				if nums[got] != wantVal {
					t.Logf("nums: %v", nums)
					t.Logf("query: [%d, %d]", l, r)
					t.Fatalf("nums[got] = %d, want %d", nums[got], wantVal)
				}
			}
		}
	})
}

func Fuzz_solve1000(f *testing.F) {
	for i := 0; i < 10; i++ {
		f.Add(int64(i))
	}
	f.Fuzz(func(t *testing.T, seed int64) {
		rand := rand.New(rand.NewSource(seed))

		n := rand.Intn(1000) + 1
		nums := make([]int32, n)

		for i := 0; i < n; i++ {
			nums[i] = rand.Int31n(1000) + 1
		}

		tab, pow := precalc(nums)

		for l := 0; l < len(nums); l++ {
			for r := l; r < len(nums); r++ {
				wantVal := findMaximum(nums, l, r)

				got := func() int {
					defer func() {
						if p := recover(); p != nil {
							t.Logf("nums: %v", nums)
							t.Logf("query: [%d, %d]", l, r)
							// t.Fatalf("panic: %v", p)
							panic(p)
						}
					}()
					return query(tab, pow, l, r)
				}()

				if !(l <= got && got <= r) {
					t.Logf("nums: %v", nums)
					t.Logf("query: [%d, %d]", l, r)
					t.Fatalf("got = %d out of range [%d, %d]", got, l, r)
				}
				if nums[got] != wantVal {
					t.Logf("nums: %v", nums)
					t.Logf("query: [%d, %d]", l, r)
					t.Fatalf("nums[got] = %d, want %d", nums[got], wantVal)
				}
			}
		}
	})
}

var bench_ans int

func Benchmark_solve(b *testing.B) {
	rand := rand.New(rand.NewSource(1))

	n := 100000
	k := 3000000

	nums := make([]int32, n)
	for i := 0; i < n; i++ {
		nums[i] = rand.Int31n(100000) + 1
	}

	for i := 0; i < b.N; i++ {
		tab, pow := precalc(nums)

		for i := 0; i < k; i++ {
			l := rand.Intn(len(nums))
			r := rand.Intn(len(nums))
			if l > r {
				l, r = r, l
			}
			bench_ans = query(tab, pow, l, r)
		}
	}

}
