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
4 5 5 2 3
`)},
			`3 0 0 1 2 
`,
			true,
		},
		{
			"2",
			args{strings.NewReader(`5
5 1 3 1 5
`)},
			`0 1 2 1 0 
`,
			true,
		},
		{
			"3",
			args{strings.NewReader(`3
6 6 6
`)},
			`0 0 0 
`,
			true,
		},
		{
			"4",
			args{strings.NewReader(`4
6 5 5 6
`)},
			`0 0 0 0 
`,
			true,
		},
		{
			"5",
			args{strings.NewReader(`2
4 5
`)},
			`0 0 
`,
			true,
		},
		{
			"8",
			args{strings.NewReader(`100
591 417 888 251 792 847 685 3 182 461 102 348 555 956 771 901 712 878 580 631 342 333 285 899 525 725 537 718 929 653 84 788 104 355 624 803 253 853 201 995 536 184 65 205 540 652 549 777 248 405 677 950 431 580 600 846 328 429 134 983 526 103 500 963 400 23 276 704 570 757 410 658 507 620 984 244 486 454 802 411 985 303 635 283 96 597 855 775 139 839 839 61 219 986 776 72 729 69 20 917
`)},
			`2 1 7 1 2 6 5 1 2 3 1 2 4 9 1 7 1 5 1 4 3 2 1 6 1 3 1 2 8 2 1 4 1 2 3 5 1 6 1 0 4 2 1 3 5 6 1 7 1 2 3 8 1 2 3 4 1 2 1 9 3 1 2 6 3 1 2 4 1 5 1 3 1 2 10 1 2 1 3 1 0 1 4 2 1 3 0 2 1 0 0 1 2 0 4 1 3 2 1 8 
`,
			true,
		},

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
