package main

import (
	"github.com/jpillora/opts"
)

type Config struct {
	Foo
	Bar
}

func main() {
	c := Config{}
	opts.Parse(&c).RunFatal()
}

type Foo struct {
	Ping string
	Pong string
}
type Bar struct {
	Ping string
	Pong string
}
