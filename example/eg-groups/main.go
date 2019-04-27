package main

import (
	"log"

	"github.com/jpillora/opts"
)

type Config struct {
	Fizz       string
	Buzz       bool
	FooConfig  `opts:"group=Foo"`
	Ping, Pong int `opts:"group=More"`
}

type FooConfig struct {
	Bar  int
	Bazz int
}

func main() {
	c := Config{}
	opts.Parse(&c)
	log.Printf("%+v", c)
}
