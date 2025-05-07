package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

func main() {
	var in io.Reader
	switch 0 {
	case 0:
		in = os.Stdin

	case 1:
		in = strings.NewReader(`3 3 7
1 2
2 3
3 1
ask 3 3
cut 1 2
ask 1 2
cut 1 3
ask 2 1
cut 2 3
ask 3 1
`)

	}
	run(in, os.Stdout)
}

func run(in io.Reader, out io.Writer) {
	br := bufio.NewReader(in)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	var n, m, k int
	_, err := fmt.Fscanln(br, &n, &m, &k)
	if err != nil {
		panic(err)
	}

	for range m {
		br.ReadLine()
	}

	type ask struct {
		cmd  string
		x, y int
	}

	asks := make([]ask, 0, k)
	for range k {
		var (
			cmd  string
			x, y int
		)
		if _, err := fmt.Fscan(br, &cmd, &x, &y); err != nil {
			panic(err)
		}
		asks = append(asks, ask{cmd, x, y})
	}

	slices.Reverse(asks)

	graph := make([]int, n+1)
	for i := range graph {
		graph[i] = i
	}

	getRoot := func(node int) (int, int) {
		var h int
		for graph[node] != node {
			h++
			node = graph[node]
		}
		return node, h
	}

	ans := make([]string, 0, k)

	setRoot := func(node int, root int) {
		for {
			next := graph[node]
			graph[node] = root
			if next == node {
				break
			}
			node = next
		}
	}

	for _, ask := range asks {
		rx, nx := getRoot(ask.x)
		ry, ny := getRoot(ask.y)

		if ask.cmd == "ask" {
			if rx == ry {
				ans = append(ans, "YES")
			} else {
				ans = append(ans, "NO")
			}
		} else { // "cut"
			if rx != ry {
				if nx < ny {
					ry = rx
				} else {
					rx = ry
				}

			}
		}

		setRoot(ask.x, rx)
		setRoot(ask.y, ry)
	}

	slices.Reverse(ans)

	bw.WriteString(strings.Join(ans, "\n"))
	bw.WriteByte('\n')
}
