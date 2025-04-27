package main

import (
	"bytes"
	"io"
	"log"
	"os"
)

const maxBits = 14

func encode(input []byte) ([]byte, *HufNode) {
	codes := LZWEncode(input, maxBits)
	tree := BuildHufTree(codes)
	dict := NewHufDict(tree)
	return HufEncode(codes, dict), tree
	// var output bytes.Buffer
	// w := NewBitWriter(&output)
	// for _, code := range codes {
	// 	w.WriteBits(uint(code), maxBits)
	// }
	// return output.Bytes()
}

func decode(input []byte, tree *HufNode) ([]byte, error) {
	// var codes []uint16
	// r := NewBitReader(bytes.NewReader(input))
	// for {
	// 	code, err := r.ReadBits(maxBits)
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	codes = append(codes, uint16(code))
	// }
	// log.Println("len codes:", len(codes))
	codes := HufDecode(input, tree)
	return LZWDecode(codes)
}

var debugEnable bool

func main() {
	_, debugEnable = os.LookupEnv("DEBUG")

	// if len(os.Args) > 1 && os.Args[1] == "encode" {
	// 	input, err := io.ReadAll(os.Stdin)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	os.Stdout.Write(encode(input))
	// 	return
	// }

	// if len(os.Args) > 1 && os.Args[1] == "decode" {
	// 	input, err := io.ReadAll(os.Stdin)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	output, err := decode(input)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	os.Stdout.Write(output)
	// 	return
	// }

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("input:", len(input), "bytes")

	encoded, tree := encode(input)
	log.Printf("encoded: %v bytes %0.2f%%", len(encoded), float64(len(encoded))/float64(len(input)))

	decoded, err := decode(encoded, tree)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("decoded:", len(input), "bytes")

	if bytes.Equal(input, decoded) {
		log.Println("ok")
	} else {
		log.Println("oops!..")
	}

	os.Stdout.Write(decoded)
}
