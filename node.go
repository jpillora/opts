package opts

import "reflect"

//item is the structure representing a
//an opt item
type item struct {
	val       reflect.Value
	name      string
	shortName string
	envName   string
	useEnv    bool
	typeName  string
	help      string
	defstr    string
}

//node is the main class, it contains
//all parsing state for a single set of
//arguments
type node struct {
	err error
	//embed item since an node can also be an item
	item
	parent   *node
	flags    []*item
	args     []*item
	arglist  *argumentlist
	optnames map[string]bool
	envnames map[string]bool
	cfgPath  string
	//subcommands
	cmdname *reflect.Value
	cmds    map[string]*node
	//help
	order                 []string
	templates             map[string]string
	repo, author, version string
	pkgrepo, pkgauthor    string
	lineWidth             int
	padAll                bool
	padWidth              int
	internalOpts          struct {
		//pretend these are in the user struct :)
		Help    bool
		Version bool
	}
	complete struct {
		enabled bool
	}
}

func newNode(val reflect.Value) *node {
	return &node{
		item: item{
			val: val,
		},
		parent: nil,
		//each cmd/cmd has its own set of names
		optnames: map[string]bool{},
		envnames: map[string]bool{},
		cmds:     map[string]*node{},
		flags:    []*item{},
		//these are only set at the root
		order:     defaultOrder(),
		templates: map[string]string{},
		//public defaults
		lineWidth: 72,
		padAll:    true,
		padWidth:  2,
	}
}

//argumentlist represends a
//named string slice
type argumentlist struct {
	item
	min int
}
