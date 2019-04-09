package main

import (
	"github.com/jpillora/opts"
)

type Config struct {
	Alpha string
	Bar
	Foo
}

type Foo struct {
	Ping string
	Pong string
}

type Bar struct {
	Zip string
	Zop string
}

func main() {
	config := Config{}
	opts.New(&config).
		Complete().
		Parse().
		RunFatal()
}
