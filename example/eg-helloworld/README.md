## helloworld example

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
$ eg-helloworld --file zip.txt --lines 42
```

<!--tmpl,chomp,code=plain:go run main.go --file zip.txt --lines 42 -->
``` plain 
2019/04/27 12:23:55 {File:zip.txt Lines:42}
```
<!--/tmpl-->

```
$ eg-helloworld --help
```

<!--tmpl,chomp,code=plain:go build -o eg-helloworld && ./eg-helloworld --help ; rm eg-helloworld -->
``` plain 

  Usage: eg-helloworld [options]

  Options:
  --file, -f   file to load
  --lines, -l  number of lines to show
  --help, -h

```
<!--/tmpl-->
