package main

import (
	"github.com/jpillora/opts"
)

type FooConfig struct {
	Ping string
	Pong string
}

type Config struct {
	Cmd string `type:"cmdname"`
}

func main() {

	c := Config{}
	oc := opts.New(&c)

	f := FooConfig{}
	of := opts.New(&f)

	oc.AddCommand(of)

	opts.Parse(&c)
}
