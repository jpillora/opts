## args example

<!--tmpl,chomp,code=go:cat main.go -->
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

<!--tmpl,chomp,code=plain:go run main.go foo bar -->
``` plain 
0 foo
1 bar
```
<!--/tmpl-->

```
$ args --help
```

<!--tmpl,chomp,code=plain:go run main.go --help -->
``` plain 

  Usage:  [options] bazzes...

  Options:
  --help, -h

```
<!--/tmpl-->
