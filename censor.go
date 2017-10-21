package textcensor

import (
	"io/ioutil"
	"strings"
)

var defaultPunctuation = " !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~，。？；：”’￥（）——、！……"

var defaultCaseSensitive = false

const bomHead = 65279

type sTree struct {
	Root *runeNode
}

type runeNode struct {
	isEnd    bool
	children map[rune]*runeNode
}

var tree = &sTree{&runeNode{false, make(map[rune]*runeNode, 1000)}}

var punctuations = getPunctuationMap(defaultPunctuation)

//SetPunctuation set some punctuations you want to ignore in the strict mode
func SetPunctuation(str string) {
	punctuations = getPunctuationMap(str)
}

func (rNode *runeNode) add(character rune, n *runeNode) {
	if rNode.children == nil {
		rNode.children = make(map[rune]*runeNode)
	}
	rNode.children[character] = n
}

func (rNode *runeNode) find(character rune) *runeNode {
	return rNode.children[character]
}

func initOneWord(str string, caseSensitive bool) {
	str = strings.TrimSpace(str)
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
		if v == bomHead {
			//fmt.Printf("bomHead = %+v\n", bomHead)
			continue
		}
		next := node.find(v)
		if next == nil {
			next = &runeNode{}
			//next.isEnd = i == l-1
			node.add(v, next)
		}
		node = next
		if i == l-1 {
			node.isEnd = true
		}
	}
}

//InitWordsByPath load all the censored word from path
//caseSensitive if true 大小写敏感，即如果FUCK在敏感词库中，fuck没有在库中，则fuck可以通过敏感词检查。
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

//InitWords init the finding tree use given wordsArr
func InitWords(wordsArr []string, caseSensitive bool) {
	//tree = &STree{&Node{false, make(map[rune]*Node, 1000)}}
	defaultCaseSensitive = caseSensitive
	for _, v := range wordsArr {
		initOneWord(v, caseSensitive)
	}
}

//CheckAndReplace check if the text contians bad word, if contains return a newText
//that replaced the bad word with replaceCharacter.
//if strict is true , some punctuations with be ignored.
func CheckAndReplace(text string, strict bool, replaceCharacter rune) (pass bool, newText string) {
	text = strings.TrimSpace(text)
	if len(text) < 1 {
		return true, text
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
			nextNode := node.find(runeArr[j])
			if nextNode == nil && strict && punctuations[runeArr[j]] {
				continue
			}
			node = nextNode
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

	return pass, string(runeArr)
}

//IsPass only check. don't replace words.
//strict if true, some Punctuation will be ignore, eg fuck f*u*c*k f^u^c^k ... can't pass.
func IsPass(text string, strict bool) bool {
	text = strings.TrimSpace(text)
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
			nextNode := node.find(runeArr[j])
			if nextNode == nil && strict && punctuations[runeArr[j]] {
				continue
			}
			node = nextNode
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

func getPunctuationMap(str string) map[rune]bool {
	m := make(map[rune]bool, len(str))
	for _, v := range str {
		m[v] = true
	}
	return m

}
