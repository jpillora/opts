## config example

<!--tmpl,chomp,code=go:cat main.go -->
``` go 
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
	c := Config{}
	opts.New(&c).
		ConfigPath("config.json").
		Parse()
	fmt.Println(c.Foo)
	fmt.Println(c.Bar)
}
```
<!--/tmpl-->

<!--tmpl,chomp,code=json:cat config.json -->
``` json 
{
	"foo": "hello",
	"bar": "world"
}
```
<!--/tmpl-->

```
$ config --bar moon
```

<!--tmpl,chomp,code=plain:go run main.go --bar moon -->
``` plain 

  Usage: main [options]

  Options:
  --foo
  --bar
  --help

  Error:
    Invalid config file: json: Unmarshal(non-pointer main.Config)

```
<!--/tmpl-->

```
$ config --help
```

<!--tmpl,chomp,code=plain:go build -o eg-config && ./eg-config --help ; rm eg-config -->
``` plain 

  Usage: eg-config [options]

  Options:
  --foo
  --bar
  --help

  Error:
    Invalid config file: json: Unmarshal(non-pointer main.Config)

```
<!--/tmpl-->
