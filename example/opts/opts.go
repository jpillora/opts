package main

import (
	"fmt"

	"github.com/jpillora/opts"
)

type Config struct {
	Bazzes opts.RepeatedStringOpt //`type:"opt"`
	Inizga []string               `type:"commalist"`
}

func main() {
	c := Config{}
	opts.New(&c).Parse()
	for i, foo := range c.Bazzes.GetSlice() {
		fmt.Println(i, foo)
	}
	for i, foo := range c.Inizga {
		fmt.Println(i, foo)
	}
}
