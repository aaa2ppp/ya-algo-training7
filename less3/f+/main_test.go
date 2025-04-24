package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func Test_run(t *testing.T) {
	test_run(t, solve)
}

func Test_run_solw(t *testing.T) {
	test_run(t, slowSolve)
}

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
			args{strings.NewReader(`2 2
1 1 1
2 2 2
`)},
			`YES`,
			true,
		},
		{
			"2",
			args{strings.NewReader(`2 2
1 1 1
1 1 2
`)},
			`NO
2 2 1`,
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

func Benchmark_run(b *testing.B) {

	b.Run("run solve14", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			func() {
				f, err := os.Open("./test_data/14")
				if err != nil {
					panic(err)
				}
				defer f.Close()
				run(f, io.Discard, solve)
			}()
		}
	})

	b.Run("run solve21", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			func() {
				f, err := os.Open("./test_data/21")
				if err != nil {
					panic(err)
				}
				defer f.Close()
				run(f, io.Discard, solve)
			}()
		}
	})

	b.Run("run slow solve14", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			func() {
				f, err := os.Open("./test_data/14")
				if err != nil {
					panic(err)
				}
				defer f.Close()
				run(f, io.Discard, slowSolve)
			}()
		}
	})

	b.Run("run slow solve21", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			func() {
				f, err := os.Open("./test_data/21")
				if err != nil {
					panic(err)
				}
				defer f.Close()
				run(f, io.Discard, slowSolve)
			}()
		}
	})
}
