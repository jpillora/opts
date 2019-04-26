## simple example

<!--tmpl,chomp,code=go:cat main.go -->
``` go 
package main

import (
	"log"

	"github.com/jpillora/opts"
)

func main() {
	config := struct {
		File  string `opts:"help=file to load"`
		Lines int    `opts:"help=number of lines to show"`
	}{}
	opts.Parse(&config)
	log.Printf("%+v", config)
}
```
<!--/tmpl-->

```
$ eg-simple --file zip.txt --lines 42
```

<!--tmpl,chomp,code=plain:go run main.go --file zip.txt --lines 42 -->
``` plain 
2019/04/26 22:15:59 {File:zip.txt Lines:42}
```
<!--/tmpl-->

```
$ eg-simple --help
```

<!--tmpl,chomp,code=plain:go build -o eg-simple && ./eg-simple --help ; rm eg-simple -->
``` plain 

  Usage: eg-simple [options]

  Options:
  --file, -f   file to load
  --lines, -l  number of lines to show
  --help, -h

```
<!--/tmpl-->
