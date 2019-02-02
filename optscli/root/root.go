package root

import (
	"github.com/jpillora/opts"
)

type root struct {
}

var rootOpt = opts.New(&root{})

func Singleton() opts.Builder {
	return rootOpt
}
