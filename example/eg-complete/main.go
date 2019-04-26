package main

import (
	"github.com/jpillora/opts"
)

type Config struct {
	Alpha   string
	Charlie string
	Delta   string
	Foo     `opts:"type=cmd"`
	Bar     `opts:"type=cmd"`
}

type Foo struct {
	Ping  string
	Pong  string
	Files []opts.File
}

type Bar struct {
	Zip string
	Zop string
	Dir opts.Dir
}

func main() {
	config := Config{}
	opts.New(&config).
		Complete().
		Parse()
}
