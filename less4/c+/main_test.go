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
			args{strings.NewReader(`push_back 1
back
exit
`)},
			`ok
1
bye
`,
			true,
		},
		{
			"2",
			args{strings.NewReader(`size
push_back 1
size
push_back 2
size
push_front 3
size
exit
`)},
			`0
ok
1
ok
2
ok
3
bye
`,
			true,
		},
		{
			"3",
			args{strings.NewReader(`push_back 3
push_front 14
size
clear
push_front 1
back
push_back 2
front
pop_back
size
pop_front
size
exit
`)},
			`ok
ok
2
ok
ok
1
ok
1
2
1
1
0
bye
`,
			true,
		},
		{
			"6",
			args{strings.NewReader(`push_back 4234
front
back
size
push_back 34234
front
back
size
push_front 4231342
front
back
size
push_back 2345
front
back
size
push_back 41234
front
back
size
push_front 423412
front
back
size
pop_back
front
back
size
pop_back
front
back
size
pop_front
front
back
size
pop_back
front
back
size
pop_front
front
back
size
pop_front
size
exit

`)},
			`ok
4234
4234
1
ok
4234
34234
2
ok
4231342
34234
3
ok
4231342
2345
4
ok
4231342
41234
5
ok
423412
41234
6
41234
423412
2345
5
2345
423412
34234
4
423412
4231342
34234
3
34234
4231342
4234
2
4231342
4234
4234
1
4234
0
bye
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
