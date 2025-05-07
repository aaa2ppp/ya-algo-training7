package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Node[T any] struct {
	prev  *Node[T]
	next  *Node[T]
	Value T
}

func (node Node[T]) Next() *Node[T] {
	return node.next
}

func (node Node[T]) Prev() *Node[T] {
	return node.prev
}

func NewNode[T any](v T) *Node[T] {
	return &Node[T]{Value: v}
}

func (node *Node[T]) Remove() {
	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}
	node.next = nil
	node.prev = nil
}

type Round[T any] struct {
	head *Node[T]
}

func (r *Round[T]) Head() *Node[T] {
	return r.head
}

func (r *Round[T]) Insert(node *Node[T]) {
	if r.head == node {
		return
	}

	node.Remove()
	if r.head == nil {
		node.next = node
		node.prev = node
	} else {
		next := r.head
		prev := next.prev

		node.next = next
		next.prev = node

		node.prev = prev
		prev.next = node
	}

	r.head = node
}

func run(in io.Reader, out io.Writer) {
	br := bufio.NewReader(in)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	var n int
	if _, err := fmt.Fscanln(br, &n); err != nil {
		panic(err)
	}

	var round Round[string]

	for i := 0; i < n; i++ {
		s, err := br.ReadString('\n')
		if err != nil {
			panic(err)
		}

		switch {
		case strings.HasPrefix(s, "Run"):
			taskName := strings.TrimSpace(strings.TrimPrefix(s, "Run"))
			round.Insert(NewNode(taskName))

		case strings.HasPrefix(s, "Alt+"):
			node := round.Head()
			for i, n := 0, strings.Count(s, "+"); i < n; i++ {
				node = node.Next()
			}
			round.Insert(node)

		default:
			panic("unknown prefix: " + strings.TrimSpace(s))
		}

		bw.WriteString(round.Head().Value)
		bw.WriteByte('\n')
	}
}

// ----------------------------------------------------------------------------

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	run(os.Stdin, os.Stdout)
}
