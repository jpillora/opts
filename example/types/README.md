## types example

<tmpl,code=go:cat types.go>
``` go 
package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jpillora/opts"
)

//custom types are also allowed
type Config struct {
	Foo  time.Duration
	Bar  MagicInt
	Bazz int
}

//as long as they implement the flag.Value interface
type MagicInt int

func (b *MagicInt) String() string {
	return strconv.Itoa(int(*b))
}

func (b *MagicInt) Set(s string) error {
	n, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*b = MagicInt(n + 42)
	return nil
}

func main() {

	c := &Config{}

	opts.Parse(c)

	fmt.Println(c.Foo.Seconds())
	fmt.Println(c.Bar)
	fmt.Println(c.Bazz)
}
```
</tmpl>
```
$ types --foo 2m --bar 5 --bazz 5
```
<tmpl,code:go run types.go --foo 2m --bar 5 --bazz 5>
``` plain 
120
47
5
```
</tmpl>
```
$ types --help
```
<tmpl,code:go run types.go --help>
``` plain 

  Usage: types [options]
  
  Options:
  --foo, -f 
  --bar, -b 
  --bazz    
  --help, -h
  
```
</tmpl>