## simple example

<!--tmpl,chomp,code=go:cat main.go -->
``` go 
package main

import (
	"fmt"

	"github.com/jpillora/opts"
)

func main() {
	config := struct {
		File  string `help:"file to load"`
		Lines int    `help:"number of lines to show"`
	}{}
	opts.Parse(&config)
	fmt.Println(config)
}
```
<!--/tmpl-->

```
$ eg-simple --file zip.txt --lines 42
```

<!--tmpl,chomp,code=plain:go run main.go --file zip.txt --lines 42 -->
``` plain 
{zip.txt 42}
```
<!--/tmpl-->

```
$ eg-simple --help
```

<!--tmpl,chomp,code=plain:go build -o eg-simple && ./eg-simple --help && rm eg-simple -->
``` plain 

  Usage: eg-simple [options]

  Options:
  --file, -f   file to load
  --lines, -l  number of lines to show
  --help, -h

```
<!--/tmpl-->
