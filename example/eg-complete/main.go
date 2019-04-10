package main

import (
	"fmt"

	"github.com/jpillora/opts"
	"github.com/posener/complete"
)

type Config struct {
	// pstring is a custom type with a Predict function
	Alpha pstring
	Bar
	Foo
}

type Foo struct {
	Ping string
	Pong string
}

type Bar struct {
	Zip string
	Zop string
}

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
	config := Config{}
	// The Name must match the exec's name for self completion to work
	opts.New(&config).Name("eg-complete").
		Complete().
		AddCommand(opts.New(&foo{}).Name("fooie").AddCommand(
			opts.New(&man{}).Name("man").AddCommand(
				opts.New(&chew{}).Name("chew")))).
		Parse().
		RunFatal()
}

func (obj *Config) Run() {
	fmt.Printf("%+v\n", obj)
}
