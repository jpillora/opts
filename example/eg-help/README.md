## help example

<!--tmpl,chomp,code=go:cat main.go -->

```go
package main

import "github.com/jpillora/opts"

type HelpConfig struct {
	Zip  string `type:"arg" help:"zip is very lorem ipsum dolor sit amet, consectetur adipiscing elit. Phasellus at commodo odio. Sed id tincidunt purus. Cras vel felis dictum, lobortis metus a, tempus tellus"`
	Foo  string `help:"this is help for foo"`
	Bar  int    `help:"and help for bar"`
	Fizz string `help:"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Phasellus at commodo odio. Sed id tincidunt purus. Cras vel felis dictum, lobortis metus a, tempus tellus"`
	Buzz bool   `help:"and help for buzz"`
}

func main() {

	c := HelpConfig{
		Foo: "42",
	}

	opts.New(&c).
		Name("help").
		Version("1.0.0").
		Repo("https://github.com/jpillora/foo").
		Parse()
}
```

<!--/tmpl-->

```
$ eg-help --help
```

<!--tmpl,chomp,code=plain:go run main.go --help -->

```plain

  Usage: help [options] <zip>

  zip is very lorem ipsum dolor sit amet, consectetur adipiscing elit.
  Phasellus at commodo odio. Sed id tincidunt purus. Cras vel felis dictum,
  lobortis metus a, tempus tellus

  Options:
  --foo, -f      this is help for foo (default 42)
  --bar, -b      and help for bar
  --fizz         Lorem ipsum dolor sit amet, consectetur adipiscing elit.
                 Phasellus at commodo odio. Sed id tincidunt purus. Cras
                 vel felis dictum, lobortis metus a, tempus tellus
  --buzz         and help for buzz
  --help, -h
  --version, -v

  Version:
    1.0.0

  Read more:
    https://github.com/jpillora/foo

```

<!--/tmpl-->
