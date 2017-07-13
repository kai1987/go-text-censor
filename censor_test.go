package textcensor

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
	{"台*湾*独*立", false, true}, //因为独立也在默认敏感词中
	{"台--@#*湾*独立", true, false},
	{"台--@#*!！……——湾*独立", true, false},
	{"你好", false, true},
	{"你好", true, true},
	{"台湾你好", true, true},
	{"FUCK", true, false},
	{"fuck", true, false},
	{"FUCK", false, false},
	{"fuck", false, false},
	{"f*u*ck", false, true},
	{"", true, true},
	{"操场", true, false},
	{"操", true, false},
	{"f    u  ck", true, false},
}

var crTests = []struct {
	in   string
	mode bool
	rpc  rune
	pass bool
	out  string
}{
	{"你好", true, '*', true, "你好"},
	{"你好操ni", true, '*', false, "你好*ni"},
	{"你好你**妹ni", true, '*', false, "你好****ni"},
	{"你好你**妹ni", false, '*', true, "你好你**妹ni"},
	{"操", false, '-', false, "-"},
	{"", false, '-', true, ""},
}

var passTestsCaseSensitive = []struct {
	in   string
	mode bool
	out  bool
}{
	{"台湾独立", true, false},
	{"台湾独立", false, false},
	{"台*湾*独立", true, false},
	{"Fuck", true, false},
	{"FuCk", true, true},
}

var crTestsWithSetPunctuation = []struct {
	in   string
	mode bool
	rpc  rune
	pass bool
	out  string
}{
	{"你好", true, '*', true, "你好"},
	{"你好操ni", true, '*', false, "你好*ni"},
	{"你好你**妹ni", true, '*', false, "你好****ni"},
	{"你好你**妹ni", false, '*', true, "你好你**妹ni"},
	{"操", false, '-', false, "-"},
	{"", false, '-', true, ""},
	{"习i猪i头", false, '-', true, "习i猪i头"},
	{"习i猪i头", true, '-', false, "-----"},
	{"习,猪,头", true, '-', true, "习,猪,头"},
	{"习i猪i头不好", true, '-', false, "-----不好"},
}

func TestInitWordsByPath(t *testing.T) {
	err := InitWordsByPath("./not_exist.txt", false)
	if err == nil {
		t.Errorf("init not exist file should have err")
	}
}

func TestIsPass(t *testing.T) {
	InitWordsByPath("./censored_words.txt", false)
	for _, v := range passTests {
		if v.out != IsPass(v.in, v.mode) {
			t.Errorf("str %s,mode %t should be %t , casesensitive", v.in, v.mode, v.out, defaultCaseSensitive)
		}
	}
}

func TestIsPassCaseSensitive(t *testing.T) {
	InitWordsByPath("./censored_words.txt", true)
	for _, v := range passTestsCaseSensitive {
		if v.out != IsPass(v.in, v.mode) {
			t.Errorf("str %s,mode %t should be %t", v.in, v.mode, v.out)
		}
	}
}

func TestCheckAndReplace(t *testing.T) {
	InitWordsByPath("./censored_words.txt", false)
	for _, v := range crTests {
		pass, out := CheckAndReplace(v.in, v.mode, v.rpc)
		if pass != v.pass || out != v.out {
			t.Errorf("message ,v:%v,got:pass:%v,out:%v", v, pass, out)
		}
	}

}
func TestCheckAndReplaceSetPunctuation(t *testing.T) {
	InitWordsByPath("./censored_words.txt", false)
	SetPunctuation("*i")
	for _, v := range crTestsWithSetPunctuation {
		pass, out := CheckAndReplace(v.in, v.mode, v.rpc)
		if pass != v.pass || out != v.out {
			t.Errorf("message ,v:%v,got:pass:%v,out:%v", v, pass, out)
		}
	}
	SetPunctuation(defaultPunctuation)

}

func BenchmarkIsPassShort(b *testing.B) {
	InitWordsByPath("./censored_words.txt", false)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsPass("完全个苏。", true)
	}
}

func BenchmarkIsPass(b *testing.B) {
	InitWordsByPath("./censored_words.txt", false)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsPass("中国台湾是好友友啊，不要打傣啊。", true)
	}
}

func BenchmarkReplace(b *testing.B) {
	InitWordsByPath("./censored_words.txt", false)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CheckAndReplace("中国台湾是好友友啊，不要打傣啊。你妹", true, '*')
	}
}
