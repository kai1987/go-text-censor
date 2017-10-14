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
	{"毛泽东", true, false},
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

var passTestsPunctuation = []struct {
	in   string
	mode bool
	out  bool
}{
	{"台湾独立", true, false},
	{"台湾独立", false, false},
	{"台*湾*独立", true, false},
	{"Fuck", true, false},
	{"Fiiabcduck", true, false},
	{"共ssdl3^&**#(#())产iiiss!!~		34党Fiiabcduck", true, false},
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
	{"习i猪i头不好", true, '-', false, "-----不好"},
	{"fuckaaa", true, '-', false, "----aaa"},
	{"fiiuaackaaa", true, '-', false, "--------aaa"},
	{"fuckbdfiiuaackaaa", true, '-', false, "----bd--------aaa"},
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
			t.Errorf("str %s,mode %t should be %t , casesensitive:%t", v.in, v.mode, v.out, defaultCaseSensitive)
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

func TestIsPassPunctuation(t *testing.T) {
	InitWordsByPath("./censored_words.txt", false)
	defaultPunctuation := "abcdefghijklmnopqrstuvwxyz !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~，。？；：”’￥（）——、！……"
	SetPunctuation(defaultPunctuation)
	for _, v := range passTestsPunctuation {
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
	defaultPunctuation := "abcdefghijklmnopqrstuvwxyz !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~，。？；：”’￥（）——、！……"
	SetPunctuation(defaultPunctuation)
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

func BenchmarkReplaceLong(b *testing.B) {
	InitWordsByPath("./censored_words.txt", false)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CheckAndReplace(`中国台湾是好友友啊，不要打傣啊。你妹,,回来不是再开发布会，也不必写公开信，也不用再讲梦想，而是把“蒙眼狂奔”后的眼罩摘下来，认真地协助现有的管理层处理当前危机。你可以做的事情真的很多。
		比如：梳理乐视系的所有业务条线，有价值者存之，业已败坏者弃之，让造血机制尽快恢复；比如：与乐视的所有债主和供应商见面——哪怕被他们咬死，也要在保镖的卫护下与他们一起坐下来，探讨尽善事宜，解决所有浮现的危机和可能存在的隐患。
		特别是那几位借钱给你的长江商学院同学，你不能让一息尚存的同学情谊和江湖意气也遭遇背叛。
		更重要的是，回来安抚和安顿每一位乐视的员工。`, true, '*')
	}
}
