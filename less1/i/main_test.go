package main

import (
	"bytes"
	"io"
	"reflect"
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
			args{strings.NewReader(`3 7
4 1 2
3 1 2
2 1 2
`)},
			`3 3
1 2 3 
`,
			true,
		},
		{
			"2",
			args{strings.NewReader(`3 7
4 1 3
3 1 2
2 1 1
`)},
			`2 2
2 3 
`,
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

			// if gotOut := out.String(); trimLines(gotOut) != trimLines(tt.wantOut) {
			// 	t.Errorf("run() = %v, want %v", gotOut, tt.wantOut)
			// }

			// XXX: Now check only first line. To full check use Test_solve
			if gotOut := out.String(); lines(trimLines(gotOut))[0] != lines(tt.wantOut)[0] {
				t.Errorf("run() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func lines(s string) []string {
	s = strings.TrimRight(s, "\n")
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t\r")
	}
	return lines
}

func trimLines(text string) string {
	lines := lines(text)
	for n := len(lines); n > 0 && lines[n-1] == ""; n-- {
		lines = lines[:n-1]
	}
	return strings.Join(lines, "\n")
}

func Test_solve(t *testing.T) {
	type args struct {
		maxVolume int
		items     []item
	}
	tests := []struct {
		name      string
		args      args
		wantCost  int
		wantItems []int
		debug     bool
	}{
		{
			"1",
			args{
				7,
				[]item{
					{1, 4, 1, 2},
					{2, 3, 1, 2},
					{3, 2, 1, 2},
				},
			},
			3, []int{1, 2, 3},
			true,
		},
		{
			"2",
			args{
				7,
				[]item{
					{1, 4, 1, 3},
					{2, 3, 1, 2},
					{3, 2, 1, 1},
				},
			},
			2, []int{2, 3},
			true,
		},
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func(v bool) { debugEnable = v }(debugEnable)
			debugEnable = tt.debug

			gotCost, gotItems := solve(tt.args.maxVolume, slices.Clone(tt.args.items))
			if gotCost != tt.wantCost {
				t.Errorf("solve() cost = %v, want %v", gotCost, tt.wantCost)
			}

			checkItems := func(gotItems []int, wantCost int) bool {
				totalCost := 0
				volume := 0
				for _, num := range gotItems {
					it := tt.args.items[num-1]
					totalCost += it.cost
					volume += it.volume
				}

				if totalCost != wantCost {
					t.Errorf("total cost of items = %v, want %v", totalCost, wantCost)
					return false
				}

				maxPressure := max(0, volume-tt.args.maxVolume)
				for _, num := range gotItems {
					it := tt.args.items[num-1]
					if it.pressure < maxPressure {
						t.Errorf("pressure of item #%v, max posible pressure %v", it.pressure, maxPressure)
						return false
					}
				}

				return true
			}

			if !reflect.DeepEqual(gotItems, tt.wantItems) && !checkItems(gotItems, tt.wantCost) {
				t.Errorf("solve() gotItems = %v, wantItems %v", gotItems, tt.wantItems)
			}
		})
	}
}
