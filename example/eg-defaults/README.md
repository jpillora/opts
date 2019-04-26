## defaults example

<!--tmpl,chomp,code=go:cat main.go -->
``` go 
package main

import (
	"fmt"

	"github.com/jpillora/opts"
)

type Config struct {
	Foo string
	Bar string `opts:"default=world"` //only changes help text
}

func main() {
	c := Config{Foo: "hello"}
	opts.Parse(&c)
	fmt.Println(c.Foo)
	fmt.Println(c.Bar)
}
```
<!--/tmpl-->

```
$ defaults --foo hello
```

<!--tmpl,chomp,code=plain:go run main.go --foo hello -->
``` plain 
hello

```
<!--/tmpl-->

```
$ defaults --help
```

<!--tmpl,chomp,code=plain:go build -o eg-defaults && ./eg-defaults --help ; rm eg-defaults -->
``` plain 

  Usage: eg-defaults [options]

  Options:
  --foo, -f   default hello
  --bar, -b   default world
  --help, -h

```
<!--/tmpl-->
