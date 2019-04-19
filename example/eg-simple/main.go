package main

import (
	"fmt"

	"github.com/jpillora/opts"
)

func main() {
	config := struct {
		File  string `opts:"help=file to load,env=FOO"`
		Lines int    `opts:"help=number of lines to show"`
	}{}
	opts.Parse(&config)
	fmt.Println(config)
}
