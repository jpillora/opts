package main

import (
	"log"

	"github.com/jpillora/opts"
)

type Config struct {
	Shark  string   `opts:"mode=arg"`
	Octopi []string `opts:"mode=arg,min=2"`
}

func main() {
	c := Config{}
	opts.New(&c).Parse()
	log.Printf("%+v", c)
}
