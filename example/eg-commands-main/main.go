package main

import (
	"log"

	"github.com/jpillora/opts"
)

type Config struct{}

func main() {
	opts.New(&Config{}).
		AddCommand(
			opts.NewNamed(&Foo{}, "foo").
				AddCommand(
					opts.NewNamed(&Bar{}, "bar"),
				),
		).
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
