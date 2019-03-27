package opts

import "reflect"

//node is the main class, it contains
//all parsing state for a single set of
//arguments
type node struct {
	//embed item since an node can also be an item
	item
	parent       *node
	cmds         map[string]*node
	flags        []*item
	args         []*item
	arglist      *argumentlist
	optnames     map[string]bool
	envnames     map[string]bool
	order        []string
	templates    map[string]string
	internalOpts struct {
		//pretend these are in the user struct :)
		Help, Version bool
	}
	cfgPath               string
	erred                 error
	cmdname               *reflect.Value
	repo, author, version string
	pkgrepo, pkgauthor    string
	//LineWidth defines where new-lines
	//are inserted into the help text
	//(defaults to 42)
	LineWidth int
	//PadAll enables padding around the
	//help text (defaults to true)
	PadAll bool
	//PadWidth defines the amount padding
	//when rendering help text (defaults to 2)
	PadWidth int
}

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

//argumentlist represends a
//named string slice
type argumentlist struct {
	item
	min int
}

func fork(parent *node, val reflect.Value) *node {
	order := make([]string, len(DefaultOrder))
	copy(order, DefaultOrder)
	tmpls := map[string]string{}
	//instantiate
	n := &node{
		item: item{
			val: val,
		},
		parent: parent,
		//each cmd/cmd has its own set of names
		optnames: map[string]bool{},
		envnames: map[string]bool{},
		cmds:     map[string]*node{},
		flags:    []*item{},
		//these are only set at the root
		order:     order,
		templates: tmpls,
		//public defaults
		LineWidth: 72,
		PadAll:    true,
		PadWidth:  2,
	}
	//all fields from val
	if val.Type().Kind() != reflect.Ptr {
		n.errorf("opts: %s should be a pointer to a struct", val.Type().Name())
		return n
	}
	n.addFields(val.Elem())
	//add help option
	g := reflect.ValueOf(&n.internalOpts).Elem()
	n.addOptArg(g.Type().Field(0), g.Field(0))
	//
	return n
}
