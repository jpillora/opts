## args example

<!--tmpl,chomp,code=go:cat main.go -->
``` go 
package main

import (
	"log"

	"github.com/jpillora/opts"
)

type Config struct {
	Shark  string   `opts:"mode=arg"`
	Octopi []string `opts:"mode=arg,min=2"`
}

func main() {
	c := Config{}
	opts.New(&c).Parse()
	log.Printf("%+v", c)
}
```
<!--/tmpl-->

```
$ args --foo hello --bar world
```

<!--tmpl,chomp,code=plain:go run main.go foo bar -->
``` plain 
2019/05/12 20:21:58 {Shark:foo Octopi:[bar]}
```
<!--/tmpl-->

```
$ args --help
```

<!--tmpl,chomp,code=plain:go build -o eg-args && ./eg-args --help ; rm eg-args -->
``` plain 

  Usage: eg-args [options] <shark> <octopus>

  allows multiple

  Options:
  --help, -h  display help

```
<!--/tmpl-->
