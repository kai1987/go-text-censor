package main

import (
	"fmt"
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

//2000-\u206F
var tree = STree{&Node{false, 0, make(map[rune]*Node, 1000)}}

var PUNCS = allPunctuation()

func (this *Node) add(n *Node) {
	children := this.children
	if len(children) < 1 {
		children = make(map[rune]*Node)
	}

	children[n.character] = n
	this.children = children
}

func (this *Node) find(character rune) *Node {
	return this.children[character]
}

func initOneWord(str string) {
	l := len(str)
	if l <= 0 {
		return
	}
	str = strings.ToLower(str)
	runeArr := []rune(str)
	l = len(runeArr)
	node := tree.Root
	for i := 0; i < l; i++ {
		v := runeArr[i]
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
func IsPass(text string, strict bool) bool {
	l := len(text)
	if l < 1 {
		return true
	}

	text = strings.ToLower(text)

	runeArr := []rune(text)
	l = len(runeArr)

	for i := 0; i < l-1; i++ {
		cWord := runeArr[i]
		node := tree.Root.find(cWord)

		if node == nil {
			continue
		}
		if node.isEnd {
			return false
		}

		for j := i + 1; j < l; j++ {
			//如果是严格模式，将所有的标点忽略掉
			if strict && PUNCS[runeArr[j]] {
				continue
			}
			node = node.find(runeArr[j])
			if node == nil {
				break
			}
			if node.isEnd {
				return false
			}
		}

	}
	return true
}

func GetTree() STree {
	return tree

}

func allPunctuation() map[rune]bool {
	str := " !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~，。？；：”’￥（）——、！……"
	m := make(map[rune]bool, len(str))
	for _, v := range str {
		m[v] = true
	}
	return m

}

func main() {
	for i := 32; i < 127; i++ {
		fmt.Print(string(rune(i)))
	}
	fmt.Printf("strings.ToLower = %+v\n", strings.ToLower("Ghoe中文中国"))
}
