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
