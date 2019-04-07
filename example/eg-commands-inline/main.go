package main

import (
	"log"

	"github.com/jpillora/opts"
)

type Config struct {
	//register commands by including them
	//in the parent struct
	Foo `help:"This text also becomes commands Description() text"`
	Bar `help:"command two of two"`
}

func main() {
	c := Config{}
	opts.NewNamed(&c, "eg-commands-inline").
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
	Ping string
	Pong string
}

func (b *Bar) Run() error {
	log.Printf("bar: %+v", b)
	return nil
}
