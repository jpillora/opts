package main

import (
	"github.com/jpillora/opts"
)

type Config struct {
	Alpha string
	Bar
	Charlie string
	Delta   string
	Foo
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
