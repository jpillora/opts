## cmds example

<!--tmpl,chomp,code=go:cat main.go -->
``` go 
package main

import (
	"log"

	"github.com/jpillora/opts"
)

type Config struct {
	//register commands by including them
	//in the parent struct
	Foo  `opts:"type=cmd,help=This text also becomes commands description text"`
	*Bar `opts:"type=cmd,help=command two of two"`
}

func main() {
	c := Config{}
	opts.NewNamed(&c, "eg-commands-inline").
		Parse().
		Run()
}

type Foo struct {
	Ping string
	Pong string
}

func (f *Foo) Run() error {
	log.Printf("foo: %+v", f)
	return nil
}

type Bar struct {
	Ping string
	Pong string
}

func (b *Bar) Run() error {
	log.Printf("bar: %+v", b)
	return nil
}
```
<!--/tmpl-->

```
$ cmds bar --zip hello --zap world
```

<!--tmpl,chomp,code=plain:go run main.go bar --zip hello --zap world -->
``` plain 

  Usage: eg-commands-inline bar [options]

  command two of two

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

<!--tmpl,chomp,code=plain:go build -o eg-commands-inline && ./eg-commands-inline --help ; rm eg-commands-inline -->
``` plain 

  Usage: eg-commands-inline [options] <command>

  Options:
  --help, -h

  Commands:
  • foo - This text also becomes commands description text
  • bar - command two of two

```
<!--/tmpl-->
