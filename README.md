# opts

The **v1** release is coming soon!

---

**A low friction command-line interface library for Go (golang)**

[![GoDoc](https://godoc.org/github.com/jpillora/opts?status.svg)](https://godoc.org/github.com/jpillora/opts)  [![CircleCI](https://circleci.com/gh/jpillora/opts.svg?style=shield&circle-token=69ef9c6ac0d8cebcb354bb85c377eceff77bfb1b)](https://circleci.com/gh/jpillora/opts)

Command-line parsing should be easy. Use configuration structs:

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

```sh
$ go build -o my-prog
$ ./my-prog --help
  Usage: my-prog [options]

  Options:
  --file, -f   file to load
  --lines, -l  number of lines to show
  --help, -h

```

```sh
$ ./my-prog -f foo -l 12
{foo 12}
```

### Features (with examples)

* Easy to use ([simple](example/simple/))
* Promotes separation of CLI code and library code ([separation](example/separation/))
* Automatically generated `--help` text via struct tags `help:"Foo bar"` ([help](example/help/))
* Commands by nesting structs ([cmds](example/cmds/))
* Default values by modifying the struct prior to `Parse()` ([defaults](example/defaults/))
* Default values from a JSON config file, unmarshalled via your config struct ([config](example/config/))
* Default values from environment, defined by your field names ([env](example/env/))
* Infers program name from executable name
* Infers sub-command names from package name
* Extensible via `flag.Value` ([customtypes](example/customtypes/))
* Customizable help text by modifying the default templates ([customhelp](example/customhelp/))
* Built-in auto-completion

Find all examples here [`example/`](./example)

### Overview

Internally, `opts` creates `flag.FlagSet`s from your configuration structs using `pkg/reflect`. So, given the following configuration:

```go
type Config struct {
	Alpha   string        `help:"a string"`
	Bravo   int           `help:"an int"`
	Charlie bool          `help:"a bool"`
	Delta   time.Duration `help:"a duration"`
}

c := Config{
	Bravo: 42,
	Delta: 2 * time.Minute,
}
```

When we perform a `opts.Parse`:

```go
opts.Parse(&c)
```

Under the hood, `opts` will *approximately* perform:

```go
set := flag.NewFlagSet("config")
set.StringVar(&c.Alpha, "", "a string")
set.IntVar(&c.Bravo, 42, "an int")
set.BoolVar(&c.Charlie, false, "a bool")
set.DurationVar(&c.Delta, 2 * time.Minute, "a duration")
set.Parse(os.Args)
```

However `opts` also does **much more**, like giving you pretty `--help` output:

```
$ ./foo --help

  Usage: foo [options]

  Options:
  --alpha, -a    a string
  --bravo, -b    an int (default 42)
  --charlie, -c  an bool
  --delta, -d    a duration (default 2m0s)
  --help, -h

```

### Package API

See [![GoDoc](https://godoc.org/github.com/jpillora/opts?status.svg)](https://godoc.org/github.com/jpillora/opts) https://godoc.org/github.com/jpillora/opts

### Struct Tag API

**opts** tries to set sane defaults so, for the most part, you'll get the desired behaviour by simply providing a configuration struct.

However, you can customise this behaviour by providing the `opts` struct
tag with a series of settings in the form of `key=value`:

```
`opts:"key=value,key=value,..."
```

Where `key` must be:

* `name` - Name is used to display the field in the help text (defaults to the field name converted to lowercase and dashes)

* `help` - Help is used to describe the field

	Note: The `help` setting is only setting that can also be set as a
	stand-alone struct tag: `help:"my text goes here"`. You must use the stand-alone tag if you wish to use `=` and `,` in your help string.

* `type` - The type assigned the field (defaults using the table below)

	Where `value` must be one of:

	* `flag` - The field will be treated as a flag. That is, an optional, named, configurable field. Set using `./program --flag-name <flag-value>`.

	* `arg` - The field will be treated as an argument. That is, a required, positional, unamed, configurable field. Set using `./program <argument-value>`.

	* `args` - The field will be treated as an argument list. Set using `./program <argument-value>`.

	* `cmd` - 	A command is nested `opts.Opt` instance, so its fields behave in exactly the same way as the parent struct.

		You can access the options of a command with `prog --prog-opt X cmd --cmd-opt Y`

		Restricted to fields with Go type `struct`.

	* `cmdname` - 	A special type which will assume the name of the selected command

		Restricted to fields with Go type `string`.

	* `embedded` - 	A special type which causes the fields of struct to be used in the current struct. Useful if you want to split your command-line options across multiple files.

		Restricted to fields with Go type `struct`.

	All fields must have a **opts** `type`. By default a struct field will be assigned a `type` depending on its field type:

  | Go Type         | Default opts `type` |      Valid `type`s       |
  | --------------- | :-----------------: | :----------------------: |
  | `int`           |       `flag`        |      `flag`, `arg`       |
  | `string`        |       `flag`        | `flag`, `arg`, `cmdname` |
  | `bool`          |       `flag`        |      `flag`, `arg`       |
  | `flag.Value`    |       `flag`        |      `flag`, `arg`       |
  | `time.Duration` |       `flag`        |      `flag`, `arg`       |
  | `[]string`      |       `flag`        |      `flag`, `args`      |
  | `struct`        |     `embedded`      |    `cmd`, `embedded`     |


* `short` - One or two letters to be used a flags "short" name. Only valid when `type` is `flag`.

* `env` - An environent variable to use as the field's **default** value. It can always be overridden by providing the appropriate flag.

	For example, `opts:"env=FOO"`. It can also be infered using the field name with simply `opts:"env"`. You can enable inference on all flags with the `opts.Opts` method `UseEnv()`.

### Other projects

Other CLI libraries which infer flags from struct tags:

* https://github.com/alexflint/go-arg Doesn't support subcommands
* https://github.com/jessevdk/go-flags Has many features though still no subcommands

### Todo

* More tests
* Multiple short options `-aux` (Requires a non-`pkg/flag` parser)

#### MIT License

Copyright © 2019 &lt;dev@jpillora.com&gt;

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
