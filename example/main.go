package main

import (
	"fmt"

	"github.com/jpillora/flag"
)

type Config struct {
	Foo string
	Bar string
}

func main() {
	c := &Config{}

	flag.New(c).Parse()

	fmt.Println(c.Foo)
	fmt.Println(c.Bar)
}
