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
