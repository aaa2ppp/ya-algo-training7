package main

import (
	"bytes"
	"fmt"
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
1 2 1 2 1
4
1 2 4 2
0 3 5 2
1 1 3 2
1 2 3 3
`)},
			`+-+
`,
			true,
		},
		{
			"2",
			args{strings.NewReader(`1
1
3
1 1 1 1
0 1 1 2
1 1 1 1`)},
			`++`,
			true,
		},
		{
			"3",
			args{strings.NewReader(`5
1 2 3 1 2  
3
1 1 4 2
0 1 5 6
1 1 2 4`)},
			`++`,
			true,
		},
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

func hash(nums []int32) uint64 {
	h := uint64(0)
	for _, v := range nums {
		h *= multiplier
		h += uint64(v)
	}
	return h
}

func Test_precalc(t *testing.T) {

	// Проверяем, что precalc правильно считает общую хеш сумму.

	tests := []struct {
		name string
		nums []int32
	}{
		{
			"1",
			[]int32{1},
		},
		{
			"1 2",
			[]int32{1, 2},
		},
		{
			"1 2 3",
			[]int32{1, 2, 3},
		},
		{
			"1 2 3 4",
			[]int32{1, 2, 3, 4},
		},
		{
			"1 2 3 4 5",
			[]int32{1, 2, 3, 4, 5},
		},
		{
			"1 2 3 4 5 6",
			[]int32{1, 2, 3, 4, 5, 6},
		},
		{
			"1 2 3 4 5 6 7",
			[]int32{1, 2, 3, 4, 5, 6, 7},
		},
		{
			"1 2 3 4 5 6 7 8",
			[]int32{1, 2, 3, 4, 5, 6, 7, 8},
		},
		{
			"rand 7",
			[]int32{6, 6, 3, 3, 9, 4, 6},
		},
		{
			"rand 10",
			[]int32{85, 96, 94, 12, 14, 47, 28, 59, 70, 7},
		},
		{
			"rand 15",
			[]int32{85, 96, 94, 12, 14, 47, 28, 59, 70, 7, 32, 5, 79, 34, 11},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := hash(tt.nums)
			tree := precalc(tt.nums)
			got := tree[0].val
			if got != want {
				t.Errorf("got1 = %v, want1 %v", got, want)
			}
		})
	}
}

func Test_query(t *testing.T) {

	// Проверяем, что query правильно считает хеш сумму на отрезке.

	tests := []struct {
		name string
		nums []int32
	}{
		{
			"1 2 3 1 2",
			[]int32{1, 2, 3, 1, 2},
		},
		{
			"1",
			[]int32{1},
		},
		{
			"1 2",
			[]int32{1, 2},
		},
		{
			"1 2 3",
			[]int32{1, 2, 3},
		},
		{
			"1 2 3 4",
			[]int32{1, 2, 3, 4},
		},
		{
			"1 2 3 4 5",
			[]int32{1, 2, 3, 4, 5},
		},
		{
			"1 2 3 4 5 6",
			[]int32{1, 2, 3, 4, 5, 6},
		},
		{
			"1 2 3 4 5 6 7",
			[]int32{1, 2, 3, 4, 5, 6, 7},
		},
		{
			"1 2 3 4 5 6 7 8",
			[]int32{1, 2, 3, 4, 5, 6, 7, 8},
		},
		{
			"rand 7",
			[]int32{6, 6, 3, 3, 9, 4, 6},
		},
		{
			"rand 10",
			[]int32{85, 96, 94, 12, 14, 47, 28, 59, 70, 7},
		},
		{
			"rand 15",
			[]int32{85, 96, 94, 12, 14, 47, 28, 59, 70, 7, 32, 5, 79, 34, 11},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := precalc(tt.nums)
			for l := 0; l < len(tt.nums); l++ {
				for r := l + 1; r <= len(tt.nums); r++ {
					want := hash(tt.nums[l:r])
					got := query(tree, l, r)
					if got != want {
						t.Logf("q [%d, %d]", l, r)
						t.Errorf("got1 = %v, want1 %v", got, want)
					}
				}
			}
		})
	}
}

func Test_update(t *testing.T) {

	// простой тест на обновление

	defer func(v bool) { debugEnable = v }(debugEnable)
	// debugEnable = true

	nums := []int32{1, 2, 3, 4, 5}
	tree := precalc(nums)

	copy(nums[0:3], []int32{6, 6, 6})
	update(tree, 0, 3, 6)

	for l := 0; l < len(nums); l++ {
		for r := l + 1; r <= len(nums); r++ {
			want := hash(nums[l:r])
			got := query(tree, l, r)
			if got != want {
				t.Logf("tree: %v", tree)
				t.Logf("q [%d, %d)", l, r)
				t.Errorf("got = %d, want %d", got, want)
			}
		}
	}
}

func updateNums(nums []int32, l, r int, val int32) {
	for i := l; i < r; i++ {
		nums[i] = val
	}
}

func test_solve(t *testing.T, seed int64, maxN int, maxVal int32) {
	rand := rand.New(rand.NewSource(seed))

	buf := &bytes.Buffer{}

	n := rand.Intn(maxN) + 1
	fmt.Fprintln(buf, n)

	nums := make([]int32, n)
	for i := 0; i < n; i++ {
		nums[i] = rand.Int31n(maxVal) + 1
	}
	fmt.Fprintln(buf, nums)

	tree := precalc(nums)

	for l := 0; l < len(nums); l++ {
		for r := l + 1; r <= len(nums); r++ {
			want := hash(nums[l:r])
			got := query(tree, l, r)
			if got != want {
				fmt.Fprintf(buf, "q [%d %d)\n", l, r)
				t.Log(buf.String())
				t.Fatalf("got1 = %v, want1 %v", got, want)
			}
		}
	}

	for i := 0; i < 10; i++ {
		l, r := rand.Intn(n), rand.Intn(n)
		if l > r {
			l, r = r, l
		}
		r++
		val := rand.Int31n(maxVal) + 1
		updateNums(nums, l, r, val)
		update(tree, l, r, val)
		fmt.Fprintf(buf, "u [%d %d) %d\n", l, r, val)
	}

	for l := 0; l < len(nums); l++ {
		for r := l + 1; r <= len(nums); r++ {
			want := hash(nums[l:r])
			got := query(tree, l, r)
			if got != want {
				fmt.Fprintf(buf, "q [%d %d)\n", l, r)
				t.Log(buf.String())
				t.Fatalf("got1 = %v, want1 %v", got, want)
			}
		}
	}
}

func Fuzz_solve10(f *testing.F) {
	const (
		maxN   = 10
		maxVal = 10
	)

	for i := 1; i < 10; i++ {
		f.Add(int64(i))
	}

	f.Fuzz(func(t *testing.T, a int64) {
		test_solve(t, a, maxN, maxVal)
	})
}

func Fuzz_solve100(f *testing.F) {
	const (
		maxN   = 100
		maxVal = 100
	)

	for i := 1; i < 10; i++ {
		f.Add(int64(i))
	}

	f.Fuzz(func(t *testing.T, a int64) {
		test_solve(t, a, maxN, maxVal)
	})
}

func Benchmark_solve(b *testing.B) {
	rand := rand.New(rand.NewSource(1))

	const (
		N      = 100000
		Q      = 100000
		maxVal = 100000
	)

	nums := make([]int32, N)
	for i := 0; i < N; i++ {
		nums[i] = rand.Int31n(maxVal) + 1
	}

	for i := 0; i < b.N; i++ {
		tree := precalc(nums)
		for j := 0; j < Q; j++ {
			{
				l, r := rand.Intn(N), rand.Intn(N)
				if l > r {
					l, r = r, l
				}
				r++
				val := rand.Int31n(maxVal) + 1
				update(tree, l, r, val)
			}

			{
				l, r := rand.Intn(N), rand.Intn(N)
				if l > r {
					l, r = r, l
				}
				r++
				query(tree, l, r)
			}
		}
	}

}
