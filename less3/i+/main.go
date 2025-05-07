package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

func encode(input string) string {
	k := 0
	for (1 << k) < len(input)+k+1 {
		k++
	}

	// копируем ввод в ввывод
	output := make([]byte, len(input)+k+1)
	n := 1
	for i, j, p := 1, 0, 0; j < len(input); i++ {
		n++

		//  пропускаем контрольные биты (степень двойки)
		if i == (1 << p) {
			p++
			continue
		}

		// при копировании преобразуем символы в числа 0/1
		// мне показалось, что так будет удобнее...
		output[i] = input[j] - '0'
		j++
	}
	output = output[:n]

	if debugEnable {
		log.Println("input :", input)
		log.Println("output:", output)
	}

	for p := 0; 1<<p < len(output); p++ {
		var sum byte

		// насчитывае еденицы в позициях, где bit(p) == 1
		for i, c := range output {
			// TODO: это можно сделать только битовыми операциями (без if)
			if i&(1<<p) != 0 {
				sum += c
			}
			// sum += byte(i>>p) & c & 1
		}

		// сумма (с учетом контрольного бита) должна быть нечетна
		output[1<<p] = (sum + 1) & 1

		if debugEnable {
			log.Println("p:", p, "sum:", sum, "->", output[1<<p])
		}
	}

	if debugEnable {
		log.Println("output:", output)
	}

	// переводим output в символы '0'/'1'
	for i := range output {
		output[i] += '0'
	}

	return unsafeString(output[1:]) // отрезаем первый ноль (не несет смысловой нагрузки)
}

func decode(input string) string {
	var k, broken int
	for p := 0; 1<<p < len(input); p++ {
		k = p
		var sum byte

		// насчитывае еденицы в позициях, где bit(p) == 1
		for i, c := range []byte(input) {
			// TODO: это можно сделать только битовыми операциями (без if)
			if (i+1)&(1<<p) != 0 { // +1 учитываем отрезанный ноль
				sum += c & 1 // 48->0, 49->1
			}
		}

		if debugEnable {
			log.Println("p:", p, "sum:", sum)
		}

		// сумма должна быть нечетна
		if sum&1 == 0 {
			// oops!..
			broken |= 1 << p
		}
	}

	if debugEnable {
		log.Println("broken:", broken)
	}

	// копируе биты данных в вывод
	var output strings.Builder
	output.Grow(len(input) - k)

	for i, p := 0, 0; i < len(input); i++ {

		// пропускаем контрольные биты (степень двойки)
		if i+1 == 1<<p { // +1 учитываем отрезанный ноль
			p++
			continue
		}

		c := input[i]

		// если бит битый, то инвертируем его
		if i+1 == broken { // +1 учитываем отрезанный ноль
			c ^= 1 // '0'=48->49, '1'=49->48
		}

		output.WriteByte(c)
	}

	return output.String()
}

func run(in io.Reader, out io.Writer) {
	br := bufio.NewReader(in)
	bw := bufio.NewWriter(out)
	defer bw.Flush()

	var cmd int
	if _, err := fmt.Fscanln(br, &cmd); err != nil {
		panic(err)
	}

	line, err := br.ReadString('\n')
	if err != nil && err != io.EOF {
		panic(err)
	}
	line = strings.TrimRight(line, " \t\r\n")

	var ans string
	switch cmd {
	case 1:
		ans = encode(line)
	case 2:
		ans = decode(line)
	default:
		panic("unknown command " + strconv.Itoa(cmd))
	}

	bw.WriteString(ans)
	bw.WriteByte('\n')
}

// ----------------------------------------------------------------------------

var _, debugEnable = os.LookupEnv("DEBUG")

func main() {
	_ = debugEnable
	run(os.Stdin, os.Stdout)
}

// ----------------------------------------------------------------------------

func unsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
