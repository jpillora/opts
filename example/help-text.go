package main

import "github.com/jpillora/flag"

func main() {

	type Config struct {
		Foo  string `help:"this is help for foo."`
		Bar  string `help:"and help for bar."`
		Fizz string `help:"and some more and more and more and more and more and more and more and more and more and more."`
		Buzz string `help:"and help for buzz."`
	}

	c := &Config{}

	f := flag.New(c)
	// f.Version = "1.0.0"
	f.Repo = "https://github.com/jpillora/foo"
	f.Parse()
}
