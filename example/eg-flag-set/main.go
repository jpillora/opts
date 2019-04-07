package main

import (
	"flag"
	"fmt"

	"github.com/golang/glog"
	"github.com/jpillora/opts"
)

type root struct {
}
type foo struct {
}
type man struct {
}
type chew struct {
}

func main() {
	ro := opts.New(&root{}).
		AddGlobalFlagSet().
		Complete().
		SetLineWidth(132).
		AddCommand(opts.New(&foo{}).Name("foo").AddCommand(
			opts.New(&man{}).Name("man").AddCommand(
				opts.New(&chew{}).Name("chew"))))
	po := ro.Parse()
	_ = po
	po.RunFatal()
}

func (rt *root) Run() {
	fmt.Printf("----\n")
	flag.CommandLine.VisitAll(func(f *flag.Flag) {
		fmt.Printf("fxxx %+v\n", f)
	})
	glog.Infof("root %v\n", rt)
}

func (obj *foo) Run() { println("foo") }
func (obj *man) Run() { println("man") }
func (obj *chew) Run() {
	println("\nchew")
}
