# opts

A minimalist CLI library for Go.

`opts` automatically creates `flag.FlagSet`s from your configuration structs using `reflect`. Given the following configuration:

``` go
type Config struct {
	A string `help:"a string"`
	B int    `help:"an int"`
}
```

Then:

``` go
c := Config{}
opts.Parse(&c)
```

Will *approximately* perform the following:

``` go
config := Config{}
set := flag.NewFlagSet("Config")
set.StringVar(&config.A, "", "a string")
set.IntVar(&config.B, 0, "an int")
set.Parse(os.Args)
```

This is quite manageable with two only options, though when we reach 20 options with subcommands, it quickly becomes a chore to keep each value and their flag in sync.

### Features

* Easy to use
* Promotes separation of CLI code and library code
* Automatically generated `--help` text
* Help text via struct tags `help:"Foo bar"`
* Subcommands by nesting structs (each struct represents a `flag.FlagSet`)
* Default values by modifying the struct prior to `Parse()`
* Default values from JSON file, unmarshalled via your config struct
* Default values from environment, defined by your field names
* Infers program name from package name (and optional repository link)
* Extensible via `flag.Value`

### [Simple Example](example/simple)

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
	opts.Parse(&c)
	fmt.Println(c.Foo)
	fmt.Println(c.Bar)
}
```

```
$ ./myprog --foo hello --bar world
hello
world
```

``` plain 
$ ./myprog --help

  Usage: myprog [options]
  
  Options:
  --foo, -f 
  --bar, -b 
  --help, -h
  
```

### More examples

* [Sub-commands](example/subcmds)
* [Args](example/arg)
* [ArgList](example/args)
* [Defaults](example/defaults)
* [Environment Variables](example/env)
* [JSON Config](example/env)
* [Custom Flag Types](example/customtypes)

### Todo

* More tests
* Sub-command help
* Mention env vars in help when enabled

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
