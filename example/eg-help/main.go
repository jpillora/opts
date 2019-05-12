package main

import "github.com/jpillora/opts"

type HelpConfig struct {
	Zip  string `opts:"mode=arg,<zip> is a required arg which lorem ipsum dolor sit amet, consectetur adipiscing elit"`
	Foo  string `help:"this is help for foo"`
	Bar  int    `help:"and help for bar"`
	Fizz string `help:"lorem ipsum dolor sit amet, consectetur adipiscing elit. Phasellus at commodo odio. Sed id tincidunt purus. Cras vel felis dictum, lobortis metus a, tempus tellus, and fizz"`
	Buzz bool   `help:"and help for buzz"`
}

func main() {

	c := HelpConfig{
		Foo: "42",
	}

	opts.New(&c).
		Name("help").
		Description("The help program demonstrates how to customise the help text").
		Version("1.0.0").
		Repo("https://github.com/jpillora/foo").
		Parse()
}
