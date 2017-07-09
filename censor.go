package main

import (
	"io/ioutil"
	"strings"
)

type STree struct {
	Root *Node
}

type Node struct {
	isEnd     bool
	character rune
	children  map[rune]*Node
}

var tree = STree{&Node{false, 0, make(map[rune]*Node, 1000)}}

func (this *Node) add(n *Node) {
	children := this.children
	if len(children) < 1 {
		children = make(map[rune]*Node)
	}

	children[n.character] = n
}

func (this *Node) find(character rune) *Node {
	return this.children[character]
}

func initOneWord(str string) {
	l := len(str)
	if l <= 0 {
		return
	}
	node := tree.Root
	for i, v := range str {
		next := node.find(v)
		if next == nil {
			next = &Node{i == l-1, v, make(map[rune]*Node)}
			node.add(next)
		}
		node = next
	}
}

//InitWords load all the censored word from path
func InitWords(path string) error {
	words, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	str := string(words)
	wordsArr := strings.Split(str, "\n")
	for _, v := range wordsArr {
		initOneWord(v)
	}
	return nil
}

//CheckAndReplace check the text contians bad word, if contains return a newText
//that replaced the bad word with replaceCharacter.
func CheckAndReplace(text string, replaceCharacter string) (pass bool, newText string, err error) {
	return true, text, nil
}

//IsPass only check. don't replace words.
func IsPass(text string) bool {
	l := len(text)
	if l < 1 {
		return true
	}

	runArr := []rune(text)

	node := tree.Root
	for i, j := 0, 1; i < l-1; {
		cWord := runArr[i]
		node = node.find(cWord)
		if node == nil {
			i++
			j = i + 1
			continue
		}
		if node.isEnd {
			return false
		}
		nWord := runArr[j]
		node = node.find(nWord)

	}
	return true
}

func GetTree() STree {
	return tree

}
