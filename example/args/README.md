## args example

<tmpl,code=go:cat args.go>
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

	c := &Config{}

	opts.Parse(c)

	fmt.Println(c.Foo)
	fmt.Println(c.Bar)
}
```
</tmpl>
```
$ args --foo hello --bar world
```
<tmpl,code:go run args.go --foo hello --bar world>
``` plain 
hello
world
```
</tmpl>
```
$ args --help
```
<tmpl,code:go run args.go --help>
``` plain 

  Usage: args [options]
  
  Options:
  --foo, -f 
  --bar, -b 
  --help, -h
  
```
</tmpl>