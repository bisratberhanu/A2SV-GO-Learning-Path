package main

import (
    "fmt"
    "strings"
)

func reverseString(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

func ispalindrom(s string) bool {
    // Convert to lowercase to make the check case-insensitive
    s = strings.ToLower(s)
    // Remove spaces to handle multi-word palindromes
    s = strings.ReplaceAll(s, " ", "")
    return s == reverseString(s)
}

func main() {
    fmt.Println(ispalindrom("A man a plan a canal Panama")) // Output: true
    fmt.Println(ispalindrom("hello"))                       // Output: false
}