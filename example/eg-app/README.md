## library separation

_`main.go`_

<!--tmpl,chomp,code=go:cat main.go -->
``` go 
package main

import (
	"github.com/jpillora/opts"
	"github.com/jpillora/opts/example/eg-app/foo"
)

//set this via ldflags
var version = "0.0.0"

func main() {
	//new app with some defaults
	app := foo.App{Ping: "hello", Pong: "world"}
	opts.
		New(&app).        //initialise
		Complete().       //enable shell-completion
		Version(version). //use version string set at compile time
		PkgRepo().        //infer the repo URL from package and include in the help text
		Parse().          //exits with 1 if any problems
		RunFatal()        //executes App's Run method
}
```
<!--/tmpl-->

_`foo/foo.go`_

<!--tmpl,chomp,code=go:cat foo/foo.go -->
``` go 
package foo

type App struct {
	//configurable fields
	Ping string
	Pong string
	Zip  int
	Zop  int
	//internal state
	bar  int
	bazz int
}

func (f *App) Run() {
	f.bar = 42 + f.Zip
	f.bazz = 21 + f.Zop
	println("Foo is running...")
}
```
<!--/tmpl-->

```sh
# build program and set 'version' at compile time
$ go build -ldflags "-X main.version=0.2.6" -o foo
$ ./foo --help
```

<!--tmpl,chomp,code=plain:go build -ldflags "-X main.version=0.2.6" -o tmp && ./tmp --help && rm tmp -->
``` plain 

  Usage: tmp [options]

  Options:
  --ping, -p       default hello
  --pong           default world
  --zip, -z
  --zop
  --help, -h
  --version, -v
  --install, -i    install shell-completion
  --uninstall, -u  uninstall shell-completion

  Version:
    0.2.6

  Read more:
    https://github.com/jpillora/opts

```
<!--/tmpl-->
