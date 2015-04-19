package main

import (
	"fmt"

	"github.com/jpillora/opts"
)

type Config struct {
	Foo string `arg:"!" help:"foo is a very important argument"`
	Bar string
}

func main() {

	c := &Config{}

	opts.New(c).Version("1.0.0").Parse()

	fmt.Println(c.Foo)
	fmt.Println(c.Bar)
}
