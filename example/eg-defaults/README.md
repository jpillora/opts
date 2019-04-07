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
	Bar string
}

func main() {

	c := Config{
		Bar: "moon",
	}

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
moon
```
<!--/tmpl-->

```
$ defaults --help
```

<!--tmpl,chomp,code=plain:go build -o eg-defaults && ./eg-defaults --help && rm eg-defaults -->
``` plain 

  Usage: eg-defaults [options]

  Options:
  --foo, -f
  --bar, -b   default moon
  --help, -h

```
<!--/tmpl-->
