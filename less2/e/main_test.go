package main

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"math/rand"
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
0 0 3 0 2
3
u 1 5
u 1 0
s 1 5 3
`)},
			`4`,
			true,
		},
		{
			"2",
			args{strings.NewReader(`7
7 0 6 0 0 0 10
1
s 3 3 2
`)},
			`-1`,
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

func findKZero(nums []int32, l, r, k int) int {
	var cnt int
	for i := l; i <= r; i++ {
		if nums[i] == 0 {
			cnt++
		}
		if cnt == k {
			return i
		}
	}
	return -1
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
			nums[i] = rand.Int31n(15) + 1
			if nums[i] > 10 {
				nums[i] = 0
			}
		}

		tree := precalc(nums)

		for i := range nums {
			for l := 0; l < len(nums); l++ {
				for r := l; r < len(nums); r++ {
					wantVal := findKZero(nums, l, r, 2)

					got := func() int {
						defer func() {
							if p := recover(); p != nil {
								t.Logf("nums: %v", nums)
								t.Logf("query: [%d, %d] %d", l, r, 2)
								// t.Fatalf("panic: %v", p)
								panic(p)
							}
						}()
						return query(tree, l, r+1, 2)
					}()

					if got != wantVal {
						t.Logf("nums: %v", nums)
						t.Logf("query: %d %d %d", l, r, 2)
						t.Fatalf("nums[got] = %d, want %d", nums[got], wantVal)
					}
				}
			}

			nums[i] += rand.Int31n(10) + 1
			update(tree, i, nums[i])
		}
	})
}

// func Fuzz_solve100(f *testing.F) {
// 	for i := 0; i < 10; i++ {
// 		f.Add(int64(i))
// 	}
// 	f.Fuzz(func(t *testing.T, seed int64) {
// 		rand := rand.New(rand.NewSource(seed))

// 		n := rand.Intn(100) + 1
// 		nums := make([]int32, n)

// 		for i := 0; i < n; i++ {
// 			nums[i] = rand.Int31n(100) + 1
// 		}

// 		tree := precalc(nums)

// 		for i := range nums {
// 			for l := 0; l < len(nums); l++ {
// 				for r := l; r < len(nums); r++ {
// 					wantVal := findMaximum(nums, l, r)

// 					got := func() int32 {
// 						defer func() {
// 							if p := recover(); p != nil {
// 								t.Logf("nums: %v", nums)
// 								t.Logf("query: [%d, %d]", l, r)
// 								// t.Fatalf("panic: %v", p)
// 								panic(p)
// 							}
// 						}()
// 						return query(tree, l, r+1)
// 					}()

// 					if got != wantVal {
// 						t.Logf("nums: %v", nums)
// 						t.Logf("query: [%d, %d]", l, r)
// 						t.Fatalf("nums[got] = %d, want %d", nums[got], wantVal)
// 					}
// 				}
// 			}

// 			nums[i] += rand.Int31n(100) + 1
// 			update(tree, i, nums[i])
// 		}
// 	})
// }

// var bench_ans int32

// func Benchmark_solve(b *testing.B) {
// 	rand := rand.New(rand.NewSource(1))

// 	n := 100000
// 	k := 30000

// 	nums := make([]int32, n)
// 	for i := 0; i < n; i++ {
// 		nums[i] = rand.Int31n(100000)
// 	}

// 	for i := 0; i < b.N; i++ {
// 		tree := precalc(nums)

// 		for i := 0; i < k; i++ {
// 			l := rand.Intn(len(nums))
// 			r := rand.Intn(len(nums))
// 			if l > r {
// 				l, r = r, l
// 			}
// 			bench_ans = query(tree, l, r)
// 			i := rand.Intn(len(nums))
// 			v := rand.Int31n(100000)
// 			update(tree, i, v)
// 		}
// 	}

// }
