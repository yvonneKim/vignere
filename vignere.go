package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var alphabet string

func help() {
	fmt.Println(`
		Examples:
			vignere --keyword kolache --decode my_ciphered_file.txt
			vignere --keyword kolache --encode file_to_be_encoded.txt
	`)
}

type Config struct {
	keyword string
	encode  bool
	file    io.Reader
}

func getConfig() Config {
	config := Config{}
	var err error
	for i := 1; i < len(os.Args)-1; i += 2 {
		flag, value := os.Args[i], os.Args[i+1]
		if flag == "--keyword" {
			config.keyword = value
		}

		if flag == "--decode" {
			config.encode = false
			config.file, err = os.Open(value)
		}

		if flag == "--encode" {
			config.encode = true
			config.file, err = os.Open(value)
		}

		if err != nil {
			log.Fatal(err)
		}
	}

	return config
}

func encode(plaintext string, keyword string) string {
	result := ""
	plaintext = strings.ToLower(plaintext)
	for i, letter := range plaintext {
		letter_index := strings.Index(alphabet, string(letter))

		keyword_letter := keyword[i%len(keyword)]
		keyword_index := strings.Index(alphabet, string(keyword_letter))

		cipher_index := (letter_index + keyword_index) % len(alphabet)
		cipher_letter := string(alphabet[cipher_index])

		result += cipher_letter
	}
	return result
}

func decode(ciphertext string, keyword string) string {
	result := ""
	ciphertext = strings.ToLower(ciphertext)
	for i, letter := range ciphertext {
		letter_index := strings.Index(alphabet, string(letter))

		keyword_letter := keyword[i%len(keyword)]
		keyword_index := strings.Index(alphabet, string(keyword_letter))

		plain_index := (letter_index - keyword_index) % len(alphabet)
		plain_letter := string(alphabet[plain_index])

		result += plain_letter
	}
	return result
}

func run(c Config) {
	s := bufio.NewScanner(c.file)

	var result string
	for s.Scan() {
		if c.encode {
			result = encode(s.Text(), c.keyword)
		} else {
			result = decode(s.Text(), c.keyword)
		}
	}

	fmt.Println(result)
}

func main() {
	alphabet = "abcdefghijklmnopqrstuvwxyz"
	config := getConfig()
	run(config)
}
