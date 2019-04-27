## groups example

<!--tmpl,chomp,code=go:cat main.go -->
``` go 
package main

import (
	"log"

	"github.com/jpillora/opts"
)

type Config struct {
	Fizz       string
	Buzz       bool
	FooConfig  `opts:"group=Foo"`
	Ping, Pong int `opts:"group=More"`
}

type FooConfig struct {
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
  --help, -h

  Foo options:
  --bar
  --bazz

  More options:
  --ping, -p
  --pong

```
<!--/tmpl-->
