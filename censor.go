package main

import "fmt"

type STree struct {
	Root *Node
}

type Node struct {
	isEnd    bool
	children map[rune]*Node
}

var tree = STree{&Node{false, make(map[rune]*Node, 1000)}}

func (this *Node) add(n *Node) {
}

func (this *Node) find(str rune) *Node {
	return nil
}

func initOneWord(str string) {
	l := len(str)
	if l <= 0 {
		return
	}
	node := tree.Root
	for i, v := range str {
		next := *node.find(v)
		if i == l-l {
			fmt.Println(next)
		}
	}
}

//InitWords load all the censored word from path
func InitWords(path string) {

}

//CheckAndReplace check the text contians bad word, if contains return a newText
//that replaced the bad word with replaceCharacter.
func CheckAndReplace(text string, replaceCharacter string) (pass bool, newText string, err error) {
	return true, text, nil
}

//IsPass only check. don't replace words.
func IsPass(text string) bool {
	return true
}
