package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	var in io.Reader
	switch 0 {
	case 0:
		in = os.Stdin

	case 1:
		in = strings.NewReader(`8
0 1
1 5
2 4
3 2
4 3
5 0
6 6
1 0
`)

	}
	run(in, os.Stdout)
}

func run(in io.Reader, out io.Writer) {
	br := bufio.NewReader(in)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	var n int
	_, err := fmt.Fscanln(br, &n)
	if err != nil {
		panic(err)
	}

	type snowman struct {
		prev   int
		weight int
	}

	snowmans := make([]snowman, 1, n)
	var totalWeight int

	for range n {
		var t, m int
		if _, err := fmt.Fscanln(br, &t, &m); err != nil {
			panic(err)
		}
		if m == 0 {
			t = snowmans[t].prev
			snowmans = append(snowmans, snowmans[t])
			totalWeight += snowmans[t].weight

		} else {

			weight := snowmans[t].weight + m
			snowmans = append(snowmans, snowman{t, weight})
			totalWeight += weight
		}
	}

	fmt.Fprintln(bw, totalWeight)
}
