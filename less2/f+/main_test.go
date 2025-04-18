package main

import (
	"bytes"
	"io"
	"slices"
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
			args{strings.NewReader(`4 5
1 2 3 4
1 1 1
1 1 3
1 1 5
0 2 3
1 1 3
`)},
			`1
3
-1
2
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

func findFirstNotLess(nums []int32, i int, v int32) int {
	for ; i < len(nums); i++ {
		if nums[i] >= v {
			return i
		}
	}
	return -1
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
		nums := make([]int32, n)
		for i := 0; i < n; i++ {
			nums[i] = rand.Int31n(maxV)
		}

		tree := precalc(nums)

		al := rand.Intn(len(nums))
		ar := rand.Intn(len(nums))
		if al > ar {
			al, ar = ar, al
		}
		ar++ // to open

		numsCopy := slices.Clone(nums)

		ui := rand.Intn(len(numsCopy))
		uv := rand.Int31n(maxV)
		numsCopy[ui] = uv
		update(tree, ui, uv)

		for v := int32(0); v < maxV; v++ {
			for i := 0; i < len(nums); i++ {
				want := findFirstNotLess(numsCopy, i, v)
				got := query(tree, i, v)
				if got != want {
					t.Logf("nums: %v", nums)
					t.Logf("set: %d, %d", ui, uv)
					t.Logf("nums: %v", numsCopy)
					t.Logf("get: %d, %d", i, v)
					t.Fatalf("got = %d, want %d", got, want)
				}
			}
		}

	})
}

// func Benchmark_solve(b *testing.B) {
// 	buf, err := os.ReadFile("test_data/9")
// 	if err != nil {
// 		panic(err)
// 	}

// 	for i := 0; i < b.N; i++ {
// 		run(bytes.NewReader(buf), io.Discard)
// 	}
// }
