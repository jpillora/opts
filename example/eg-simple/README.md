## simple example

<!--tmpl,chomp,code=go:cat main.go -->
``` go 
package main

import (
	"fmt"

	"github.com/jpillora/opts"
)

type Config struct {
	Foo string
	Bar string
}

func main() {
	c := Config{}
	opts.Parse(&c)
	fmt.Println(c.Foo)
	fmt.Println(c.Bar)
}
```
<!--/tmpl-->

```
$ eg-simple --foo hello --bar world
```

<!--tmpl,chomp,code=plain:go run main.go --foo hello --bar world -->
``` plain 
hello
world
```
<!--/tmpl-->

```
$ eg-simple --help
```

<!--tmpl,chomp,code=plain:go build -o eg-simple && ./eg-simple --help && rm eg-simple -->
``` plain 

  Usage: eg-simple [options]

  Options:
  --foo, -f
  --bar, -b
  --help, -h

```
<!--/tmpl-->
