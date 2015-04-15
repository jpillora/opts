package main

import (
	"fmt"

	"github.com/jpillora/flag"
)

func main() {

	type Config struct {
		Foo string
		Bar string
	}

	c := &Config{}

	flag.New(c).Parse()

	fmt.Println(c.Foo)
	fmt.Println(c.Bar)
}
