package main

import (
	"github.com/jpillora/opts"
	"github.com/jpillora/opts/example/application/foo"
)

//set this via ldflags
var version = "0.0.0"

func main() {
	//new app with some defaults
	app := foo.App{Ping: "hello", Pong: "world"}
	opts.
		New(&app).        //infer CLI name from package name "foo"
		Version(version). //use version string set at compile time
		PkgRepo().        //infer the repo URL and include in the help text
		Parse().          //exits with 1 if any problems
		RunFatal()        //executes App's Run method
}
