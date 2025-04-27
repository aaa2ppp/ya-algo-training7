package main

import (
	"bytes"
	"container/heap"
	"io"
	"slices"
)

type HufNode struct {
	Let   uint16
	Cnt   int
	Left  *HufNode
	Right *HufNode
}

func NewHufLeaf(c uint16, cnt int) *HufNode {
	return &HufNode{
		Let: c,
		Cnt: cnt,
	}
}

func NewHufNode(a, b *HufNode) *HufNode {
	return &HufNode{
		Cnt:   a.Cnt + b.Cnt,
		Left:  a,
		Right: b,
	}
}

func (node *HufNode) IsLeaf() bool {
	return node.Left == nil && node.Right == nil
}

type HufQueue []*HufNode

func (q HufQueue) Len() int           { return len(q) }
func (q HufQueue) Less(i, j int) bool { return q[i].Cnt < q[j].Cnt }
func (q HufQueue) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }

func (q *HufQueue) Push(x any) {
	*q = append(*q, x.(*HufNode))
}

func (q *HufQueue) Pop() any {
	old := *q
	n := len(old)
	x := old[n-1]
	old[n-1] = nil
	*q = old[:n-1]
	return x
}

func BuildHufTree(text []uint16) *HufNode {
	freq := make(map[uint16]int, 256)
	for _, let := range text {
		freq[let]++
	}
	var q HufQueue
	for let, cnt := range freq {
		if cnt > 0 {
			heap.Push(&q, NewHufLeaf(uint16(let), cnt))
		}
	}
	for q.Len() > 1 {
		a := heap.Pop(&q).(*HufNode)
		b := heap.Pop(&q).(*HufNode)
		heap.Push(&q, NewHufNode(a, b))
	}
	return heap.Pop(&q).(*HufNode)
}

type HufDict map[uint16][]byte

func NewHufDict(tree *HufNode) HufDict {
	dict := make(map[uint16][]byte, 256)

	var dfs func(node *HufNode, bits []byte)
	dfs = func(node *HufNode, bits []byte) {
		if node.IsLeaf() {
			dict[node.Let] = slices.Clone(bits)
		}
		if node.Left != nil {
			dfs(node.Left, append(bits, 0))
		}
		if node.Right != nil {
			dfs(node.Right, append(bits, 1))
		}
	}
	dfs(tree, nil)
	return dict
}

func HufEncode(input []uint16, dict HufDict) []byte {
	var output bytes.Buffer
	w := NewBitWriter(&output)

	for _, c := range input {
		for _, bit := range dict[c] {
			w.WriteBit(uint(bit))
		}
	}

	return output.Bytes()
}

func HufDecode(input []byte, tree *HufNode) []uint16 {
	var output []uint16

	r := NewBitReader(bytes.NewReader(input))

	node := tree
	for {
		c, err := r.ReadBit()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		if c == 0 {
			node = node.Left
		} else {
			node = node.Right
		}
		if node.IsLeaf() {
			output = append(output, node.Let)
			node = tree
		}
	}

	return output
}
