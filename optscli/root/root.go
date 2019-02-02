package root

import (
	"github.com/jpillora/opts"
)

type root struct {
}

var rootOpt = opts.New(&root{})

func Singleton() *opts.Opts {
	return rootOpt
}
