package main

import (
	"bytes"
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
			args{strings.NewReader(`1
0100010000111101
`)},
			`010010010100001111101`,
			true,
		},
		{
			"2",
			args{strings.NewReader(`2
010010010100001111101`)},
			`0100010000111101`,
			true,
		},
		{
			"3",
			args{strings.NewReader(`2
010011010100001111101`)}, // 6 bit is broken
			`0100010000111101`,
			true,
		},
		{
			"4",
			args{strings.NewReader(`2
010010010000001111101`)}, // 10 bit is broken
			`0100010000111101`,
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

func Test_encode(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		debug bool
	}{
		{
			"1",
			args{"0100010000111101"},
			"010010010100001111101",
			true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func(v bool) { debugEnable = v }(debugEnable)
			debugEnable = tt.debug
			if got := encode(tt.args.line); got != tt.want {
				t.Errorf("encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decode(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		debug bool
	}{
		{
			"1",
			args{"010010010100001111101"},
			"0100010000111101",
			true,
		},
		{
			"2",
			//    010010010100001111101
			//          v
			args{"010011010100001111101"},
			"0100010000111101",
			true,
		},
		{
			"2",
			//    010010010100001111101
			//              v
			args{"010010010000001111101"},
			"0100010000111101",
			true,
		},
		{
			"100000",
			//    010010010100001111101
			//              v
			args{"010010010000001111101"},
			"0100010000111101",
			true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func(v bool) { debugEnable = v }(debugEnable)
			debugEnable = tt.debug
			if got := decode(tt.args.line); got != tt.want {
				t.Errorf("decode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_encode100000(t *testing.T) {
	origin := generate(rand.New(rand.NewSource(1)), 100000)
	encoded := encode(origin)

	if len(encoded) > 100017 {
		t.Fatalf("encoded lan must be <=100017, you may have forgotten to trim the first zero.")
	}

	for _, c := range encoded {
		switch c {
		case '0', '1':
		default:
			t.Fatal("encoded must contains only 0 and")
		}
	}
}

func generate(rand *rand.Rand, n int) string {
	var sb strings.Builder
	for i := 0; i < n; {
		num := rand.Uint64()
		for j := 0; j < 64 && i < n; j++ {
			sb.WriteByte(byte(num&1) + '0')
			i++
		}
	}
	return sb.String()
}

func test_encodeDecode(t *testing.T, seed int64, maxLen int) {
	rand := rand.New(rand.NewSource(seed))
	n := rand.Intn(maxLen) + 1
	origin := generate(rand, n)
	encoded := encode(origin)

	for _, c := range encoded {
		switch c {
		case '0', '1':
		default:
			t.Log("origin :", origin)
			t.Log("encoded:", encoded)
			t.Fatal("encoded must contains only 0 and")
		}
	}

	if decoded := decode(encoded); decoded != origin {
		for _, c := range encoded {
			switch c {
			case '0', '1':
			default:
				t.Log("origin :", origin)
				t.Log("encoded:", encoded)
				t.Fatal("decoded must contains only 0 and")
			}
		}

		t.Log("origin :", origin)
		t.Log("encoded:", encoded)
		t.Log("decoded:", decoded)
		t.Fatal("fail decode")
	}

	broken := []byte(encoded)
	bit := rand.Intn(len(broken))
	broken[bit] ^= 1
	if unsafeString(broken) == encoded {
		t.Fatal("oops!.. no broken")
	}

	if decoded := decode(unsafeString(broken)); decoded != origin {
		for _, c := range encoded {
			switch c {
			case '0', '1':
			default:
				t.Log("origin :", origin)
				t.Log("encoded:", encoded)
				t.Fatal("decoded must contains only 0 and")
			}
		}

		t.Log("origin :", origin)
		t.Log("encoded:", encoded)
		t.Log("brokBit:", bit)
		t.Log("broken :", unsafeString(broken))
		t.Fatal("fail restore")
	}
}

func Fuzz_encodeDecode10(f *testing.F) {
	const maxLen = 10
	for i := 0; i < 10; i++ {
		f.Add(int64(i))
	}
	f.Fuzz(func(t *testing.T, seed int64) {
		test_encodeDecode(t, seed, maxLen)
	})
}

func Fuzz_encodeDecode100(f *testing.F) {
	const maxLen = 100
	for i := 0; i < 10; i++ {
		f.Add(int64(i))
	}
	f.Fuzz(func(t *testing.T, seed int64) {
		test_encodeDecode(t, seed, maxLen)
	})
}

func Fuzz_encodeDecode1000(f *testing.F) {
	const maxLen = 1000
	for i := 0; i < 10; i++ {
		f.Add(int64(i))
	}
	f.Fuzz(func(t *testing.T, seed int64) {
		test_encodeDecode(t, seed, maxLen)
	})
}

func Fuzz_encodeDecode100000(f *testing.F) {
	const maxLen = 1000000
	for i := 0; i < 10; i++ {
		f.Add(int64(i))
	}
	f.Fuzz(func(t *testing.T, seed int64) {
		test_encodeDecode(t, seed, maxLen)
	})
}

var bench_out string

func Benchmark_encode(b *testing.B) {
	origin := generate(rand.New(rand.NewSource(1)), 100000)
	for i := 0; i < b.N; i++ {
		bench_out = encode(origin)
	}
}

func Benchmark_decode(b *testing.B) {
	origin := generate(rand.New(rand.NewSource(1)), 100000)
	encoded := encode(origin)
	for i := 0; i < b.N; i++ {
		bench_out = decode(encoded)
	}
}
