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
			args{strings.NewReader(`3
7 10 11
`)},
			`14 9 7 
`, // 13 10 7 тоже верный ответ (нужна более интеллектуальная проверка)
			true,
		},
		{
			"2",
			args{strings.NewReader(`3
7 10 3
`)},
			`impossible
`,
			true,
		},
		// 		{
		// 			"10",
		// 			args{strings.NewReader(`10
		// 9 27 79 95 52 20 76 88 57 15`)},
		// 			``,
		// 			true,
		// 		},
		// 		{
		// 			"19",
		// 			args{strings.NewReader(`50
		// 142248252569430780 133324907501390179 983947614974662124 512504946295209842 39075343779535902 912901592249295650 67898179883791095 520516499312895360 51420401458264395 574343816776192833 526020783734316114 860436503630917972 348921198005401590 507786190494489987 324703645786487230 992245742652340813 148896847560920130 599376160399393180 563085128591974750 964319185706948288 736206337881553707 702824168949630807 695653003819762771 920985796801184095 346357621001855544 952528903560082100 872866132237844240 153857886759239894 809654254270367452 349139887501274437 720210528913917212 518213634238538696 292380132335184750 759779631671019398 996461293938371677 851508323119609400 319978316143175550 567838986200920603 38320196453973021 968243359383024370 795611966917317852 808941385439456457 449217222338967809 44320396055586310 728183842408367634 375853386119954859 277865520987756775 209559833148758218 134648355746659608 617882951076076658
		// `)},
		// 			``,
		// 			true,
		// 		},
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
