package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	in := os.Stdin
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
	graph := make([][]int, n+1)
	for i := 1; i < len(graph); i++ {
		var j int
		_, err := fmt.Fscan(br, &j)
		if err != nil {
			panic(err)
		}
		graph[i] = append(graph[i], j)
		graph[j] = append(graph[j], i)
	}
	visited := make([]bool, n+1)
	var dfs func(node int)
	dfs = func(node int) {
		if visited[node] {
			return
		}
		visited[node] = true
		for _, neig := range graph[node] {
			dfs(neig)
		}
	}
	var count int
	for node := 1; node < len(graph); node++ {
		if !visited[node] {
			count++
			dfs(node)
		}
	}
	fmt.Fprintln(bw, count)
}
