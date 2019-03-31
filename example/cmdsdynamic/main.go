package main

import (
	"github.com/jpillora/opts"
)

type Config struct {
}

func main() {

	c := Config{}
	f := Foo{}
	opts.New(&c).
		AddCommand(opts.New(&f)).
		Parse()
}

type Foo struct {
	Ping string
	Pong string
}
