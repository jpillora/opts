## args example

<!--tmpl,chomp,code=go:cat main.go -->
``` go 
package main

import (
	"log"

	"github.com/jpillora/opts"
)

type Config struct {
	Octopi []string `opts:"type=arg,min=2"`
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
2019/04/26 22:15:48 SINGLE: octopus

  Usage: main [options] <octopus>

  allows multiple

  Options:
  --help, -h

  Error:
    Unexpected arguments: [bar]

```
<!--/tmpl-->

```
$ args --help
```

<!--tmpl,chomp,code=plain:go build -o eg-args && ./eg-args --help ; rm eg-args -->
``` plain 
2019/04/26 22:15:48 SINGLE: octopus

  Usage: eg-args [options] <octopus>

  allows multiple

  Options:
  --help, -h

```
<!--/tmpl-->
