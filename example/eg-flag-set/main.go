package main

import (
	"github.com/golang/glog"
	"github.com/jpillora/opts"
)

func main() {
	opts.New(&app{}).
		ImportGlobalFlagSet().
		Complete().
		SetLineWidth(90).
		Parse().
		RunFatal()
}

type app struct {
}

func (a *app) Run() {
	glog.Infof("hello from app via glog\n")
}
