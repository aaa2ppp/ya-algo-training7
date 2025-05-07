package main

import (
	"bytes"
	"io"
	"log"
	"os"

	"ya-training7/less3/j+/lzw/huffman"
	"ya-training7/less3/j+/lzw/lzw"
)

const maxBits = 13

func encode(input []byte) []byte {
	codes := lzw.Encode(input, maxBits)
	return huffman.Encode(codes)
}

func decode(input []byte) ([]byte, error) {
	codes, err := huffman.Decode[uint16](input)
	if err != nil {
		return nil, err
	}
	return lzw.Decode(codes), nil
}

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("input:", len(input), "bytes")

	encoded := encode(input)
	log.Printf("encoded: %v bytes %0.2f%%", len(encoded), float64(len(encoded))/float64(len(input))*100)

	decoded, err := decode(encoded)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	log.Println("decoded:", len(input), "bytes")

	if bytes.Equal(input, decoded) {
		log.Println("ok")
	} else {
		log.Println("oops!..")
	}

	os.Stdout.Write(decoded)
}
