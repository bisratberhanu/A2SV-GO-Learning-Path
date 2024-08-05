package main

import (
	"fmt"
	"strings"
)

func freqCounter(s string) map[string]int {
    freq := make(map[string]int)
    words := strings.Fields(strings.ToLower(s))
    for _, word := range words {
        freq[word]++
    }
    return freq
}

func main() {
    text := "Hello world hello"
    fmt.Println(freqCounter(text))
}