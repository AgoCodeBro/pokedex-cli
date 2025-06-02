package main

import ("fmt"
		"strings"
		)

func main() {
	fmt.Printf("Hello, World!")
}

func cleanInput(text string) []string {
	var result []string
	text = strings.ToLower(text)
	result = strings.Fields(text)

	return result
}
