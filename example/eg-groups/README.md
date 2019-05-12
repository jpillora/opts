## groups example

<!--tmpl,chomp,code=go:cat main.go -->
``` go 
package main

import (
	"log"

	"github.com/jpillora/opts"
)

type Config struct {
	Fizz string
	Buzz bool
	//Foo has an implicit `opts:"mode=embedded,group=Foo"`.
	//Can be unset with `opts:"group="`.
	FooBar
	Ping, Pong int `opts:"group=More"`
}

type FooBar struct {
	Bar  int
	Bazz int
}

func main() {
	c := Config{}
	opts.Parse(&c)
	log.Printf("%+v", c)
}
```
<!--/tmpl-->

Group order in the help text is first-use order

```
$ eg-groups --help
```

<!--tmpl,chomp,code=plain:go build -o eg-groups && ./eg-groups --help ; rm eg-groups -->
``` plain 

  Usage: eg-groups [options]

  Options:
  --fizz, -f
  --buzz, -b
  --help, -h  display help

  Foo options:
  --bar
  --bazz

  More options:
  --ping, -p
  --pong

```
<!--/tmpl-->
