## custom help example

<!--tmpl,code=go:cat main.go -->
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
	//see default templates and the default template order
	//in the opts/help.go file
	o := opts.New(&c).
		DocAfter("usage", "mytext", "\nthis is a some text!\n"). //add new entry
		Repo("myfoo.com/bar").
		DocSet("repo", "\nMy awesome repo:\n  {{.Repo}}"). //change existing entry
		Parse()

	fmt.Println(o.Help())
}
```
<!--/tmpl-->

```
$ eg-help-advanced --help
```

<!--tmpl,code=plain:go build -o eg-help-advanced && ./eg-help-advanced --help && rm eg-help-advanced -->
``` plain 

  Usage: eg-help-advanced [options]

  this is a some text!

  Options:
  --foo, -f
  --bar, -b
  --help, -h

  My awesome repo:
    myfoo.com/bar
```
<!--/tmpl-->
