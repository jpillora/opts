package main

import (
	"fmt"

	"github.com/jpillora/opts"
)

type Config struct {
	Foo string
	Bar string
}

func main() {

	c := &Config{
		Bar: "moon",
	}

	opts.New(c).Parse()

	fmt.Println(c.Foo)
	fmt.Println(c.Bar)
}
