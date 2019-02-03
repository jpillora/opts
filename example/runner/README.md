## dyncmds example

<tmpl,code=go:cat dyncmds.go>
``` go 
package main

import (
	"fmt"

	"github.com/jpillora/opts"
)

type FooConfig struct {
	Ping string
	Pong string
}

//config
type Config struct {
	Cmd string `type:"cmdname"`
	//command (external struct)
	Foo FooConfig
	//command (inline struct)
	Bar struct {
		Zip string
		Zap string
	}
}

type Ping struct {
	Fizz int
}

func main() {
	c := Config{}
	p1 := &Ping{}
	p2 := &Ping{}
	p3 := &Ping{}
	opts.New(&c).
		AddSubCmd("pings", p1).
		/**/ GetSubCmd("bar").
		/**/ AddSubCmd("pong", p2).
		/*  */ GetSubCmd("pong").
		/*  */ AddSubCmd("pongs", p3).
		/**/ Parent().
		Parent().
		Parse()
	fmt.Println(c.Cmd)
	fmt.Printf("c  %+v\n", c)
	fmt.Printf("pings %+v\n", p1)
	fmt.Printf("pong  %+v\n", p2)
	fmt.Printf("pongs %+v\n", p3)
}
```
</tmpl>
```
$ ./dyncmds bar pong pongs -f 12
```
<tmpl,code:go run dyncmds.go bar pong pongs -f 12>
``` bar.pong.pongs
c  {Cmd:bar.pong.pongs Foo:{Ping: Pong:} Bar:{Zip: Zap:}}
pings &{Fizz:0}
pong  &{Fizz:0}
pongs &{Fizz:12}
```
</tmpl>
```
$ cmds --help
```
<tmpl,code:go run dyncmds.go --help>
``` plain 

  Usage: dyncmds [options] <command>

  Options:
  --help, -h

  Commands:
  • foo
  • bar
  • pings

```
</tmpl>