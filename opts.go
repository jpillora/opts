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
)

//Opts is the main class, it contains
//all parsing state for a single set of
//arguments
type Opts struct {
	//embed item since an Opts can also be an item
	item
	parent       *Opts
	subcmds      map[string]*Opts
	opts         []*item
	args         []*item
	arglist      *argumentlist
	shorts       map[string]bool
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
	//public format settings
	LineWidth int  //42
	PadAll    bool //true
	PadWidth  int  //2
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
	name      string
	shortName string
	envName   string
	useEnv    bool
	typeName  string
	help      string
	hasDef    bool
}

//Creates a new Opts instance and Parses it
func Parse(config interface{}) *Opts {
	return New(config).Parse()
}

//New creates a new Opts instance
func New(config interface{}) *Opts {
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

func fork(parent *Opts, c reflect.Value) *Opts {
	//copy default ordering
	var order []string

	//root only
	if parent == nil {
		order = make([]string, len(DefaultOrder))
		copy(order, DefaultOrder)
	}

	//instantiate
	o := &Opts{
		item: item{
			val: c,
		},
		parent:    parent,
		shorts:    map[string]bool{},
		subcmds:   map[string]*Opts{},
		opts:      []*item{},
		order:     order,
		templates: map[string]string{},
		//public defaults
		LineWidth: 72,
		PadAll:    true,
		PadWidth:  2,
	}

	t := c.Type()
	k := t.Kind()
	//must be pointer (meaningless to modify a copy of the struct)
	if k != reflect.Ptr {
		o.errorf("opts: %s should be a pointer to a struct", t.Name())
		return o
	}

	c = c.Elem()
	t = c.Type()
	k = t.Kind()
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
		k := sf.Type.Kind()
		switch k {
		case reflect.Ptr, reflect.Struct:
			o.addSubcmd(sf, val)
		case reflect.Slice:
			if sf.Type.Elem().Kind() != reflect.String {
				o.errorf("arglist must be of type []string")
				return o
			}
			o.addArgs(sf, val)
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
		default:
			o.errorf("Struct field '%s' has unsupported type: %s", sf.Name, k)
			return o
		}
	}

	//add help option
	g := reflect.ValueOf(&o.internalOpts).Elem()
	o.addOptArg(g.Type().Field(0), g.Field(0))

	return o
}

func (o *Opts) errorf(format string, args ...interface{}) *Opts {
	//only store the first error
	if o.erred == nil {
		o.erred = fmt.Errorf(format, args...)
	}
	return o
}

func (o *Opts) addSubcmd(sf reflect.StructField, val reflect.Value) {

	if o.arglist != nil {
		o.errorf("argslists and subcommands cannot be used together")
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
	// log.Printf("define subcmd: %s =====", subname)
	sub := fork(o, val)
	sub.name = name
	sub.help = sf.Tag.Get("help")
	o.subcmds[name] = sub
}

func (o *Opts) addArgs(sf reflect.StructField, val reflect.Value) {

	if len(o.subcmds) > 0 {
		o.errorf("argslists and subcommands cannot be used together")
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
			name: name,
			help: sf.Tag.Get("help"),
		},
		min: min,
	}
}

var durationType = reflect.TypeOf(time.Second)

func (o *Opts) addOptArg(sf reflect.StructField, val reflect.Value) {

	i := &item{val: val}

	//find name
	i.name = sf.Tag.Get("name")
	if i.name == "" {
		i.name = camel2dash(sf.Name) //default to struct field name
	}

	//assume int64s are durations
	if sf.Type.Kind() == reflect.Int64 &&
		!sf.Type.AssignableTo(durationType) {
		o.errorf("int64 field '%s' must be of type time.Duration", i.name)
		return
	}

	i.envName = sf.Tag.Get("env")
	if i.envName != "" {
		i.useEnv = true
	}
	if i.envName == "" || i.envName == "!" {
		i.envName = camel2const(sf.Name)
	}

	//assume opt, unless arg tag is present
	t := sf.Tag.Get("type")
	if t == "" {
		t = "opt"
	}

	//get help text
	i.help = sf.Tag.Get("help")

	//display default, when non-zero-val
	i.hasDef = false
	def := val.Interface()
	if def != reflect.Zero(sf.Type).Interface() {
		i.hasDef = true
		if i.help == "" {
			i.help = fmt.Sprintf("default %v", def)
		} else {
			i.help += fmt.Sprintf(" (default %v)", def)
		}
	}

	switch t {
	case "opt":
		//only options have short names
		n := i.name[0:1]
		if _, ok := o.shorts[n]; ok {
			n = ""
		} else {
			o.shorts[n] = true
		}
		i.shortName = n

		// log.Printf("define option: %s %s", name, sf.Type)
		o.opts = append(o.opts, i)
	case "arg":
		//TODO allow other types in 'arg' fields
		if sf.Type.Kind() != reflect.String {
			o.errorf("arg '%s' type must be a string", i.name)
			return
		}
		o.args = append(o.args, i)
	default:
		o.errorf("Invalid optype: %s", t)
	}
}

func (o *Opts) Name(name string) *Opts {
	o.name = name
	return o
}

func (o *Opts) Version(version string) *Opts {
	//add version option
	g := reflect.ValueOf(&o.internalOpts).Elem()
	o.addOptArg(g.Type().Field(1), g.Field(1))
	o.version = version
	return o
}

func (o *Opts) PkgRepo() *Opts {
	if o.pkgrepo == "" {
		return o.errorf("Package repository could not be infered")
	}
	o.Repo(o.pkgrepo)
	return o
}

func (o *Opts) Repo(repo string) *Opts {
	o.repo = repo
	return o
}

func (o *Opts) PkgAuthor() *Opts {
	if o.pkgrepo == "" {
		return o.errorf("Package author could not be infered")
	}
	o.Author(o.pkgauthor)
	return o
}

func (o *Opts) Author(author string) *Opts {
	o.author = author
	return o
}

// //Doc inserts a text block at the end of the help text
// func (o *Opts) Doc(paragraph string) *Opts {
// 	return o
// }

// //DocAfter inserts a text block after the specified help entry
// func (o *Opts) DocAfter(id, paragraph string) *Opts {
// 	return o
// }

// //DecSet replaces the specified
// func (o *Opts) DocSet(id, paragraph string) *Opts {
// 	return o
// }

//ConfigPath defines a path to a JSON file which matches
//the structure of the provided config. Environment variables
//override JSON Config variables.
func (o *Opts) ConfigPath(path string) *Opts {
	o.cfgPath = path
	return o
}

//ConfigPath defines a path to a JSON file which matches
//the structure of the provided config. Environment variables
//override JSON Config variables.
func (o *Opts) UseEnv() *Opts {
	o.useEnv = true
	return o
}

//Parse with os.Args
func (o *Opts) Parse() *Opts {
	return o.ParseArgs(os.Args[1:])
}

//ParseArgs with the provided arguments
func (o *Opts) ParseArgs(args []string) *Opts {
	if err := o.Process(args); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
	return o
}

//Process is the same as ParseArgs except
//it returns an error on failure
func (o *Opts) Process(args []string) error {

	//cannot be processed - already encountered error
	if o.erred != nil {
		return o.erred
	}

	//1. set config via JSON file
	if o.cfgPath != "" {
		b, err := ioutil.ReadFile(o.cfgPath)
		if err == nil {
			v := o.val.Interface() //*struct
			err = json.Unmarshal(b, v)
			if err != nil {
				return fmt.Errorf("Invalid config file: %s", err)
			}
		}
	}

	flagset := flag.NewFlagSet(o.name, flag.ContinueOnError)
	flagset.SetOutput(ioutil.Discard)

	for _, opt := range o.opts {
		// TODO remove debug
		// log.Printf("parse prepare option: %s", opt.name)

		//2. set config via environment
		envVal := ""
		if opt.useEnv || o.useEnv {
			envVal = os.Getenv(opt.envName)
		}
		//3. set config via Go's pkg/flags
		addr := opt.val.Addr().Interface()
		switch addr := addr.(type) {
		case flag.Value:
			flagset.Var(addr, opt.name, "")
			if opt.shortName != "" {
				flagset.Var(addr, opt.shortName, "")
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
			return fmt.Errorf("Option '%s' has unsupported type", opt.name)
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

	//internal opts
	if o.internalOpts.Help {
		return errors.New(o.Help())
	} else if o.internalOpts.Version {
		fmt.Println(o.version)
		os.Exit(0)
	}

	//fill each individual arg
	args = flagset.Args()
	for i, argument := range o.args {
		if len(args) > 0 {
			str := args[0]
			args = args[1:]
			argument.val.SetString(str)
		} else if !argument.hasDef {
			//not-set and no default!
			return fmt.Errorf("Argument #%d '%s' is missing", i+1, argument.name)
		}
	}

	//use subcommand? peek at args
	if len(o.subcmds) > 0 && len(args) > 0 {
		a := args[0]
		//matching subcommand, use it
		if sub, exists := o.subcmds[a]; exists {
			//user wants name to be set
			if o.cmdname != nil {
				o.cmdname.SetString(a)
			}
			return sub.Process(args[1:])
		}
	}

	//fill arglist? assign remaining as slice
	if o.arglist != nil {
		if len(args) < o.arglist.min {
			return fmt.Errorf("Too few arguments (expected %d, got %d)", o.arglist.min, len(args))
		}
		o.arglist.val.Set(reflect.ValueOf(args))
		args = nil
	}

	//we *should* have consumed all args
	if len(args) != 0 {
		return fmt.Errorf("Unexpected arguments: %+v", args)
	}

	return nil
}
