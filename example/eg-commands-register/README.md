## commands using the `Register` pattern

_`main.go`_

<!--tmpl,chomp,code=go:cat main.go -->
``` go 
package main

import (
	"github.com/jpillora/opts"
	"github.com/jpillora/opts/example/eg-commands-register/foo"
)

type cmd struct{}

func main() {
	c := cmd{}
	o := opts.New(&c)
	foo.Register(o)
	o.Parse().RunFatal()
}
```
<!--/tmpl-->

_`foo/cmd.go`_

<!--tmpl,chomp,code=go:cat foo/cmd.go -->
``` go 
package foo

import (
	"log"

	"github.com/jpillora/opts"
	"github.com/jpillora/opts/example/eg-commands-register/bar"
)

func Register(parent opts.Opts) {
	c := cmd{}
	//NOTE: default name for a subcommand is its package name ("foo")
	o := opts.New(&c)
	bar.Register(o)
	parent.AddCommand(o)
}

type cmd struct {
	Ping string
	Pong string
}

func (f *cmd) Run() error {
	log.Printf("foo: %+v", f)
	return nil
}
```
<!--/tmpl-->

```sh
$ ./eg-commands-register --help
```

<!--tmpl,chomp,code=plain:go build -o eg-commands-register && ./eg-commands-register --help ; rm eg-commands-register -->
``` plain 

  Usage: eg-commands-register [options] <command>

  Options:
  --help, -h  display help

  Commands:
  • cmd

```
<!--/tmpl-->
