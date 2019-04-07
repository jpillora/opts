## cmds example

<!--tmpl,chomp,code=go:cat main.go -->
``` go 
package main

import (
	"github.com/jpillora/opts"
)

type Config struct {
	Foo
	Bar
}

func main() {
	c := Config{}
	opts.Parse(&c).RunFatal()
}

type Foo struct {
	Ping string
	Pong string
}
type Bar struct {
	Ping string
	Pong string
}
```
<!--/tmpl-->

```
$ cmds bar --zip hello --zap world
```

<!--tmpl,chomp,code=plain:go run main.go bar --zip hello --zap world -->
``` plain 

  Usage:  bar [options]

  Options:
  --ping, -p
  --pong
  --help, -h

  Error:
    flag provided but not defined: -zip

```
<!--/tmpl-->

```
$ cmds --help
```

<!--tmpl,chomp,code=plain:go run main.go --help -->
``` plain 

  Usage:  [options] <command>

  Options:
  --help, -h

  Commands:
  • foo
  • bar

```
<!--/tmpl-->
