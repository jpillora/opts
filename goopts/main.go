package main

import (
	"fmt"
	"os"

	"github.com/jpillora/opts"
	"github.com/jpillora/opts/goopts/internal/initopts"
)

var (
	Version string
	Data    string
	Commit  string
)

type root struct {
}

func main() {
	r := root{}
	ro := opts.New(&r).Name("goopts").
		EmbedGlobalFlagSet().
		Complete().
		Version(Version)

	initopts.Register(ro)

	err := ro.Parse().Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "run error %v\n", err)
		os.Exit(2)
	}
}
