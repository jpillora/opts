package main

import (
	"fmt"

	"github.com/jpillora/opts"
	"github.com/posener/complete"
)

type Config struct {
	Ping string   `predict:"files:*.go"`
	Pong string   `predict:"dirs"`
	Zip  []string `predict:"none" type:"commalist"`
	Zop  pstring
	Zing *opts.RepeatedStringOpt
	Zang *opts.RepeatedStringOpt
	Zong opts.RepeatedStringOpt
}

// ./eg-complete --zing a --zing b --zing c
// &{Ping: Pong: Zip: Zop:{val:} Zing:[a b c]}

type pstring struct {
	val string
}

func (*pstring) Predict(args complete.Args) []string {
	return []string{"a", "b", "c"}
}
func (p *pstring) Set(val string) (err error) {
	p.val = val
	return
}
func (p *pstring) String() string {
	return p.val
}

type foo struct {
}
type man struct{}
type chew struct{}

func main() {
	zang := &opts.RepeatedStringOpt{}
	zang.Set("abc")
	config := Config{Zang: zang}
	opts.New(&config).
		Complete().
		AddCommand(opts.New(&foo{}).Name("foo").AddCommand(
			opts.New(&man{}).Name("man").AddCommand(
				opts.New(&chew{}).Name("chew")))).
		Parse().
		RunFatal()
}

func (obj *Config) Run() {
	fmt.Printf("%+v\n", obj)
}
