## arg example

<!--tmpl,chomp,code=go:cat main.go -->

```go
package main

import (
	"fmt"

	"github.com/jpillora/opts"
)

type Config struct {
	Foo string `type:"arg" help:"foo is a very important argument"`
	Bar string
}

func main() {

	c := Config{}

	opts.New(&c).Parse()

	fmt.Println(c.Foo)
	fmt.Println(c.Bar)
}
```

<!--/tmpl-->

```
$ eg-arg --foo hello --bar world
```

<!--tmpl,chomp,code=plain:go run main.go --foo hello --bar world -->

```plain

  Usage:  [options] <foo>

  foo is a very important argument

  Options:
  --bar, -b
  --help, -h

  Error:
    flag provided but not defined: -foo

```

<!--/tmpl-->

```
$ arg --help
```

<!--tmpl,chomp,code=plain:go run main.go --help -->

```plain

  Usage:  [options] <foo>

  foo is a very important argument

  Options:
  --bar, -b
  --help, -h

```

<!--/tmpl-->
