## eg-commands-inline example

<!--tmpl,chomp,code=go:cat main.go -->
``` go 
package main

import (
	"log"

	"github.com/jpillora/opts"
)

type Config struct{}

func main() {
	opts.New(&Config{}).
		AddCommand(
			opts.NewNamed(&Foo{}, "foo").
				AddCommand(
					opts.NewNamed(&Bar{}, "bar"),
				),
		).
		Parse().
		RunFatal()
}

type Foo struct {
	Ping string
	Pong string
}

func (f *Foo) Run() error {
	log.Printf("foo: %+v", f)
	return nil
}

type Bar struct {
	Zip string
	Zop string
}

func (b *Bar) Run() error {
	log.Printf("bar: %+v", b)
	return nil
}
```
<!--/tmpl-->

```
$ eg-commands-inline foo bar --zip 2
```

<!--tmpl,chomp,code=plain:go run main.go foo bar --zip 2 -->
``` plain 
2019/04/19 21:49:41 bar: &{Zip:2 Zop:}
```
<!--/tmpl-->

```
$ eg-commands-inline --help
```

<!--tmpl,chomp,code=plain:go build -o eg-commands-inline && ./eg-commands-inline --help && rm eg-commands-inline -->
``` plain 

  Usage: eg-commands-inline [options] <command>

  Options:
  --help, -h

  Commands:
  • foo

```
<!--/tmpl-->
