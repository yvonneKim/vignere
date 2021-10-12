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

// http://pi.math.cornell.edu/~mec/2003-2004/cryptography/subs/frequencies.html
var etaoinsrhdlucmfywgpbvkxqjz = map[rune]int{
	'e': 1202,
	't': 910,
	'a': 812,
	'o': 768,
	'i': 731,
	'n': 695,
	's': 628,
	'r': 602,
	'h': 592,
	'd': 432,
	'l': 398,
	'u': 288,
	'c': 271,
	'm': 261,
	'f': 230,
	'y': 211,
	'w': 209,
	'g': 203,
	'p': 182,
	'b': 149,
	'v': 111,
	'k': 69,
	'x': 17,
	'q': 11,
	'j': 10,
	'z': 7,
}

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

func detect_offset(frequencies map[rune]int) int {
	var most_frequent_letter rune
	var highest_count int
	for letter, count := range frequencies {
		if count > highest_count {
			most_frequent_letter = letter
			highest_count = count
		}
	}

	index_difference := int(most_frequent_letter) - int('e')
	if index_difference < 0 {
		index_difference = int('e') - int(most_frequent_letter)
	}
	return index_difference
}

func count_frequencies(ciphertext string, key_len int) map[rune]int {
	result := make(map[rune]int)
	for i := 0; i < len(ciphertext); i += key_len {
		result[rune(ciphertext[i])] += 1
	}
	return result
}

func solve(ciphertext string) string {
	for key_len := 1; key_len < 10; key_len++ {
		freqs := count_frequencies(ciphertext, key_len)
		offset := detect_offset(freqs)
		offset_letter := string(int('a') + offset)
	}

	return ""
}

func sanitize(text string) string {
	result := ""
	for _, letter := range strings.ToLower(text) {
		if letter < 97 || letter > 122 {
			continue
		}
		result += string(letter)
	}

	return result
}

func encode(plaintext string, keyword string) string {
	result := ""
	plaintext = sanitize(plaintext)
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

		if plain_index < 0 {
			plain_index = len(alphabet) + plain_index
		}
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
			fmt.Println(result)
			continue
		}

		if c.keyword != "" {
			result = decode(s.Text(), c.keyword)
			fmt.Println(result)
			continue
		}
		result = solve(s.Text())
	}
}

func main() {
	alphabet = "abcdefghijklmnopqrstuvwxyz"
	config := getConfig()
	run(config)
}
