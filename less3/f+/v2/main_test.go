package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testData = "../test_data"

func Test_run(t *testing.T) {
	test_run(t, solve)
}

func Test_run_solw(t *testing.T) {
	test_run(t, slowSolve)
}

type args struct {
	in io.Reader
}

type test struct {
	name    string
	args    args
	wantOut string
	debug   bool
}

func test_run(t *testing.T, solve solveFunc) {
	tests := []test{
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
			test_run_tt(t, solve, tt)
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

func test_run_tt(t *testing.T, solve solveFunc, tt test) {
	defer func(v bool) { debugEnable = v }(debugEnable)
	debugEnable = tt.debug
	out := &bytes.Buffer{}
	run(tt.args.in, out, solve)
	if gotOut := out.String(); trimLines(gotOut) != trimLines(tt.wantOut) {
		t.Errorf("run() = %v, want %v", gotOut, tt.wantOut)
	}
}

func Test_run_solve_data(t *testing.T) {
	list := []string{"14", "19", "21", "24"}
	for _, num := range list {
		name := "solve" + num
		t.Run(name, func(t *testing.T) {
			wantOut, err := os.ReadFile(filepath.Join(testData, num+".a"))
			if err != nil {
				t.Fatal(err)
			}

			in, err := os.Open(filepath.Join(testData, num))
			if err != nil {
				t.Fatal(err)
			}
			defer in.Close()

			test_run_tt(t, solve, test{
				"solve" + num,
				args{in},
				unsafeString(wantOut),
				false,
			})
		})
	}
}

func Benchmark_run(b *testing.B) {
	list := []string{"14", "19", "21", "24"}

	for _, num := range list {
		name := "solve" + num
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				func() {
					f, err := os.Open(filepath.Join(testData, num))
					if err != nil {
						panic(err)
					}
					defer f.Close()
					run(f, io.Discard, solve)
				}()
			}
		})
	}

	for _, num := range list {
		name := "slow solve" + num
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				func() {
					f, err := os.Open(filepath.Join(testData, num))
					if err != nil {
						panic(err)
					}
					defer f.Close()
					run(f, io.Discard, slowSolve)
				}()
			}
		})
	}
}
