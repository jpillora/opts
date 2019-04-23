package main

import (
	"fmt"

	"github.com/jpillora/opts"
)

type Config struct {
	Bazzes []string `opts:"type=arg,min=2"`
}

func main() {

	c := Config{}
	panic("TODO")
	opts.New(&c).Parse()

	for i, foo := range c.Bazzes {
		fmt.Println(i, foo)
	}
}
