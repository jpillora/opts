package main

import (
	"fmt"

	"github.com/jpillora/flag"
)

func main() {

	type Config struct {
		Foo  string `help:"this is help for foo"`
		Bar  string `help:"and help for bar"`
		Fizz string `help:"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Phasellus at commodo odio. Sed id tincidunt purus. Cras vel felis dictum, lobortis metus a, tempus tellus"`
		Buzz string `help:"and help for buzz"`
	}

	c := &Config{
		Buzz: "42",
	}

	h := flag.
		New(c).
		Version("1.0.0").
		Repo("https://github.com/jpillora/foo").
		Author("jpillora").
		Help()

	fmt.Print(h)
}
