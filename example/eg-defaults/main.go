package main

import (
	"fmt"

	"github.com/jpillora/opts"
)

type Config struct {
	Foo string
	Bar string `opts:"default=world"` //only changes help text
}

func main() {
	c := Config{Foo: "hello"}
	opts.Parse(&c)
	fmt.Println(c.Foo)
	fmt.Println(c.Bar)
}
