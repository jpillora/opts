## help example

<!--tmpl,chomp,code=go:cat main.go -->
``` go 
package main

import "github.com/jpillora/opts"

type HelpConfig struct {
	Zip  string `opts:"type=arg,<zip> is a required arg which lorem ipsum dolor sit amet, consectetur adipiscing elit"`
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
```
<!--/tmpl-->

```
$ eg-help --help
```

<!--tmpl,chomp,code=plain:go build -o eg-help && ./eg-help --help ; rm eg-help -->
``` plain 

  Usage: help [options] <zip>

  The help program demonstrates how to customise the help text

  Version:
    1.0.0

  Read more:
    https://github.com/jpillora/foo

  Error:
    field 'Zip' unused opts keys: [ consectetur adipiscing elit <zip> is a required arg which lorem ipsum dolor sit amet]

```
<!--/tmpl-->
