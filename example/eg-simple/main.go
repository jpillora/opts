package main

import (
	"log"

	"github.com/jpillora/opts"
)

func main() {
	config := struct {
		File  string `opts:"help=file to load"`
		Lines int    `opts:"help=number of lines to show"`
	}{}
	opts.Parse(&config)
	log.Printf("%+v", config)
}
