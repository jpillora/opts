package opts

import (
	"flag"
	"reflect"
	"strings"
)

var flagValueType = reflect.TypeOf((*flag.Value)(nil)).Elem()

type Opts interface {
	//configure this opts node
	Name(name string) Opts
	Version(version string) Opts
	Repo(repo string) Opts
	PkgRepo() Opts
	Author(author string) Opts
	PkgAuthor() Opts
	SetPadWidth(p int) Opts
	DisablePadAll() Opts
	SetLineWidth(l int) Opts
	ConfigPath(path string) Opts
	UseEnv() Opts
	DocBefore(target, newID, template string) Opts
	DocAfter(target, newID, template string) Opts
	DocSet(id, template string) Opts
	//parse this opts node and its children
	Parse() ParsedOpts
	ParseArgs(args []string) ParsedOpts
	//for testing
	process([]string) error
}

type ParsedOpts interface {
	//display this node in the help-output form
	Help() string
}

//New creates a new node instance
func New(config interface{}) Opts {
	v := reflect.ValueOf(config)
	//nil parent -> root command
	n := fork(nil, v)
	if n.erred != nil {
		//opts has already encounted an error
		return n
	}
	//attempt to infer package name, repo, author
	pkgpath := v.Elem().Type().PkgPath()
	parts := strings.Split(pkgpath, "/")
	if len(parts) >= 3 {
		n.pkgauthor = parts[1]
		n.Name(parts[2])
		switch parts[0] {
		case "github.com", "bitbucket.org":
			n.pkgrepo = "https://" + strings.Join(parts[0:3], "/")
		}
	}
	return n
}

//Parse is shorthand for New(config).Parse()
func Parse(config interface{}) ParsedOpts {
	return New(config).Parse()
}
