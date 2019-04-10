package main

import (
	"github.com/jpillora/opts"
	"github.com/jpillora/opts/example/eg-commands-register/foo"
)

type cmd struct{}

func main() {
	c := cmd{}
	//NOTE: since no name is set, the file name will
	//be used to call .Name("eg-commands-register")
	o := opts.New(&c)
	foo.Register(o)
	o.Parse().RunFatal()
}
