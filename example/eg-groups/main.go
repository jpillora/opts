package main

import (
	"log"

	"github.com/jpillora/opts"
)

type Config struct {
	Fizz string
	Buzz bool
	//Foo has an implicit `opts:"type=embedded,group=Foo"`.
	//Could be be merged with config by unsetting group `opts:"group="`.
	Foo
	Ping, Pong int `opts:"group=More"`
}

type Foo struct {
	Bar  int
	Bazz int
}

func main() {
	c := Config{}
	opts.Parse(&c)
	log.Printf("%+v", c)
}
