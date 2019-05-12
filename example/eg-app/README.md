## library separation

_`main.go`_

<!--tmpl,chomp,code=go:cat main.go -->
``` go 
package main

import (
	"github.com/jpillora/opts"
	"github.com/jpillora/opts/example/eg-app/foo"
)

//set this via ldflags (see https://stackoverflow.com/q/11354518)
var version = "0.0.0"

func main() {
	//new app with some defaults
	app := foo.App{Ping: "hello", Pong: "world"}
	opts.
		New(&app).        //initialise
		Complete().       //enable shell-completion
		Version(version). //use version string set at compile time
		PkgRepo().        //infer the repo URL from package and include in the help text
		Parse().          //where the magic happens, exits with 1 on error
		RunFatal()        //executes App's Run method, exits with 1 on error
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

func (a *App) Run() {
	a.bar = 42 + a.Zip
	a.bazz = 21 + a.Zop
	println("App is running: %+v", a)
}
```
<!--/tmpl-->

```sh
# build program and set 'version' at compile time
$ go build -ldflags "-X main.version=0.2.6" -o foo
$ ./foo --help
```

<!--tmpl,chomp,code=plain:go build -ldflags "-X main.version=0.2.6" -o eg-app && ./eg-app --help ; rm eg-app -->
``` plain 

  Usage: eg-app [options]

  Options:
  --ping, -p       default hello
  --pong           default world
  --zip, -z
  --zop
  --version, -v    display version
  --help, -h       display help

  Completion options:
  --install, -i    install fish-completion
  --uninstall, -u  uninstall fish-completion

  Version:
    0.2.6

  Read more:
    https://github.com/jpillora/opts

```
<!--/tmpl-->
