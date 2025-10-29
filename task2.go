package main

import (
	"fmt"
	s "strings"
	"unicode"
)

var print = fmt.Println

func filter(input string) string {
	var filtered []rune
	for _,r:= range input {
		if unicode.IsPunct(r) { continue }
		r=unicode.ToLower(r)
		filtered= append(filtered,r)
	}
	return string(filtered)
}

func word_frequency(target string) map[string]int {
	words := s.Fields(target)
	counter := make(map[string]int)
	
	
	for _, word := range words {
		word = filter(word)
		counter[word]++
	}

	return counter
}



func palindromeChecker(str string) bool {
	str = s.ToLower(str)


	var filtered []rune
	for _, r := range str {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			filtered = append(filtered, r)
		}
	}

	n := len(filtered)
	for i := 0; i < n/2; i++ {
		if filtered[i] != filtered[n-1-i] { return false }
	}

	return true
}

func main() {
	print(word_frequency("the angry the "))
	print(palindromeChecker("heleh"))
}