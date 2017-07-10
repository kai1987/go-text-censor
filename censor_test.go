package main

import (
	"testing"
)

var passTests = []struct {
	in   string
	mode bool
	out  bool
}{
	{"台湾独立", true, false},
	{"台湾独立", false, false},
	{"台*湾*独立", true, false},
	{"台*湾*独立", false, true},
	{"台--@#*湾*独立", true, false},
	{"台--@#*!！……——湾*独立", true, false},
	{"你好", false, true},
	{"你好", true, true},
	{"台湾你好", true, true},
	{"FUCK", true, false},
	{"fuck", true, false},
	{"f    u  ck", true, false},
	{"FUCK", false, false},
	{"fuck", false, false},
	{"f*u*ck", false, true},
	{"", true, true},
}

func TestInitWords(t *testing.T) {

	InitWords("./censored_words.txt")

	for _, v := range passTests {
		if v.out != IsPass(v.in, v.mode) {
			t.Errorf("str %s,mode %t should be %t", v.in, v.mode, v.out)
		}

	}

}
