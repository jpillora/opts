package main

import (
	"github.com/jpillora/opts"
	"github.com/jpillora/opts/example/application/foo"
)

//set this via ldflags
var version = "0.0.0"

func main() {
	//configuration with defaults
	app := foo.App{
		Ping: "hello",
		Pong: "world",
	}
	//parse config, note the library version, and extract the
	//repository link from the config package import path
	opts.New(&app).
		Name("foo").
		Version(version).
		PkgRepo(). //includes the infered URL to the package in the help text
		Parse().
		Run()
}
