## customtypes example

<!--tmpl,chomp,code=go:cat main.go -->
``` go 
package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jpillora/opts"
)

//MagicInt is a valid custom type since it implements the flag.Value interface
type MagicInt int

func (b MagicInt) String() string {
	return "{" + strconv.Itoa(int(b)) + "}"
}

func (b *MagicInt) Set(s string) error {
	n, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*b = MagicInt(n + 42)
	return nil
}

type Config struct {
	Mmm   []MagicInt
	Bar   time.Duration
	Zee   bool
	Files []opts.File
	Dir   opts.Dir
}

func main() {
	c := Config{}
	opts.Parse(&c)
	fmt.Printf("%+v\n", c)
}
```
<!--/tmpl-->

```sh
#NOTE: 5 + 42 = 47
$ eg-custom-flag --foo 2m --bar 5 --bazz 5
```

<!--tmpl,chomp,code=plain:go run main.go --foo 2m --bar 5 --bazz 5 -->
``` plain 

  Usage: main [options]

  Options:
  --mmm, -m   allows multiple
  --bar, -b
  --zee, -z
  --file, -f  allows multiple
  --dir, -d
  --help, -h  display help

  Error:
    flag provided but not defined: -foo

```
<!--/tmpl-->

```
$ eg-custom-flag --help
```

<!--tmpl,chomp,code=plain:go install && eg-custom-flag --help ; rm $(which eg-custom-flag) -->
``` plain 

  Usage: eg-custom-flag [options]

  Options:
  --mmm, -m   allows multiple
  --bar, -b
  --zee, -z
  --file, -f  allows multiple
  --dir, -d
  --help, -h  display help

```
<!--/tmpl-->
