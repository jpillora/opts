package initopts

import (
	"fmt"

	"github.com/jpillora/opts"
)

type initOpts struct {
}

func Register(parent opts.Opts) opts.Opts {
	in := initOpts{}
	parent.AddCommand(opts.New(&in).Name("init"))
	return parent
}

func (in *initOpts) Run() error {
	fmt.Printf("#init %v\n", in)
	return nil
}
