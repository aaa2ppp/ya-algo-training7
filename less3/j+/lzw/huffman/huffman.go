package huffman

import (
	"bytes"
	"cmp"
	"container/heap"
	"log"
	"os"
	"slices"
	"unsafe"

	"ya-training7/less3/j+/lzw/bitio"
)

var _, debugEnable = os.LookupEnv("DEBUG")

type Char interface {
	~uint | ~uint8 | ~uint16 | ~uint32
}

type Node[C Char] struct {
	left  *Node[C]
	right *Node[C]
	cnt   int
	char  C
}

func newLeaf[C Char](c C, cnt int) *Node[C] {
	return &Node[C]{
		char: c,
		cnt:  cnt,
	}
}

func newNode[C Char](a, b *Node[C]) *Node[C] {
	return &Node[C]{
		cnt:   a.cnt + b.cnt,
		left:  a,
		right: b,
	}
}

func (node *Node[C]) isLeaf() bool {
	return node.left == nil && node.right == nil
}

type queue[C Char] []*Node[C]

func (q queue[C]) Len() int           { return len(q) }
func (q queue[C]) Less(i, j int) bool { return q[i].cnt < q[j].cnt }
func (q queue[C]) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }

func (q *queue[C]) Push(x any) {
	*q = append(*q, x.(*Node[C]))
}

func (q *queue[C]) Pop() any {
	old := *q
	n := len(old)
	x := old[n-1]
	old[n-1] = nil
	*q = old[:n-1]
	return x
}

func buildTree[C Char](text []C) *Node[C] {
	// TODO: len(text) == 0

	freq := make(map[C]int, 256)
	for _, code := range text {
		freq[code]++
	}

	// TODO: len(freq) == 1

	var q queue[C]
	for code, cnt := range freq {
		if cnt > 0 {
			heap.Push(&q, newLeaf(code, cnt))
		}
	}

	for q.Len() > 1 {
		a := heap.Pop(&q).(*Node[C])
		b := heap.Pop(&q).(*Node[C])
		heap.Push(&q, newNode(a, b))
	}

	return heap.Pop(&q).(*Node[C])
}

type Dict[C Char] map[C][]byte

func dictFromTree[C Char](tree *Node[C]) Dict[C] {
	dict := make(map[C][]byte, 256)

	var dfs func(node *Node[C], bits []byte)
	dfs = func(node *Node[C], bits []byte) {
		if node.isLeaf() {
			if debugEnable {
				log.Printf("%d: %v", node.char, bits)
			}
			dict[node.char] = slices.Clone(bits)
			return
		}
		if node.left != nil {
			dfs(node.left, append(bits, 0))
		}
		if node.right != nil {
			dfs(node.right, append(bits, 1))
		}
	}

	dfs(tree, nil)
	return dict
}

func (dict Dict[C]) normalize() {
	if debugEnable {
		log.Println("== Dict.normalize")
	}

	type item struct {
		char C
		len  int
	}

	items := make([]item, 0, len(dict))
	for code, bits := range dict {
		items = append(items, item{code, len(bits)})
	}

	slices.SortFunc(items, func(a, b item) int {
		return cmp.Or(a.len-b.len, int(a.char)-int(b.char))
	})

	var (
		code   uint
		curLen int = 1
	)

	for _, item := range items {
		code <<= item.len - curLen
		curLen = item.len
		// if debugEnable {
		// 	log.Printf("-> %0"+strconv.Itoa(curLen)+"b", code)
		// }

		bits := dict[item.char]
		for i, n := 0, len(bits); i < n; i++ {
			bits[i] = byte((code >> (n - i - 1)) & 1)
		}

		if debugEnable {
			if debugEnable {
				log.Printf("%d: %v", item.char, bits)
			}
		}

		code++
	}
}

// writeTo function should only be called for a normalized dictionary!
func (dict Dict[C]) writeTo(w *bitio.Writer) error {
	var (
		maxChar  C
		charSize   = int(unsafe.Sizeof(maxChar)) * 8
		minChar  C = C(1)<<C(charSize) - 1
		maxLen   uint8
		lenSize  = int(unsafe.Sizeof(maxLen)) * 8
	)

	for char, bits := range dict {
		minChar = min(minChar, char)
		maxChar = max(maxChar, char)
		maxLen = max(maxLen, uint8(len(bits)))
	}

	if debugEnable {
		log.Printf("minChar:%d(%d) maxChar:%d(%d) maxLen:%d(%d)", minChar, charSize, maxChar, charSize, maxLen, lenSize)
	}

	if err := w.WriteBits(uint(minChar), charSize); err != nil {
		return err
	}

	if err := w.WriteBits(uint(maxChar), charSize); err != nil {
		return err
	}

	if err := w.WriteBits(uint(maxLen), lenSize); err != nil {
		return err
	}

	for char := minChar; char <= maxChar; char++ {
		n := len(dict[char])
		if debugEnable {
			log.Printf("write %d:%b(%d)", char, n, n)
		}
		if err := w.WriteBits(uint(n), int(maxLen)); err != nil {
			return err
		}
	}
	return nil
}

func readDictFrom[C Char](r *bitio.Reader) (Dict[C], error) {
	var (
		minChar  C
		maxChar  C
		charSize = int(unsafe.Sizeof(maxChar)) * 8
		maxLen   uint8
		lenSize  = int(unsafe.Sizeof(maxLen)) * 8
	)

	{
		v, err := r.ReadBits(charSize)
		if err != nil {
			return nil, err
		}
		minChar = C(v)
	}

	{
		v, err := r.ReadBits(charSize)
		if err != nil {
			return nil, err
		}
		maxChar = C(v)
	}

	{
		v, err := r.ReadBits(lenSize)
		if err != nil {
			return nil, err
		}
		maxLen = uint8(v)
	}

	if debugEnable {
		log.Printf("minChar:%d maxChar:%d maxLen:%d", minChar, maxChar, maxLen)
	}

	type item struct {
		char C
		len  int
	}

	items := make([]item, 0, maxChar-minChar+1)
	for char := minChar; char <= maxChar; char++ {
		n, err := r.ReadBits(int(maxLen))
		if err != nil {
			return nil, err
		}
		if debugEnable {
			log.Printf("read %d:%b(%d)", char, n, n)
		}
		if n != 0 {
			items = append(items, item{char, int(n)})
		}
	}

	slices.SortFunc(items, func(a, b item) int {
		return cmp.Or(a.len-b.len, int(a.char)-int(b.char))
	})

	dict := make(Dict[C], int(maxChar)+1)

	var (
		code   uint
		curLen int = 1
	)

	for _, item := range items {
		code <<= item.len - curLen
		curLen = item.len

		bits := make([]byte, curLen)
		for i, n := 0, len(bits); i < n; i++ {
			bits[i] = byte((code >> (n - i - 1)) & 1)
		}

		dict[item.char] = bits

		if debugEnable {
			if debugEnable {
				log.Printf("%d: %v", item.char, bits)
			}
		}

		code++
	}

	return dict, nil
}

func treeFromDict[C Char](dict Dict[C]) *Node[C] {
	tree := &Node[C]{}

	for char, bits := range dict {
		node := tree
		for _, bit := range bits {
			switch bit {
			case 0:
				if node.left == nil {
					node.left = &Node[C]{}
				}
				node = node.left
			case 1:
				if node.right == nil {
					node.right = &Node[C]{}
				}
				node = node.right
			}
			node.char = char
		}
	}

	return tree
}

func BuildDict[C Char](text []C) Dict[C] {
	dict := dictFromTree(buildTree(text))
	dict.normalize()
	return dict
}

func Encode[C Char](input []C) []byte {
	if debugEnable {
		log.Println("== Encode")
	}

	var output bytes.Buffer
	w := bitio.NewWriter(&output)

	// Резервируем 3 бита под кол-во значащих бит в последнем байте
	w.WriteBits(0, 3)

	dict := dictFromTree(buildTree(input))
	dict.normalize()
	dict.writeTo(w)

	for _, code := range input {
		for _, bit := range dict[code] {
			w.WriteBit(uint(bit))
		}
	}

	w.Close()

	// Запоминаем кол-во значащих бит в последнем байте. Если 0 - все биты значащие
	if n := byte(w.Cnt() & 7); n != 0 {
		if debugEnable {
			log.Printf("cnt: %d, last byte: %d", w.Cnt(), n)
		}
		output.Bytes()[0] |= n
	}

	return output.Bytes()
}

func Decode[C Char](input []byte) ([]C, error) {
	if debugEnable {
		log.Println("== Decode")
	}

	var output []C
	r := bitio.NewReader(bytes.NewReader(input))

	// Вычисляем количество бит с учетом кол-ва значащих бит в последнем байте
	cnt := len(input) * 8
	if n, _ := r.ReadBits(3); n != 0 {
		cnt -= (8 - int(n))
	}
	if debugEnable {
		log.Printf("cnt: %d", cnt)
	}

	dict, err := readDictFrom[C](r)
	if err != nil {
		return nil, err
	}

	tree := treeFromDict(dict)
	node := tree
	for i := r.Cnt(); i < cnt; i++ {
		c, err := r.ReadBit()
		if err != nil {
			return output, err
		}
		if c == 0 {
			node = node.left
		} else {
			node = node.right
		}
		if node.isLeaf() {
			output = append(output, node.char)
			node = tree
		}
	}

	return output, nil
}
