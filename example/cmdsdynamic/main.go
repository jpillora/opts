package main

import (
	"log"

	"github.com/jpillora/opts"
)

type Config struct{}

func main() {

	bar := Bar{}
	b := opts.New(&bar).Name("bar")

	foo := Foo{}
	f := opts.New(&foo).Name("foo").AddCommand(b)

	config := Config{}
	opts.New(&config).
		Name("root").
		AddCommand(f).
		Parse().
		RunFatal()
}

type Foo struct {
	Ping string
	Pong string
}

func (f *Foo) Run() error {
	log.Printf("foo: %+v", f)
	return nil
}

type Bar struct {
	Zip string
	Zop string
}

func (b *Bar) Run() error {
	log.Printf("bar: %+v", b)
	return nil
}
