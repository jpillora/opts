package bar

import (
	"log"

	"github.com/jpillora/opts"
)

func Register(parent opts.Opts) {
	c := cmd{}
	o := opts.New(&c)
	parent.AddCommand(o)
}

type cmd struct {
	Zip string
	Zop string
}

func (b *cmd) Run() error {
	log.Printf("bar: %+v", b)
	return nil
}
