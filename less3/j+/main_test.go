package main

import (
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"

	"math/rand"
)

func Test_run_io(t *testing.T) {
	pack := func(b []byte) []byte { return b }
	unpack := pack

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
			"1 pack",
			args{strings.NewReader(`pack
abacabaca
`)},
			`9
97 98 97 99 97 98 97 99 97
`,
			true,
		},
		{
			"2 unpack",
			args{strings.NewReader(`unpack
9
97 98 97 99 97 98 97 99 97
`)},
			`abacabaca`,
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
			run(tt.args.in, out, pack, unpack)
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

func Test_run(t *testing.T) {
	text, err := os.ReadFile("norm/eye_of_world.txt")
	if err != nil {
		panic(err)
	}

	rand := rand.New(rand.NewSource(1))

	for i := 0; i < 100; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			n := rand.Intn(100001) + 100000
			l := rand.Intn(len(text) - n)

			var input bytes.Buffer
			input.WriteString("pack\n")
			input.Write(text[l : l+n])
			input.WriteByte('\n')

			// t.Logf("%s", input.Bytes())

			var encoded bytes.Buffer
			encoded.WriteString("unpack\n")
			run(bytes.NewReader(input.Bytes()), &encoded, encode, decode)

			// TODO: проверить степень сжатия
			size, err := strconv.Atoi(strings.TrimSpace(
				string(bytes.Split(encoded.Bytes(), []byte("\n"))[1])))
			if err != nil {
				t.Fatal("first line must be number")
			}
			if s := float64(size) / float64(n); s > 0.5 {
				t.Errorf("too much %0.2f%%", s)
			}
			// t.Logf("%s", encoded.Bytes())

			defer func(v bool) { debugEnable = v }(debugEnable)
			debugEnable = true

			var decoded bytes.Buffer
			run(bytes.NewReader(encoded.Bytes()), &decoded, encode, decode)

			if !bytes.Equal(text[l:l+n], bytes.TrimRight(decoded.Bytes(), "\r\n")) {
				// t.Logf("%s", text[l:l+n])
				// t.Logf("%s", decoded.Bytes())
				t.Error("input and decoded not equal!")
			}

		})
	}
}
