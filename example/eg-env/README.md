## env example

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
	//In this case UseEnv() is equivalent to
	//adding `env:"FOO"` and `env:"BAR"` tags
	opts.New(&c).UseEnv().Parse()
	fmt.Println(c.Foo)
	fmt.Println(c.Bar)
}
```
<!--/tmpl-->

```
$ export FOO=hello
$ export BAR=world
$ go run env.go
```

<!--tmpl,chomp,code=plain:(export FOO=hello && export BAR=world && go run main.go) -->
``` plain 
hello
world
```
<!--/tmpl-->

```
$ eg-env --help
```

<!--tmpl,chomp,code=plain:go build -o eg-env && ./eg-env --help && rm eg-env -->
``` plain 

  Usage: eg-env [options]

  Options:
  --foo, -f   env FOO
  --bar, -b   env BAR
  --help, -h

```
<!--/tmpl-->
