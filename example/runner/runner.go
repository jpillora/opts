package main

import (
	"fmt"

	"github.com/jpillora/opts"
)

type FooConfig struct {
	Ping string
	Pong string
}

//config
type Config struct {
	Cmd string `type:"cmdname"`
	//command (external struct)
	Foo FooConfig
	//command (inline struct)
	Bar struct {
		Zip string
		Zap string
	}
}

type Ping struct {
	Fizz int
}

func (p *Ping) Run() error {
	fmt.Printf("Running Ping %+v\n", p)
	return nil
}

func main() {
	c := Config{}
	p1 := &Ping{}
	p2 := &Ping{}
	p3 := &Ping{}
	opts.New(&c).
		AddSubCmd("pings", p1).
		/**/ GetSubCmd("bar").
		/**/ AddSubCmd("pong", p2).
		/*  */ GetSubCmd("pong").
		/*  */ AddSubCmd("pongs", p3).
		/**/ Parent().
		Parent().
		Parse().Run()
}
