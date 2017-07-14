# GO-Text-Censor

A fast censored word check and replace package for Go.<br/>
Support case sensitive or not .<br/>
Support strict mode to ignore none-sence character like *|- ......

## Getting Started

* using `go get`

		go get github.com/kai1987/go-text-censor

* via cloning this repository:

	  git clone git@github.com:kai1987/go-text-censor.git ${GOPATH}/src/github.com/kai1987/go-text-censor


### Usage

all the usages are in censor_test.go

```
textcensor.InitWordsByPath("yourfilepath",false)

isPass := textcensor.IsPass(text,true)

// if you want to replace

isPass,newText := textcensor.CheckAndReplace(text,true,'*')

```


## Running the tests

```
go test -coverprofile textcensor
PASS
coverage: 100.0% of statements
ok    github.com/go-text-censor 0.013s
```

```
go test
```
if you want see the benchmark use
```
go test -bench=.

here is the benchmark on my MacBook Air

BenchmarkIsPassShort-4     2000000         710 ns/op
BenchmarkIsPass-4          1000000        2089 ns/op
BenchmarkReplace-4          500000        2694 ns/op
PASS
ok    github.com/go-text-censor 5.642s

it's ns/op so it's realy fast.

```

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details


----------I'm a naughty separate line, please use Chinese under me -----------



# GO-Text-Censor

这个包用来检查或者过滤敏感词的。你懂的。

支持大小写敏感设置，建议设置为不敏感，即false,

支持去除无用的字符比如*|, 默认排除了大部分英文标点和中文标点

速度见benchmark. 很快的。

