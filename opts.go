package opts

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/posener/complete/cmd/install"

	"github.com/posener/complete"
)

// var flagValuePtrType = reflect.TypeOf((*flag.Value)(nil)).Elem()
// var flagValueValType = reflect.TypeOf((flag.Value)(nil)).Elem()

type Config interface{}

type Helper interface {
	Help() string
}

type Runner interface {
	Run() error
}

type Predictable interface {
	flag.Value
	complete.Predictor
}

type Builder interface {
	AddSubCmd(name string, cmd Config) Builder
	SubCmd(name string, cmd Config) Builder
	GetSubCmd(name string) Builder
	Parent() Builder
	Complete(name string) Builder
	Name(name string) Builder
	Version(version string) Builder
	Repo(repo string) Builder
	PkgRepo() Builder
	Author(author string) Builder
	PkgAuthor() Builder
	SetPadWidth(p int) Builder
	DisablePadAll() Builder
	SetLineWidth(l int) Builder
	DocBefore(target, newid, template string) Builder
	DocAfter(target, newid, template string) Builder
	DocSet(id, template string) Builder
	ConfigPath(path string) Builder
	UseEnv() Builder
	Parse() Configured
	ParseArgs(args []string) Configured
	process(args []string) (*Opts, []string, error)
}

type Configured interface {
	Help() string
	Run() error
	IsRunner() bool
	Config() Config
}

//Opts is the main class, it contains
//all parsing state for a single set of
//arguments
type Opts struct {
	item
	parent       *Opts
	cmds         map[string]*Opts
	completeCom  *complete.Complete
	completeCmd  complete.Command
	completeExec string
	opts         []*item
	args         []*item
	arglist      *argumentlist
	optnames     map[string]bool
	envnames     map[string]bool
	order        []string
	templates    map[string]string
	internalOpts struct {
		//pretend these are in the user struct :)
		Help, Version, InstallCompletetion, UninstallCompletetion bool
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

//argumentlist represends a
//named string slice
type argumentlist struct {
	item
	min int
}

//item is the structure representing a
//an opt item
type item struct {
	val       reflect.Value
	sf        reflect.StructField
	name      string
	shortName string
	envName   string
	useEnv    bool
	typeName  string
	help      string
	defstr    string
}

//New creates a new Opts instance
func New(config interface{}) Builder {
	v := reflect.ValueOf(config)
	//nil parent -> root command
	o := fork(nil, v)
	if o.erred != nil {
		//opts has already encounted an error
		return o
	}

	//attempt to infer package name, repo, author
	pkgpath := v.Elem().Type().PkgPath()
	parts := strings.Split(pkgpath, "/")
	if len(parts) >= 3 {
		o.pkgauthor = parts[1]
		o.Name(parts[2])
		switch parts[0] {
		case "github.com", "bitbucket.org":
			o.pkgrepo = "https://" + strings.Join(parts[0:3], "/")
		}
	}

	return o
}

func secondPass(o2 *Opts) {
	var recurse func(or *Opts)
	recurse = func(or *Opts) {
		for _, opt := range or.opts {
			//should generate shortname?
			if len(opt.name) >= 3 && opt.sf.Tag.Get("short") != "-" && opt.shortName == "" {
				Log("adding short for %s\n", opt.name)
				//not already taken?
				if s := opt.name[0:1]; !or.optnames[s] {
					opt.shortName = s
					or.optnames[s] = true
					Log("'%s' %s %+v\n", s, opt.name, or.completeCmd.Flags["--"+opt.name])
					pre := or.completeCmd.Flags["--"+opt.name]
					or.completeCmd.Flags["-"+s] = pre
				}
			}
		}
		for _, node := range or.cmds {
			recurse(node)
		}
	}
	recurse(o2)
	//
	Log("completeExec '%s'\n", o2.completeExec)
	if o2.completeExec != "" {
		o2.completeCom = complete.New(o2.completeExec, o2.completeCmd)
		var recurse func(or *Opts)
		recurse = func(or *Opts) {
			for name, node := range or.cmds {
				or.completeCmd.Sub[name] = node.completeCmd
				recurse(node)
			}
		}
		recurse(o2)
	}
}

//Parse(&config) is shorthand for New(&config).Parse()
func Parse(config interface{}) Configured {
	return New(config).Parse()
}

func fork(parent *Opts, val reflect.Value) *Opts {
	//TODO allow order and template per cmd
	//for now, there is only the root
	var order []string = nil
	var tmpls map[string]string = nil
	if parent == nil {
		order = make([]string, len(DefaultOrder))
		copy(order, DefaultOrder)
		tmpls = map[string]string{}
	} else {
		order = parent.order
		tmpls = parent.templates
	}
	//instantiate
	o := &Opts{
		item: item{
			val: val,
		},
		parent: parent,
		//each cmd/cmd has its own set of names
		optnames: map[string]bool{},
		envnames: map[string]bool{},
		cmds:     map[string]*Opts{},
		opts:     []*item{},
		//these are only set at the root
		order:     order,
		templates: tmpls,
		//public defaults
		LineWidth: 72,
		PadAll:    true,
		PadWidth:  2,
		completeCmd: complete.Command{
			Sub:         complete.Commands{},
			Flags:       complete.Flags{},
			GlobalFlags: complete.Flags{},
			// Args:
		},
	}
	//all fields from val
	if val.Type().Kind() != reflect.Ptr {
		o.errorf("opts: %s should be a pointer to a struct", val.Type().Name())
		return o
	}
	o.addFields(val.Elem())
	//add help option
	g := reflect.ValueOf(&o.internalOpts).Elem()
	o.addOptArg(g.Type().Field(0), g.Field(0))
	return o
}

func (o *Opts) addFields(c reflect.Value) *Opts {
	t := c.Type()
	k := t.Kind()
	//deref pointer
	if k == reflect.Ptr {
		c = c.Elem()
		t = c.Type()
		k = t.Kind()
	}
	if k != reflect.Struct {
		o.errorf("opts: %s should be a pointer to a struct (got %s)", t.Name(), k)
		return o
	}
	//parse struct fields
	for i := 0; i < c.NumField(); i++ {
		val := c.Field(i)
		//ignore unexported
		if !val.CanSet() {
			continue
		}
		sf := t.Field(i)
		//ignore `opts:"-"`
		if sf.Tag.Get("opts") == "-" {
			continue
		}
		k := sf.Type.Kind()
		if _, ok := val.Interface().(flag.Value); ok {
			o.addOptArg(sf, val)
			continue
		}
		if _, ok := val.Addr().Interface().(flag.Value); ok {
			o.addOptArg(sf, val)
			continue
		}

		// if sf.Type.Implements(flagValuePtrType) || sf.Type.Implements(flagValueValType) {
		// 	o.addOptArg(sf, val)
		// 	continue
		// }
		switch k {
		case reflect.Ptr, reflect.Struct:
			if sf.Tag.Get("type") == "embedded" {
				o.addFields(val)
			} else {
				o.addCmd(sf, val)
			}
		case reflect.Slice:
			if sf.Type.Elem().Kind() != reflect.String {
				o.errorf("slice (list) types must be []string")
				return o
			}
			if sf.Tag.Get("type") == "commalist" {
				o.addOptArg(sf, val)
			} else if sf.Tag.Get("type") == "spacelist" {
				o.addOptArg(sf, val)
			} else {
				o.addArgList(sf, val)
			}
		case reflect.Bool, reflect.String, reflect.Int, reflect.Int64:
			if sf.Tag.Get("type") == "cmdname" {
				if k != reflect.String {
					o.errorf("cmdname field '%s' must be a string", sf.Name)
					return o
				}
				o.cmdname = &val
			} else {
				o.addOptArg(sf, val)
			}
		case reflect.Interface:
			o.errorf("Struct field '%s' interface type must implement flag.Value", sf.Name)
			return o
		default:
			o.errorf("Struct field '%s' has unsupported type: %s - '%T'\n", sf.Name, k, val)
			panic(o.erred)
			return o
		}
	}

	return o
}

func (o *Opts) errorf(format string, args ...interface{}) *Opts {
	//only store the first error
	if o.erred == nil {
		o.erred = fmt.Errorf(format, args...)
	}
	return o
}

func (o *Opts) addCmd(sf reflect.StructField, val reflect.Value) {
	if o.arglist != nil {
		o.errorf("argslists and commands cannot be used together")
		return
	}
	//requires address
	switch sf.Type.Kind() {
	case reflect.Ptr:
		//if nil ptr, auto-create new struct
		if val.IsNil() {
			ptr := reflect.New(val.Type().Elem())
			val.Set(ptr)
		}
	case reflect.Struct:
		val = val.Addr()
	}
	name := sf.Tag.Get("name")
	if name == "" || name == "!" {
		name = camel2dash(sf.Name) //default to struct field name
	}
	// o.complete.Sub[name] =
	// log.Printf("define cmd: %s =====", subname)
	sub := fork(o, val)
	sub.name = name
	sub.help = sf.Tag.Get("help")
	o.cmds[name] = sub
}

func (o *Opts) addArgList(sf reflect.StructField, val reflect.Value) {

	if len(o.cmds) > 0 {
		o.errorf("argslists and commands cannot be used together")
		return
	}
	if o.arglist != nil {
		o.errorf("only 1 arglist field is allowed ('%s' already defined)", o.arglist.name)
		return
	}

	name := sf.Tag.Get("name")
	if name == "" || name == "!" {
		name = camel2dash(sf.Name) //default to struct field name
	}

	if val.Len() != 0 {
		o.errorf("arglist '%s' is required so it should not be set. "+
			"If you'd like to set a default, consider using an option instead.", name)
		return
	}

	min, _ := strconv.Atoi(sf.Tag.Get("min"))

	//insert
	o.arglist = &argumentlist{
		item: item{
			val:  val,
			sf:   sf,
			name: name,
			help: sf.Tag.Get("help"),
		},
		min: min,
	}
}

var durationType = reflect.TypeOf(time.Second)

func (o *Opts) addOptArg(sf reflect.StructField, val reflect.Value) *item {
	Log("addOptArg %T\n", val.Interface())
	//assume opt, unless arg tag is present
	t := sf.Tag.Get("type")
	if t == "" {
		t = "opt"
	}
	//
	i := &item{
		val:      val,
		sf:       sf,
		typeName: t,
	}
	//find name
	i.name = sf.Tag.Get("name")
	if i.name == "" {
		i.name = camel2dash(sf.Name) //default to struct field name
	}
	//specific environment name
	i.envName = sf.Tag.Get("env")
	if i.envName != "" {
		if o.envnames[i.envName] {
			o.errorf("option env name '%s' already in use", i.name)
			return nil
		}
		o.envnames[i.envName] = true
		i.useEnv = true
	}
	//opt names cannot clash with each other
	if o.optnames[i.name] {
		o.errorf("option name '%s' already in use", i.name)
		return nil
	}
	o.optnames[i.name] = true
	//get help text
	i.help = sf.Tag.Get("help")
	//the **displayed** default, use 'default' tag, otherwise infer
	defstr := sf.Tag.Get("default")
	switch {
	case defstr != "":
		i.defstr = defstr
	case val.Kind() == reflect.Slice:
		i.defstr = ""
	case val.Kind() != reflect.Ptr && val.CanAddr():
		if val.Addr().Interface() != reflect.Zero(val.Type()).Interface() {
			i.defstr = fmt.Sprintf("%v", val.Interface())
		}
	case val.Interface() != reflect.Zero(sf.Type).Interface():
		i.defstr = fmt.Sprintf("%v", val.Interface())
	}
	// if defstr != "" {
	// 	i.defstr = defstr
	// } else if val.Kind() == reflect.Slice {
	// 	i.defstr = ""
	// } else if def := val.Interface(); def != reflect.Zero(sf.Type).Interface() {
	// 	//not the zero-value, stringify!
	// 	i.defstr = fmt.Sprintf("%v", def)
	// }
	//
	switch t {
	case "opt", "commalist", "spacelist":
		//options can also set short names
		var predictor complete.Predictor
		// if p, ok := o.val.Interface().(complete.Predictor); ok {
		// 	predictor = p
		// }
		predict := sf.Tag.Get("predict")
		if pable, ok := val.Interface().(complete.Predictor); ok {
			Log("impls predictable %T\n", pable)
			if predict != "" {
				panic("predict tag set on Predictable field " + i.name)
			}
			predictor = pable
		} else if val.CanAddr() {
			if pable, ok := val.Addr().Interface().(complete.Predictor); ok {
				Log("impls predictable %T\n", pable)
				if predict != "" {
					panic("predict tag set on Predictable field " + i.name)
				}
				predictor = pable
			}
		}
		if predictor == nil {
			Log("default predictable %T\n", val.Interface())
			switch {
			case predict == "":
				predictor = complete.PredictAnything
			case predict == "any":
				predictor = complete.PredictAnything
			case predict == "none":
				predictor = complete.PredictNothing
			case predict == "dirs":
				predictor = complete.PredictDirs("*")
			case predict == "files":
				predictor = complete.PredictFiles("*")
			case strings.HasPrefix(predict, "dirs:"):
				predictor = complete.PredictDirs(predict[len("dirs:"):])
			case strings.HasPrefix(predict, "files:"):
				predictor = complete.PredictFiles(predict[len("files:"):])
			default:
				panic("bad predict '" + predict + "'")
			}
		}
		o.completeCmd.Flags["--"+i.name] = predictor
		if short := sf.Tag.Get("short"); short != "" && short != "-" {
			if o.optnames[short] {
				o.errorf("option short name '%s' already in use", short)
				return nil
			} else {
				o.optnames[i.shortName] = true
				i.shortName = short
				o.completeCmd.Flags["-"+i.shortName] = predictor
			}
		}
		// log.Printf("define option: %s %s", name, sf.Type)
		o.opts = append(o.opts, i)
	case "arg":
		//TODO allow other types in 'arg' fields
		if sf.Type.Kind() != reflect.String {
			o.errorf("arg '%s' type must be a string", i.name)
			return nil
		}
		o.args = append(o.args, i)
	default:
		o.errorf("Invalid optype: %s", t)
	}
	return i
}

// Config returns the struct, eg the one passed to New, AddSubCmd or a sub command
func (o *Opts) Config() Config {
	return o.item.val.Interface()
}

// IsRunner returns true if the struct implements the Runner interface
func (o *Opts) IsRunner() bool {
	_, ok := o.Config().(Runner)
	return ok
}

// Run calls Run on the struct. If the struct isn't a runner the help is printed and program exists
// The error return is from the struct's run command.
func (o *Opts) Run() error {
	cmd := o.item.val.Interface()
	runner, ok := cmd.(Runner)
	if !ok {
		fmt.Fprintf(os.Stderr, "'%s' is not runable\n%s", o.name, o.Help())
		os.Exit(1)
	}
	return runner.Run()
}

//SubCmd add a subcmd and return new subcmd opt
func (o *Opts) SubCmd(name string, cmd Config) Builder {
	o.AddSubCmd(name, cmd)
	return o.GetSubCmd(name)
}

//AddSubCmd and return the orginal opts
func (o *Opts) AddSubCmd(name string, cmd Config) Builder {
	if o.arglist != nil {
		o.errorf("argslists and commands cannot be used together")
		panic(o.erred)
	}
	if _, exists := o.cmds[name]; exists {
		o.errorf("command already exists with name '%s'. Maybe in struct or dynamic", name)
		panic(o.erred)
	}
	val := reflect.ValueOf(cmd)
	if val.Kind() != reflect.Ptr {
		o.errorf("cmd must be a ptr. Got %T", cmd)
		panic(o.erred)
	}
	sub := fork(o, val)
	sub.name = name
	if helper, ok := cmd.(Helper); ok {
		sub.help = helper.Help()
	}
	o.cmds[name] = sub
	return o
}

//GetSubCmd used to get a dynamic subcommand inorder to add subcmds to it.
func (o *Opts) GetSubCmd(name string) Builder {
	sub, exists := o.cmds[name]
	if !exists {
		o.errorf("no command exists with name '%s'", name)
		panic(o.erred)
	}
	return sub
}

// Parent returns the parent opts
func (o *Opts) Parent() Builder {
	return o.parent
}

// Complete add --install and --uninstall to manage shell completions
// The name in the exec name this is completing, should only be used on a root opts
func (o *Opts) Complete(name string) Builder {
	Log("Complete '%s'\n", name)
	o.completeExec = name
	in := reflect.ValueOf(&o.internalOpts).Elem()
	// un/install
	itm1 := o.addOptArg(in.Type().Field(2), in.Field(2))
	itm2 := o.addOptArg(in.Type().Field(3), in.Field(3))
	if !o.optnames["i"] {
		o.optnames["i"] = true
		itm1.shortName = "i"
		o.completeCmd.Flags["-i"] = complete.PredictNothing
	}
	if !o.optnames["u"] {
		o.optnames["u"] = true
		itm2.shortName = "u"
		o.completeCmd.Flags["-u"] = complete.PredictNothing
	}
	return o
}

//Name sets the name of the program
func (o *Opts) Name(name string) Builder {
	o.name = name
	return o
}

//Version sets the version of the program
//and renders the 'version' template in the help text
func (o *Opts) Version(version string) Builder {
	//add version option
	g := reflect.ValueOf(&o.internalOpts).Elem()
	itm := o.addOptArg(g.Type().Field(1), g.Field(1))
	o.version = version
	if !o.optnames["v"] {
		o.optnames["v"] = true
		itm.shortName = "v"
		o.completeCmd.Flags["-v"] = complete.PredictNothing
	}
	return o
}

//Repo sets the repository link of the program
//and renders the 'repo' template in the help text
func (o *Opts) Repo(repo string) Builder {
	o.repo = repo
	return o
}

//PkgRepo infers the repository link of the program
//from the package import path of the struct (So note,
//this will not work for 'main' packages)
func (o *Opts) PkgRepo() Builder {
	if o.pkgrepo == "" {
		return o.errorf("Package repository could not be infered")
	}
	o.Repo(o.pkgrepo)
	return o
}

//Author sets the author of the program
//and renders the 'author' template in the help text
func (o *Opts) Author(author string) Builder {
	o.author = author
	return o
}

//PkgRepo infers the repository link of the program
//from the package import path of the struct (So note,
//this will not work for 'main' packages)
func (o *Opts) PkgAuthor() Builder {
	if o.pkgrepo == "" {
		return o.errorf("Package author could not be infered")
	}
	o.Author(o.pkgauthor)
	return o
}

//Set the padding width
func (o *Opts) SetPadWidth(p int) Builder {
	o.PadWidth = p
	return o
}

//Disable auto-padding
func (o *Opts) DisablePadAll() Builder {
	o.PadAll = false
	return o
}

//Set the line width (defaults to 72)
func (o *Opts) SetLineWidth(l int) Builder {
	o.LineWidth = l
	return o
}

//DocBefore inserts a text block before the specified template
func (o *Opts) DocBefore(target, newid, template string) Builder {
	return o.docOffset(0, target, newid, template)
}

//DocAfter inserts a text block after the specified template
func (o *Opts) DocAfter(target, newid, template string) Builder {
	return o.docOffset(1, target, newid, template)
}

func (o *Opts) docOffset(offset int, target, newid, template string) *Opts {
	if _, ok := o.templates[newid]; ok {
		o.errorf("new template already exists: %s", newid)
		return o
	}
	for i, id := range o.order {
		if id == target {
			o.templates[newid] = template
			index := i + offset
			rest := []string{newid}
			if index < len(o.order) {
				rest = append(rest, o.order[index:]...)
			}
			o.order = append(o.order[:index], rest...)
			return o
		}
	}
	o.errorf("target template not found: %s", target)
	return o
}

//DecSet replaces the specified template
func (o *Opts) DocSet(id, template string) Builder {
	if _, ok := DefaultTemplates[id]; !ok {
		if _, ok := o.templates[id]; !ok {
			o.errorf("template does not exist: %s", id)
			return o
		}
	}
	o.templates[id] = template
	return o
}

//ConfigPath defines a path to a JSON file which matches
//the structure of the provided config. Environment variables
//override JSON Config variables.
func (o *Opts) ConfigPath(path string) Builder {
	o.cfgPath = path
	return o
}

//UseEnv enables an implicit "env" struct tag option on
//all struct fields, the name of the field is converted
//into an environment variable with the transform
//`FooBar` -> `FOO_BAR`.
func (o *Opts) UseEnv() Builder {
	o.useEnv = true
	return o
}

//Parse with os.Args
func (o *Opts) Parse() Configured {
	return o.ParseArgs(os.Args[1:])
}

//ParseArgs with the provided arguments
func (o *Opts) ParseArgs(args []string) Configured {
	secondPass(o)
	Log("completeCom '%v'\n", o.completeCom)
	if o.completeCom != nil && o.completeCom.Complete() {
		Log("completion called \n")
		os.Exit(0)
	}
	var sub *Opts
	var err error
	if sub, _, err = o.process(args); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		_ = sub
		os.Exit(1)
	}
	return sub
}

//Process is the same as ParseArgs except
//it returns an error on failure
func (o *Opts) process(args []string) (*Opts, []string, error) {
	//cannot be processed - already encountered error - programmer error
	if o.erred != nil {
		return o, []string{}, fmt.Errorf("[opts] Process error: %s", o.erred)
	}
	//1. set config via JSON file
	if o.cfgPath != "" {
		b, err := ioutil.ReadFile(o.cfgPath)
		if err == nil {
			v := o.val.Interface() //*struct
			err = json.Unmarshal(b, v)
			if err != nil {
				o.erred = fmt.Errorf("Invalid config file: %s", err)
				return o, []string{}, errors.New(o.Help())
			}
		}
	}
	flagset := flag.NewFlagSet(o.name, flag.ContinueOnError)
	flagset.SetOutput(ioutil.Discard)
	//pre-loop through the options and
	//add shortnames and env names where possible
	for _, opt := range o.opts {
		env := camel2const(opt.name)
		if o.useEnv &&
			(opt.envName == "" || opt.envName == "!") &&
			opt.name != "help" && opt.name != "version" &&
			opt.name != "install-completetion" && opt.name != "uninstall-completetion" &&
			!o.envnames[env] {
			opt.envName = env
		}
	}
	for _, opt := range o.opts {
		// TODO remove debug
		// log.Printf("parse prepare option: %s", opt.name)
		//2. set config via environment
		envVal := ""
		if opt.useEnv || o.useEnv {
			envVal = os.Getenv(opt.envName)
		}
		if valset, ok := opt.val.Interface().(flag.Value); ok {
			// valset.Set(&valset, opt.name, "")
			// if opt.shortName != "" {
			// 	flagset.Var(&valset, opt.shortName, "")
			// }
			flagset.Var(valset, opt.name, "")
			if opt.shortName != "" {
				flagset.Var(valset, opt.shortName, "")
			}
		} else {
			//3. set config via Go's pkg/flags
			addr := opt.val.Addr().Interface()
			switch addr := addr.(type) {
			case flag.Value:
				flagset.Var(addr, opt.name, "")
				if opt.shortName != "" {
					flagset.Var(addr, opt.shortName, "")
				}
			case *[]string:
				sep := ""
				switch opt.typeName {
				case "commalist":
					sep = ","
				case "spacelist":
					sep = " "
				}
				fv := &sepList{sep: sep, strs: addr}
				flagset.Var(fv, opt.name, "")
				if opt.shortName != "" {
					flagset.Var(fv, opt.shortName, "")
				}
			case *bool:
				str2bool(envVal, addr)
				flagset.BoolVar(addr, opt.name, *addr, "")
				if opt.shortName != "" {
					flagset.BoolVar(addr, opt.shortName, *addr, "")
				}
			case *string:
				str2str(envVal, addr)
				flagset.StringVar(addr, opt.name, *addr, "")
				if opt.shortName != "" {
					flagset.StringVar(addr, opt.shortName, *addr, "")
				}
			case *int:
				str2int(envVal, addr)
				flagset.IntVar(addr, opt.name, *addr, "")
				if opt.shortName != "" {
					flagset.IntVar(addr, opt.shortName, *addr, "")
				}
			case *time.Duration:
				flagset.DurationVar(addr, opt.name, *addr, "")
				if opt.shortName != "" {
					flagset.DurationVar(addr, opt.shortName, *addr, "")
				}
			default:
				return o, []string{}, fmt.Errorf("[opts] Option '%s' has unsupported type - '%T'\n", opt.name, addr)
			}
		}
	}
	// log.Printf("parse %+v", args)
	//set user config
	err := flagset.Parse(args)
	if err != nil {
		//insert flag errors into help text
		o.erred = err
		o.internalOpts.Help = true
	}
	//internal opts (--help and --version)
	if o.internalOpts.Help {
		fmt.Println(o.Help())
		os.Exit(0)
	}
	if o.internalOpts.Version {
		fmt.Println(o.version)
		os.Exit(0)
	}
	if o.internalOpts.InstallCompletetion {
		err := install.Install(o.completeExec)
		if err != nil {
			fmt.Println("error installing " + err.Error())
			os.Exit(1)
		}
		fmt.Println("installed")
		os.Exit(0)
	}
	if o.internalOpts.UninstallCompletetion {
		err := install.Uninstall(o.completeExec)
		if err != nil {
			fmt.Println("error uninstalling " + err.Error())
			os.Exit(1)
		}
		fmt.Println("uninstalled")
		os.Exit(0)
	}
	//fill each individual arg
	args = flagset.Args()
	for i, argument := range o.args {
		if len(args) > 0 {
			str := args[0]
			args = args[1:]
			argument.val.SetString(str)
		} else if argument.defstr == "" {
			//not-set and no default!
			o.erred = fmt.Errorf("Argument #%d '%s' has no default value", i+1, argument.name)
			return o, []string{}, errors.New(o.Help())
		}
	}
	//use command? peek at args
	if len(o.cmds) > 0 && len(args) > 0 {
		a := args[0]
		//matching command, use it
		if sub, exists := o.cmds[a]; exists {
			subo, cmdnames, err := sub.process(args[1:])
			//user wants name to be set
			if o.cmdname != nil {
				cmdname := a
				for i := len(cmdnames) - 1; i >= 0; i-- {
					cmdname = cmdname + "." + cmdnames[i]
				}
				o.cmdname.SetString(cmdname)
			}
			return subo, append(cmdnames, a), err
		}
	}
	//fill arglist? assign remaining as slice
	if o.arglist != nil {
		if len(args) < o.arglist.min {
			o.erred = fmt.Errorf("Too few arguments (expected %d, got %d)", o.arglist.min, len(args))
			return o, []string{}, errors.New(o.Help())
		}
		o.arglist.val.Set(reflect.ValueOf(args))
		args = nil
	}
	//we *should* have consumed all args at this point.
	//this prevents:  ./foo --bar 42 -z 21 ping --pong 7
	//where --pong 7 is ignored
	if len(args) != 0 {
		o.erred = fmt.Errorf("Unexpected arguments: %+v", args)
		return o, []string{}, errors.New(o.Help())
	}
	return o, []string{}, nil
}
