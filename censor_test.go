package main

import (
	"fmt"
	"testing"
)

func TestInitWords(t *testing.T) {

	InitWords("./censored_words.txt")
	fmt.Printf("tree = %+v\n", GetTree())
}
