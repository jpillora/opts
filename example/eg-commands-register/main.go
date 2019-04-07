package main

import (
	"github.com/jpillora/opts"
	"github.com/jpillora/opts/example/eg-commands-register/foo"
)

type cmd struct{}

func main() {
	c := cmd{}
	o := opts.New(&c)
	//NOTE: Name() will be set to the binary name 'eg-commands-register'
	foo.Register(o)
	o.Parse().RunFatal()
}
