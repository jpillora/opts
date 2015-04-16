package flag

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
)

//Flag is the main class, it contains
//all parsing state for a single set of
//arguments
type Flag struct {
	config     reflect.Value
	parent     *Flag
	subs       map[string]*Flag
	opts       []*option
	shorts     map[string]bool
	globalOpts struct {
		Help, Version bool
	}
	erred   error
	cmdname *reflect.Value

	name, version string
	repo, author  string
	//public format settings
	LineWidth int  //42
	PadAll    bool //true
	PadWidth  int  //2
	Templates map[string]string
}

//option is the structure representing a
//single flag, exposed for templating
type option struct {
	val         reflect.Value
	name        string
	shortName   string
	displayName string //calculated
	typeName    string
	help        string
}

//New creates a new Flag instance
func New(config interface{}) *Flag {

	v := reflect.ValueOf(config)

	//nil parent -> root command
	f := fork(nil, v)

	//attempt to infer package name and author
	parts := strings.Split(v.Type().PkgPath(), "/")
	if len(parts) >= 2 {
		f.Name(parts[len(parts)-1])
		f.Author(parts[len(parts)-2])
	}

	return f
}

func fork(parent *Flag, c reflect.Value) *Flag {

	//instantiate
	f := &Flag{
		config: c,
		parent: parent,
		shorts: map[string]bool{},
		subs:   map[string]*Flag{},
		opts:   []*option{},
		//public defaults
		LineWidth: 42,
		PadAll:    true,
		PadWidth:  2,
		Templates: map[string]string{},
	}

	t := c.Type()
	k := t.Kind()

	//must be pointer (meaningless to modify a copy of the struct)
	if k != reflect.Ptr {
		f.errorf("flag.New(config): config should be a pointer (%s) to a struct", k)
		return f
	}

	c = c.Elem()
	t = c.Type()
	k = t.Kind()

	if k != reflect.Struct {
		f.errorf("flag.New(config): config should be a pointer to a struct (%s)", k)
		return f
	}

	//parse struct fields
	for i := 0; i < c.NumField(); i++ {
		val := c.Field(i)
		sf := t.Field(i)
		switch sf.Type.Kind() {
		case reflect.Ptr, reflect.Struct:
			f.addSubcmd(sf, val)
		case reflect.Bool, reflect.String, reflect.Int:
			if sf.Tag.Get("cmd") != "" {
				f.cmdname = &val
			} else if sf.Tag.Get("arg") != "" {
				f.addArgument(sf, val)
			} else {
				f.addOption(sf, val)
			}

		default:
			f.errorf("Struct field '%s' has unsupported type: %s",
				sf.Name, sf.Type.Kind().String())
			return f
		}
	}

	//add help option
	g := reflect.ValueOf(&f.globalOpts).Elem()
	f.addOption(g.Type().Field(0), g.Field(0))

	return f
}

func (f *Flag) errorf(format string, args ...interface{}) {
	f.erred = fmt.Errorf(format, args...)
}

func (f *Flag) addSubcmd(sf reflect.StructField, val reflect.Value) {
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

	subname := camel2dash(sf.Name)
	// log.Printf("define subcmd: %s =====", subname)
	sub := fork(f, val)
	sub.name = subname
	f.subs[subname] = sub
}

func (f *Flag) addOption(sf reflect.StructField, val reflect.Value) {

	name := camel2dash(sf.Name)

	n := name[0:1]
	if _, ok := f.shorts[n]; ok {
		n = ""
	} else {
		f.shorts[n] = true
	}

	help := sf.Tag.Get("help")

	//display default when set
	def := val.Interface()
	if def != reflect.Zero(sf.Type).Interface() {
		help += fmt.Sprintf(" (default %v)", def)
	}

	// log.Printf("define option: %s %s", name, sf.Type)

	f.opts = append(f.opts, &option{
		val:       val,
		name:      name,
		shortName: n,
		help:      help,
	})
}

func (f *Flag) addArgument(sf reflect.StructField, val reflect.Value) {

}

func (f *Flag) Name(name string) *Flag {
	f.name = name
	return f
}

func (f *Flag) Version(version string) *Flag {
	//add version option
	g := reflect.ValueOf(&f.globalOpts).Elem()
	f.addOption(g.Type().Field(1), g.Field(1))

	f.version = version
	return f
}

func (f *Flag) Repo(repo string) *Flag {
	f.repo = repo
	return f
}

func (f *Flag) Author(author string) *Flag {
	f.author = author
	return f
}

//Parse with os.Args
func (f *Flag) Parse() *Flag {
	return f.ParseArgs(os.Args[1:])
}

//ParseArgs with the provided arguments
func (f *Flag) ParseArgs(args []string) *Flag {
	if err := f.Process(args); err != nil {
		log.Fatal(err)
	}
	return f
}

//Process is the same as ParseArgs except
//it returns an error instead of calling log.Fatal
func (f *Flag) Process(args []string) error {

	//already errored
	if f.erred != nil {
		return f.erred
	}

	//use this command
	flagset := flag.NewFlagSet(f.name, flag.ContinueOnError)
	flagset.Usage = func() {
		fmt.Fprint(os.Stdout, f.Help())
	}

	for _, opt := range f.opts {
		// log.Printf("parse prepare option: %s", opt.name)
		//take address, not value
		addr := opt.val.Addr().Interface()
		switch addr := addr.(type) {
		case *bool:
			flagset.BoolVar(addr, opt.name, *addr, "")
			if opt.shortName != "" {
				flagset.BoolVar(addr, opt.shortName, *addr, "")
			}
		case *string:
			flagset.StringVar(addr, opt.name, *addr, "")
			if opt.shortName != "" {
				flagset.StringVar(addr, opt.shortName, *addr, "")
			}
		case *int:
			flagset.IntVar(addr, opt.name, *addr, "")
			if opt.shortName != "" {
				flagset.IntVar(addr, opt.shortName, *addr, "")
			}
		default:
			panic("Unexpected logic error")
		}
	}

	// log.Printf("parse %+v", args)
	//set user config
	err := flagset.Parse(args)
	if err != nil {
		return err
	}

	if f.globalOpts.Help {
		flagset.Usage()
	}

	if f.globalOpts.Version {
		fmt.Println(f.version)
		os.Exit(0)
	}

	args = flagset.Args()

	//peek at args, maybe use subcommand
	if len(args) > 0 {
		a := args[0]
		//matching subcommand, use it
		if sub, exists := f.subs[a]; exists {
			//user wants name to be set
			if f.cmdname != nil {
				f.cmdname.SetString(a)
			}
			return sub.Process(args[1:])
		}
	}

	return nil
}
