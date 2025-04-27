package main

import (
	"bytes"
	"container/heap"
	"io"
	"log"
	"os"
	"unsafe"
)

const maxCount = (1 << 14)

func encode(input []byte) []byte {
	var output []byte
	var outByte byte
	var bitCnt int

	writeBit := func(v byte) {
		outByte |= (v & 1) << (bitCnt & 7)
		bitCnt++
		if bitCnt&7 == 0 {
			output = append(output, outByte)
			outByte = 0
		}
	}

	writeBits := func(v int, n int) {
		for i := 0; i < n; i++ {
			writeBit(byte(v))
			v >>= 1
		}
	}

	items := make([]Item, 0)
	index := make(map[string]*Item)
	queue := PriorityQueue{}

	items = append(items, Item{id: 0, word: ""})
	index[""] = &items[0]
	heap.Push(&queue, &items[0])

	l, r, p := 0, 0, 0

	prev := 0
	for r < len(input) {
		w := unsafeString(input[l : r+1])
		cur, ok := index[w]
		if ok {
			prev = cur.id
			r++
			continue
		}

		if len(items) > (1 << p) {
			p++
		}
		// log.Printf("%d (%d) '%c'", prev, p, input[r])
		writeBits(prev, p)
		writeBits(int(input[r])-'a'+1, 5)
		// log.Println("bitCnt:", bitCnt)

		items[prev].priority = r
		heap.Fix(&queue, items[prev].index)

		prev = 0
		r++
		l = r

		if len(items) < maxCount {
			items = append(items, Item{id: len(items), word: w})
			index[w] = &items[len(items)-1]
			heap.Push(&queue, &items[len(items)-1])
		} else {
			item := queue[0]
			item.word = w
			item.priority = 0
			heap.Fix(&queue, item.index)
		}
	}

	log.Println("words count:", len(items), "p:", p)

	if prev != 0 {
		// log.Printf("%d (%d)", prev, p)
		writeBits(prev, p)
		// log.Println("bitCnt:", bitCnt)
	}

	if bitCnt&7 != 0 {
		output = append(output, outByte)
	}

	return output // todo
}

func decode(input []byte) []byte {
	var bitCnt int
	readBit := func() int {
		if bitCnt>>3 >= len(input) {
			return -1
		}
		x := (int(input[bitCnt>>3]) >> (bitCnt & 7)) & 1
		bitCnt++
		return x
	}

	readBits := func(n int) int {
		var x int
		for i := 0; i < n; i++ {
			bit := readBit()
			if bit == -1 {
				return -1
			}
			x |= bit << i
		}
		return x
	}

	items := make([]Item, 0)
	items = append(items, Item{id: 0, word: ""})
	queue := PriorityQueue{}
	heap.Push(&queue, &items[0])

	var output bytes.Buffer
	p := 0
	for {
		if len(items) > (1 << p) {
			p++
		}

		prev := readBits(p)
		if prev < 0 {
			break
		}
		w := items[prev].word
		output.WriteString(w)

		c := readBits(5)
		if c <= 0 {
			break
		}
		c += 'a' - 1
		output.WriteByte(byte(c))

		items[prev].priority = output.Len()
		heap.Fix(&queue, items[prev].index)

		if len(items) == maxCount {
			continue
		}

		w += string(byte(c))
		if len(items) < maxCount {
			items = append(items, Item{id: len(items), word: w})
			heap.Push(&queue, &items[len(items)-1])
		} else {
			item := queue[0]
			item.word = w
			item.priority = 0
			heap.Fix(&queue, item.index)
		}
	}

	return output.Bytes()
}

func main() {

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 && os.Args[1] == "encode" {
		os.Stdout.Write(encode(input))
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "decode" {
		os.Stdout.Write(decode(input))
		return
	}

	log.Println("input len:", len(input))

	encoded := encode(input)
	log.Printf("encoded len: %v %0.2f%%", len(encoded),
		float64(len(encoded))/float64(len(input))*100)

	decoded := decode(encoded)
	log.Println("decoded len:", len(decoded))

	if bytes.Equal(input, decoded) {
		log.Println("ok")
	} else {
		log.Println("decoded != input")
	}

	// if _, err := os.Stdout.Write(decoded); err != nil {
	// 	log.Fatal(err)
	// }
}

func unsafeString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

type Item struct {
	id       int
	word     string
	priority int
	index    int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[j].id == 0 || pq[i].id != 0 && (pq[i].priority < pq[j].priority ||
		pq[i].priority == pq[j].priority && pq[i].word < pq[j].word)
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}
