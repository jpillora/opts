package main

import (
	"fmt"

	"github.com/jpillora/opts"
)

type multiOpt struct {
	vals []string
}

func (mo *multiOpt) Set(val string) error {
	mo.vals = append(mo.vals, val)
	return nil
}

func (mo *multiOpt) String() string {
	return fmt.Sprintf("%v", mo.vals)
}

type Config struct {
	Bazzes *multiOpt //`type:"opt"`
}

func main() {

	c := Config{
		Bazzes: &multiOpt{},
	}

	opt := opts.New(&c)
	// fmt.Printf("%+v\n", opt)
	opt.Parse()

	for i, foo := range c.Bazzes.vals {
		fmt.Println(i, foo)
	}
}
