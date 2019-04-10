package foo

import (
	"log"

	"github.com/jpillora/opts"
	"github.com/jpillora/opts/example/eg-commands-register/bar"
)

func Register(parent opts.Opts) {
	c := cmd{}
	//NOTE: default name for all subcommands is the package name ("foo")
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
