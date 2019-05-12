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
	panic("TODO")
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
panic: TODO

goroutine 1 [running]:
main.main()
	/Users/jpillora/Code/Go/src/github.com/jpillora/opts/example/eg-config/main.go:15 +0x39
exit status 2
```
<!--/tmpl-->

```
$ config --help
```

<!--tmpl,chomp,code=plain:go build -o eg-config && ./eg-config --help ; rm eg-config -->
``` plain 
panic: TODO

goroutine 1 [running]:
main.main()
	/Users/jpillora/Code/Go/src/github.com/jpillora/opts/example/eg-config/main.go:15 +0x39
```
<!--/tmpl-->
