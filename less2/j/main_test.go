package main

import (
	"bytes"
	"io"
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

func hash(nums []int32) uint64 {
	h := uint64(1)
	for _, v := range nums {
		h *= multiplier
		h += uint64(v)
	}
	return h
}

func Test_precalc(t *testing.T) {

	// Проверяем, что precalc правильно считает общую хеш сумму.
	// NOTE: длины последовательностей должны быть 2^n-1, т.к. дерево отрезков растягивает
	//  массив до степени двойки. Еще один элемент нам нуже для лидирующей еденицы.

	tests := []struct {
		name string
		nums []int32
	}{
		{
			"1 2 3",
			[]int32{1, 2, 3},
		},
		{
			"1 2 3 4 5",
			[]int32{1, 2, 3, 4, 5, 6, 7},
		},
		{
			"rand 7",
			[]int32{6, 6, 3, 3, 9, 4, 6},
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
				t.Errorf("got = %v, want %v", got, want)
			}
		})
	}
}
