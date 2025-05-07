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
		in = strings.NewReader(`4 5
1 2
1 3
2 3
3 4
4 1
`)
	case 2:
		in = strings.NewReader(`50 110
34 17
32 43
22 12
22 37
19 6
33 10
40 20
7 24
46 31
28 2
50 40
42 5
29 26
45 45
10 36
17 2
22 49
6 22
12 40
15 44
3 48
22 8
6 6
45 13
29 14
33 7
49 31
5 49
42 1
35 47
13 6
50 3
8 36
49 23
13 50
41 3
9 43
33 45
48 12
39 9
6 11
27 44
5 17
35 34
18 27
1 19
40 21
27 7
13 43
30 2
6 11
33 2
17 24
21 14
18 22
46 7
43 16
22 32
14 39
29 4
23 50
29 3
45 21
3 2
22 17
27 18
25 10
22 17
3 30
27 32
49 11
42 14
20 5
38 42
22 6
28 43
14 27
39 17
11 43
3 39
37 37
28 29
6 4
38 12
32 7
25 42
50 1
29 50
35 35
6 46
33 48
24 22
40 19
35 4
3 19
28 34
26 13
30 46
43 21
6 40
10 9
24 13
12 7
30 50
3 45
16 16
1 14
5 49
40 46
8 35
`)
	}
	run(in, os.Stdout)
}

func run(in io.Reader, out io.Writer) {
	br := bufio.NewReader(in)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	var n, m int
	_, err := fmt.Fscanln(br, &n, &m)
	if err != nil {
		panic(err)
	}

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

	var count int
	for i := 0; i < m; i++ {
		var x, y int
		if _, err := fmt.Fscanln(br, &x, &y); err != nil {
			panic(err)
		}
		count++

		rx, nx := getRoot(x)
		ry, ny := getRoot(y)
		if rx != ry {
			n--
			if n == 1 {
				break
			}

			if nx < ny {
				ry = rx
			} else {
				rx = ry
			}
		}

		setRoot(x, rx)
		setRoot(y, ry)
	}

	fmt.Fprintln(bw, count)
}
