package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func randBool() bool {
	return rand.Intn(2) == 0
}

func duplicateVowel(word []byte, i int) []byte {
	return append(word[:i+1], word[i:]...)
}

func removeVowel(word []byte, i int) []byte {
	return append(word[:i], word[i+1:]...)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		word := []byte(s.Text())
		if randBool() {
			vI := -1
			for i, char := range word {
				switch char {
				case 'a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U':
					if randBool() {
						vI = i
					}
				}
			}
			if vI >= 0 {
				if randBool() {
					word = duplicateVowel(word, vI)
				} else {
					word = removeVowel(word, vI)
				}
			}
		}
		fmt.Println(string(word))
	}
}
