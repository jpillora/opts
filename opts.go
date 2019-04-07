package opts

import (
	"flag"
	"reflect"
)

type Opts interface {
	//configure this opts node
	Name(name string) Opts
	Version(version string) Opts
	ConfigPath(path string) Opts
	UseEnv() Opts
	Complete() Opts
	AddFlagSet(*flag.FlagSet) Opts
	AddGoCommandLineFlagSet() Opts
	//documentation
	Repo(repo string) Opts
	Author(author string) Opts
	PkgRepo() Opts
	PkgAuthor() Opts
	DocBefore(target, newID, template string) Opts
	DocAfter(target, newID, template string) Opts
	DocSet(id, template string) Opts
	DisablePadAll() Opts
	SetPadWidth(padding int) Opts
	SetLineWidth(width int) Opts
	//subcommands
	AddCommand(Opts) Opts
	//parse this opts node and its children
	Parse() ParsedOpts
	ParseArgs(args []string) ParsedOpts
	parse([]string) error
}

type ParsedOpts interface {
	//display this node in the help-output form
	Help() string
	//whether the Run method of the parsed command exists
	IsRunnable() bool
	//execute the Run method of the parsed command
	Run() error
	//execute the Run method of the parsed command
	RunFatal()
}

//New creates a new node instance
func New(config interface{}) Opts {
	return newNode(reflect.ValueOf(config))
}

//Parse is shorthand for New(config).Parse()
func Parse(config interface{}) ParsedOpts {
	return New(config).Parse()
}
