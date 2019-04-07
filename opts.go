package opts

import (
	"flag"
	"reflect"
)

//Opts is configuration command. It represents a node
//in a tree of commands and subcommands. Use the AddCommand
//method to add subcommands (child nodes) to this command.
type Opts interface {
	//configure this opts node
	Name(name string) Opts
	Version(version string) Opts
	ConfigPath(path string) Opts
	UseEnv() Opts
	Complete() Opts
	//documentation
	Description(desc string) Opts
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
	//import standard library flags
	ImportFlagSet(*flag.FlagSet) Opts
	ImportGlobalFlagSet() Opts
	//subcommands
	AddCommand(Opts) Opts
	//parse this opts node and its children
	Parse() ParsedOpts
	ParseArgs(args []string) ParsedOpts
	parse([]string) error
}

type ParsedOpts interface {
	//Display this node in the help-output form
	Help() string
	//Whether the Run method of the parsed command exists
	IsRunnable() bool
	//Execute the Run method of the parsed command.
	//The target Run method must be one of:
	//- Run() error
	//- Run()
	Run() error
	//Execute the Run method of the parsed command and
	//exit(1) with the returned error.
	RunFatal()
}

//New creates a new Opts instance
func New(config interface{}) Opts {
	return newNode(reflect.ValueOf(config))
}

//NewNamed creates a new Opts instance with the given name.
//It is shorthand for `New().Name(name)`
func NewNamed(config interface{}, name string) Opts {
	return New(config).Name(name)
}

//Parse is shorthand for `New(config).Parse()`
func Parse(config interface{}) ParsedOpts {
	return New(config).Parse()
}
