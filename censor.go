package main

import (
	"io/ioutil"
	"strings"
)

var defaultPunctuation = " !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~，。？；：”’￥（）——、！……"

var defaultCaseSensitive = false

type STree struct {
	Root *Node
}

type Node struct {
	isEnd     bool
	character rune
	children  map[rune]*Node
}

var tree = STree{&Node{false, 0, make(map[rune]*Node, 1000)}}

var PUNCS = getPunctuationMap(defaultPunctuation)

func SetPunctuation(str string) {
	PUNCS = getPunctuationMap(str)
}

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

func initOneWord(str string, caseSensitive bool) {
	l := len(str)
	if l <= 0 {
		return
	}
	if !caseSensitive {
		str = strings.ToLower(str)
	}
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
		if i == l-1 {
			node.isEnd = true
		}
	}
}

//InitWords load all the censored word from path
//caseSensitive if true 大小写敏感，即如果FUCK在敏感词库中，如果fuck没有在库中，则fuck可以通过敏感词检查。
//file should be like example
func InitWordsByPath(path string, caseSensitive bool) error {
	words, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	str := string(words)
	wordsArr := strings.Split(str, "\n")
	InitWords(wordsArr, caseSensitive)
	return nil
}

func InitWords(wordsArr []string, caseSensitive bool) {
	defaultCaseSensitive = caseSensitive
	for _, v := range wordsArr {
		initOneWord(v, caseSensitive)
	}
}

//CheckAndReplace check the text contians bad word, if contains return a newText
//that replaced the bad word with replaceCharacter.
func CheckAndReplace(text string, strict bool, replaceCharacter rune) (pass bool, newText string, err error) {
	if len(text) < 1 {
		return true, text, nil
	}
	if !defaultCaseSensitive {
		text = strings.ToLower(text)

	}
	runeArr := []rune(text)
	l := len(runeArr)

	pass = true

	for i := 0; i < l; i++ {
		cWord := runeArr[i]
		node := tree.Root.find(cWord)

		if node == nil {
			continue
		}
		if node.isEnd {
			runeArr[i] = replaceCharacter
			pass = false
			continue
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
				for ri := i; ri <= j; ri++ {
					runeArr[ri] = replaceCharacter
					pass = false
				}
			}
		}

	}

	return pass, string(runeArr), nil
}

//IsPass only check. don't replace words.
//strict if true, some Punctuation will be ignore, eg fuck f*u*c*k f^u^c^k ... can't pass.
func IsPass(text string, strict bool) bool {
	if len(text) < 1 {
		return true
	}
	if !defaultCaseSensitive {
		text = strings.ToLower(text)
	}
	runeArr := []rune(text)
	l := len(runeArr)

	for i := 0; i < l; i++ {
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

func getPunctuationMap(str string) map[rune]bool {
	m := make(map[rune]bool, len(str))
	for _, v := range str {
		m[v] = true
	}
	return m

}
