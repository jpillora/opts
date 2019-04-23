package opts

import (
	"flag"
	"reflect"
)

//node is the main class, it contains
//all parsing state for a single set of
//arguments
type node struct {
	err error
	//embed item since an node can also be an item
	item
	parent    *node
	flags     []*item
	flagNames map[string]bool
	args      []*item
	envNames  map[string]bool
	cfgPath   string
	//external flagsets
	flagsets []*flag.FlagSet
	//subcommands
	cmd     *node
	cmdname *string
	cmds    map[string]*node
	//help
	order                       []string
	templates                   map[string]string
	repo, author, version, desc string
	repoInfer, authorInfer      bool
	lineWidth                   int
	padAll                      bool
	padWidth                    int
	//pretend these are in the user struct :)
	internalOpts struct {
		Help      bool
		Version   bool
		Install   bool `help:"install shell-completion"`
		Uninstall bool `help:"uninstall shell-completion"`
	}
	complete bool
}

func newNode(val reflect.Value) *node {
	return &node{
		item:   item{val: val},
		parent: nil,
		//each cmd/cmd has its own set of names
		flagNames: map[string]bool{},
		envNames:  map[string]bool{},
		cmds:      map[string]*node{},
		flags:     []*item{},
		//these are only set at the root
		order:     defaultOrder(),
		templates: map[string]string{},
		//public defaults
		lineWidth: 72,
		padAll:    true,
		padWidth:  2,
	}
}
