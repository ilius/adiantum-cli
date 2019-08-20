package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ilius/go-askpass"
	"lukechampine.com/adiantum"
)

func encodeFromStdin(noNewline bool) {
	input, _ := ioutil.ReadAll(os.Stdin)
	if len(input) == 0 {
		return
	}
	input = bytes.TrimRight(input, "\n")
	if len(input) < 16 {
		zeros := bytes.Repeat([]byte{0}, 16-len(input))
		input = append(input, zeros...)
	}
	keyStr, err := askpass.Askpass("Password", true)
	if err != nil {
		panic(err)
	}
	keyInput := []byte(keyStr)
	if len(keyInput) > 32 {
		panic(fmt.Errorf("Password is %d bytes, more than 32 bytes", len(keyInput)))
	}
	key := make([]byte, 32)
	copy(key[:len(keyInput)], keyInput)
	cipher := adiantum.New(key)
	tweak := make([]byte, 12) // can be any length
	output := cipher.Encrypt(input, tweak)
	if noNewline {
		fmt.Print(string(output))
	} else {
		fmt.Println(string(output))
	}
}

func decodeFromStdin(noNewline bool) {
	input, _ := ioutil.ReadAll(os.Stdin)
	if len(input) == 0 {
		return
	}
	input = bytes.TrimRight(input, "\n")
	if len(input) < 16 {
		zeros := bytes.Repeat([]byte{0}, 16-len(input))
		input = append(input, zeros...)
	}
	keyStr, err := askpass.Askpass("Password", false)
	if err != nil {
		panic(err)
	}
	keyInput := []byte(keyStr)
	if len(keyInput) > 32 {
		panic(fmt.Errorf("Password is %d bytes, more than 32 bytes", len(keyInput)))
	}
	key := make([]byte, 32)
	copy(key[:len(keyInput)], keyInput)
	cipher := adiantum.New(key)
	tweak := make([]byte, 12) // can be any length
	output := cipher.Decrypt(input, tweak)
	// TODO: add a flag to print hex-encoded
	if noNewline {
		fmt.Print(string(output))
	} else {
		fmt.Println(string(output))
	}
	// fmt.Println(hex.EncodeToString(output))
}

func main() {
	decodeFlag := flag.Bool(
		"d",
		false,
		"Decode:\nchunk32 -d",
	)

	noNewlineFlag := flag.Bool(
		"n",
		false,
		"Do not print newline at the end (mostly useful for decode)\nchunk32 -d -n",
	)

	flag.Parse()

	noNewline := noNewlineFlag != nil && *noNewlineFlag
	if decodeFlag != nil && *decodeFlag {
		decodeFromStdin(noNewline)
	} else {
		encodeFromStdin(noNewline)
	}
}
