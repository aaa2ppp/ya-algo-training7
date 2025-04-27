package main

import (
	"bytes"
	"cmp"
	"container/heap"
	"io"
	"log"
	"os"
	"slices"
	"strings"
	"unsafe"
)

type Node struct {
	Let   byte
	Cnt   int
	Left  *Node
	Right *Node
}

func NewLeaf(c byte, cnt int) *Node {
	return &Node{
		Let: c,
		Cnt: cnt,
	}
}

func NewNode(a, b *Node) *Node {
	return &Node{
		Cnt:   a.Cnt + b.Cnt,
		Left:  a,
		Right: b,
	}
}

func (node *Node) IsLeaf() bool {
	return node.Left == nil && node.Right == nil
}

type Queue []*Node

func (q Queue) Len() int           { return len(q) }
func (q Queue) Less(i, j int) bool { return q[i].Cnt < q[j].Cnt }
func (q Queue) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }

func (q *Queue) Push(x any) {
	*q = append(*q, x.(*Node))
}

func (q *Queue) Pop() any {
	old := *q
	n := len(old)
	x := old[n-1]
	old[n-1] = nil
	*q = old[:n-1]
	return x
}

func BuildTree(text []byte) *Node {
	freq := make([]int, 256)
	for _, let := range text {
		freq[let]++
	}
	var q Queue
	for let, cnt := range freq {
		if cnt > 0 {
			heap.Push(&q, NewLeaf(byte(let), cnt))
		}
	}
	for q.Len() > 1 {
		a := heap.Pop(&q).(*Node)
		b := heap.Pop(&q).(*Node)
		heap.Push(&q, NewNode(a, b))
	}
	return heap.Pop(&q).(*Node)
}

type Dict [][]byte

func NewDict(tree *Node) Dict {
	dict := make([][]byte, 256)

	var dfs func(node *Node, bits []byte)
	dfs = func(node *Node, bits []byte) {
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

func encode(input []byte, dict Dict) []byte {
	var (
		output  []byte
		bitCnt  int
		outByte byte
	)

	output = append(output, 0) // количество значащих битов последнего байта

	for _, c := range input {
		for _, bit := range dict[c] {
			outByte |= (bit << (bitCnt & 7))
			bitCnt++
			if bitCnt&7 == 0 {
				output = append(output, outByte)
				outByte = 0
			}
		}
	}

	if bitCnt&7 != 0 {
		output = append(output, outByte)
	}
	output[0] = byte(bitCnt & 7)

	log.Println("len:", len(output), "bitCnt:", bitCnt)

	return output
}

func decode(input []byte, tree *Node) []byte {
	var output []byte
	node := tree
	bitCnt := (len(input) - 1) * 8
	if input[0] != 0 {
		bitCnt -= 8
		bitCnt += int(input[0])
	}
	log.Println("len:", len(input), "0=", input[0], "bitCnt:", bitCnt)

	for _, c := range input[1:] {
		for i := 0; i < 8 && bitCnt > 0; i++ {
			if c&(1<<i) == 0 {
				node = node.Left
			} else {
				node = node.Right
			}
			if node.IsLeaf() {
				output = append(output, node.Let)
				node = tree
			}
			bitCnt--
		}
	}

	log.Println("bitCnt:", bitCnt)

	return output
}

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	log.Println("input len:", len(input))

	tree := BuildTree(input)
	dict := NewDict(tree)

	type item struct {
		let  byte
		bits []byte
	}
	var items []item
	for let, bits := range dict {
		items = append(items, item{byte(let), bits})
	}
	slices.SortFunc(items, func(a, b item) int {
		return cmp.Or(len(a.bits)-len(b.bits), strings.Compare(unsafeString(a.bits), unsafeString(b.bits)))
	})
	for i := range items {
		let, bits := items[i].let, items[i].bits
		if bits != nil {
			log.Printf("'%c' %v\n", let, bits)
		}
	}

	if len(os.Args) > 1 && os.Args[1] == "encode" {
		os.Stdout.Write(encode(input, dict))
		return
	}

	// if len(os.Args) > 1 && os.Args[1] == "decode" {
	// 	os.Stdout.Write(decode(input))
	// 	return
	// }

	encoded := encode(input, dict)
	log.Println("encoded len:", len(encoded))

	decoded := decode(encoded, tree)
	log.Println("decoded len:", len(decoded))

	if bytes.Equal(input, decoded) {
		log.Println("ok")
	} else {
		log.Println("decode != input")
	}
}

func unsafeString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
