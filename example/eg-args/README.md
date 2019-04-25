## args example

<!--tmpl,code=go:cat main.go -->
``` go 
package main

import (
	"fmt"

	"github.com/jpillora/opts"
)

type Config struct {
	Bazzes []string `min:"2"`
}

func main() {

	c := Config{}

	opts.New(&c).Parse()

	for i, foo := range c.Bazzes {
		fmt.Println(i, foo)
	}
}
```
<!--/tmpl-->

```
$ args --foo hello --bar world
```

<!--tmpl,code=plain:go run main.go foo bar -->
``` plain 
0 foo
1 bar
```
<!--/tmpl-->

```
$ args --help
```

<!--tmpl,code=plain:go build -o eg-args && ./eg-args --help && rm eg-args -->
``` plain 

  Usage: eg-args [options] bazzes...

  Options:
  --help, -h

```
<!--/tmpl-->
