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

type foo struct{}
type man struct{}
type chew struct{}

func main() {
	config := Config{}
	opts.New(&config).
		Complete().
		AddCommand(opts.New(&foo{}).Name("foo").AddCommand(
			opts.New(&man{}).Name("man").AddCommand(
				opts.New(&chew{}).Name("chew")))).
		Parse().
		RunFatal()
}
