package main

import (
	"github.com/jpillora/opts"
)

type Config struct {
	Ping string
	Pong string
	Zip  string
	Zop  string
}

func main() {
	config := Config{}
	opts.New(&config).
		Complete().
		Parse().
		RunFatal()
}
