# opts

A minimalist, yet powerful CLI library for Go

:warning: In progress

---

### Examples

Simple

<tmpl,code=go:cat example/simple/simple.go>
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

	opts.New(c).Parse()

	fmt.Println(c.Foo)
	fmt.Println(c.Bar)
}
```
</tmpl>

```
simple --foo hello --bar world
```
<tmpl,code:go run example/simple/simple.go --foo hello --bar world>
``` plain 
hello
world
```
</tmpl>

```
simple --help
```
<tmpl,code:go run example/simple/simple.go --help>
``` plain 

  Usage: simple [options]
  
  Options:
  --foo, -f 
  --bar, -b 
  --help, -h
  
```
</tmpl>

Defaults

<tmpl,code=go:cat example/defaults/defaults.go>
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

	c := &Config{
		Bar: "moon",
	}

	opts.New(c).Parse()

	fmt.Println(c.Foo)
	fmt.Println(c.Bar)
}
```
</tmpl>
```
defaults --foo hello
```
<tmpl,code:go run example/defaults/defaults.go --foo hello>
``` plain 
hello
moon
```
</tmpl>
```
defaults --help
```
<tmpl,code:go run example/defaults/defaults.go --help>
``` plain 

  Usage: defaults [options]
  
  Options:
  --foo, -f 
  --bar, -b    (default moon).
  --help, -h
  
```
</tmpl>


<!-- 
tmpl,code=go:cat example/defaults/defaults.go></tmpl>
```
defaults --foo hello
```
tmpl,code:go run example/defaults/defaults.go --foo hello></tmpl>
```
defaults --help
```
tmpl,code:go run example/defaults/defaults.go --help></tmpl>
 -->

#### MIT License

Copyright © 2015 &lt;dev@jpillora.com&gt;

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
'Software'), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
