## subcmds example

<tmpl,code=go:cat subcmds.go>
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
	//optionally store current command name
	Cmd string `type:"cmdname"`
	//subcommand (external struct)
	Foo FooConfig
	//subcommand (inline struct)
	Bar struct {
		Zip string
		Zap string
	}
}

func main() {

	c := Config{}

	opts.Parse(&c)

	fmt.Println(c.Cmd)
	fmt.Println(c.Bar.Zip)
	fmt.Println(c.Bar.Zap)
}
```
</tmpl>
```
$ subcmds bar --zip hello --zap world
```
<tmpl,code:go run subcmds.go bar --zip hello --zap world>
``` plain 
bar
hello
world
```
</tmpl>
```
$ subcmds --help
```
<tmpl,code:go run subcmds.go --help>
``` plain 

  Usage: subcmds [options]
  
  Options:
  --help, -h
  
```
</tmpl>