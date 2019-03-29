package opts

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

//Parse with os.Args
func (n *node) Parse() ParsedOpts {
	return n.ParseArgs(os.Args[1:])
}

//ParseArgs with the provided arguments
func (n *node) ParseArgs(args []string) ParsedOpts {
	if err := n.parse(args); err != nil {
		msg := ""
		if oe, ok := err.(*optsError); ok {
			log.Printf("OPTS-ERR: %s", err)
			msg = oe.Error() + "\n"
		} else {
			log.Printf("USER-ERR: %s", err)
			msg = n.Help()
		}
		fmt.Fprintf(os.Stderr, msg)
		os.Exit(1)
	}
	return n
}

//IsRunnable
func (n *node) IsRunnable() bool {
	//TODO
	return false
}

//Run the parsed configuration
func (n *node) Run() error {
	//TODO
	return nil
}

//Run the parsed configuration
func (n *node) RunFatal() {
	//TODO
}

//parse is the same as ParseArgs except
//it returns an error on failure
func (n *node) parse(args []string) error {
	if n.err != nil {
		return n.err
	}
	v := n.item.val
	//all fields from val
	if v.Type().Kind() != reflect.Ptr && v.Type().Elem().Kind() != reflect.Struct {
		return n.errorf("%s should be a pointer to a struct", v.Type().Name())
	}
	//add struct
	if err := n.addStructFields(v.Elem()); err != nil {
		return err
	}
	//add help flag
	n.setHelpFlag()
	//find defaults from config's package
	n.setPkgDefaults()
	//1. set config via JSON file, unmarshal it into the struct
	if n.cfgPath != "" {
		b, err := ioutil.ReadFile(n.cfgPath)
		if err == nil {
			v := n.val.Interface() //*struct
			err = json.Unmarshal(b, v)
			if err != nil {
				return fmt.Errorf("Invalid config file: %s", err)
			}
		}
	}
	//pre-loop through the options and
	//add shortnames and env names where possible
	for _, opt := range n.flags {
		//should generate shortname?
		if len(opt.name) >= 3 && opt.shortName == "" {
			//not already taken?
			if s := opt.name[0:1]; !n.optnames[s] {
				opt.shortName = s
				n.optnames[s] = true
			}
		}
		env := camel2const(opt.name)
		if n.useEnv && (opt.envName == "" || opt.envName == "!") &&
			opt.name != "help" && opt.name != "version" &&
			!n.envnames[env] {
			opt.envName = env
		}
	}
	//link each flag to fields in the underlying struct
	flagset := flag.NewFlagSet(n.name, flag.ContinueOnError)
	flagset.SetOutput(ioutil.Discard)
	for _, opt := range n.flags {
		//2. set config via environment
		envVal := ""
		if opt.useEnv || n.useEnv {
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
			return fmt.Errorf("[opts] Option '%s' has unsupported type", opt.name)
		}
	}
	//set user config
	err := flagset.Parse(args)
	if err != nil {
		//insert flag errors into help text
		n.err = err
		n.internalOpts.Help = true
	}
	//internal opts (--help and --version)
	if n.internalOpts.Help {
		fmt.Println(n.Help())
		os.Exit(0)
	} else if n.internalOpts.Version {
		fmt.Println(n.version)
		os.Exit(0)
	}
	//fill each individual arg
	args = flagset.Args()
	for i, argument := range n.args {
		if len(args) > 0 {
			str := args[0]
			args = args[1:]
			argument.val.SetString(str)
		} else if argument.defstr == "" {
			//not-set and no default!
			return fmt.Errorf("Argument #%d '%s' has no default value", i+1, argument.name)
		}
	}
	//use command? peek at args
	if len(n.cmds) > 0 && len(args) > 0 {
		a := args[0]
		//matching command, use it
		if sub, exists := n.cmds[a]; exists {
			//user wants command name to be set on their struct?
			if n.cmdname != nil {
				n.cmdname.SetString(a)
			}
			//recurse!
			return sub.parse(args[1:])
		}
	}
	//fill arglist? assign remaining as slice
	if n.arglist != nil {
		if len(args) < n.arglist.min {
			return fmt.Errorf("Too few arguments (expected %d, got %d)", n.arglist.min, len(args))
		}
		n.arglist.val.Set(reflect.ValueOf(args))
		args = nil
	}
	//we *should* have consumed all args at this point.
	//this prevents:  ./foo --bar 42 -z 21 ping --pong 7
	//where --pong 7 is ignored
	if len(args) != 0 {
		return fmt.Errorf("Unexpected arguments: %+v", args)
	}
	return nil
}

func (n *node) addStructFields(c reflect.Value) error {
	t := c.Type()
	k := t.Kind()
	//deref pointer
	if k == reflect.Ptr {
		c = c.Elem()
		t = c.Type()
		k = t.Kind()
	}
	if k != reflect.Struct {
		return n.errorf("opts: %s should be a pointer to a struct (got %s)", t.Name(), k)
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
		//is a pkg/flag type
		k := sf.Type.Kind()
		if sf.Type.Implements(flagValueType) {
			err := n.addOptArg(sf, val)
			if err != nil {
				return err
			}
			continue
		}
		//reflect to find flag type
		var err error
		switch k {
		case reflect.Ptr, reflect.Struct:
			if sf.Tag.Get("type") == "embedded" {
				err = n.addStructFields(val)
			} else {
				err = n.addCmd(sf, val)
			}
		case reflect.Slice:
			if sf.Type.Elem().Kind() != reflect.String {
				err = n.errorf("slice (list) types must be []string")
			} else if sf.Tag.Get("type") == "commalist" {
				err = n.addOptArg(sf, val)
			} else if sf.Tag.Get("type") == "spacelist" {
				err = n.addOptArg(sf, val)
			} else {
				err = n.addArgList(sf, val)
			}
		case reflect.Bool, reflect.String, reflect.Int, reflect.Int64:
			if sf.Tag.Get("type") == "cmdname" {
				if k != reflect.String {
					err = n.errorf("cmdname field '%s' must be a string", sf.Name)
				} else {
					n.cmdname = &val
				}
			} else {
				err = n.addOptArg(sf, val)
			}
		case reflect.Interface:
			err = n.errorf("Struct field '%s' interface type must implement flag.Value", sf.Name)
		default:
			err = n.errorf("Struct field '%s' has unsupported type: %s", sf.Name, k)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *node) addCmd(sf reflect.StructField, val reflect.Value) error {
	if n.arglist != nil {
		return n.errorf("argslists and commands cannot be used together")
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
	// log.Printf("define cmd: %s =====", subname)
	sub := newNode(val)
	sub.name = name
	sub.help = sf.Tag.Get("help")
	sub.parent = n
	n.cmds[name] = sub
	return nil
}

func (n *node) addArgList(sf reflect.StructField, val reflect.Value) error {
	if len(n.cmds) > 0 {
		return n.errorf("argslists and commands cannot be used together")
	}
	if n.arglist != nil {
		return n.errorf("only 1 arglist field is allowed ('%s' already defined)", n.arglist.name)
	}
	name := sf.Tag.Get("name")
	if name == "" || name == "!" {
		name = camel2dash(sf.Name) //default to struct field name
	}
	if val.Len() != 0 {
		return n.errorf("arglist '%s' is required so it should not be set. "+
			"If you'd like to set a default, consider using an option instead.", name)
	}
	min, _ := strconv.Atoi(sf.Tag.Get("min"))
	//insert (there can only be one)
	n.arglist = &argumentlist{
		item: item{
			val:  val,
			name: name,
			help: sf.Tag.Get("help"),
		},
		min: min,
	}
	return nil
}

func (n *node) addOptArg(sf reflect.StructField, val reflect.Value) error {
	//assume opt, unless arg tag is present
	t := sf.Tag.Get("type")
	if t == "" {
		t = "opt"
	}
	i := &item{
		val:      val,
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
		if n.envnames[i.envName] {
			return n.errorf("option env name '%s' already in use", i.name)
		}
		n.envnames[i.envName] = true
		i.useEnv = true
	}
	//opt names cannot clash with each other
	if n.optnames[i.name] {
		return n.errorf("option name '%s' already in use", i.name)
	}
	n.optnames[i.name] = true
	//get help text
	i.help = sf.Tag.Get("help")
	//the **displayed** default, use 'default' tag, otherwise infer
	defstr := sf.Tag.Get("default")
	if defstr != "" {
		i.defstr = defstr
	} else if val.Kind() == reflect.Slice {
		i.defstr = ""
	} else if def := val.Interface(); def != reflect.Zero(sf.Type).Interface() {
		//not the zero-value, stringify!
		i.defstr = fmt.Sprintf("%v", def)
	}
	switch t {
	case "opt", "commalist", "spacelist":
		//options can also set short names
		if short := sf.Tag.Get("short"); short != "" {
			if n.optnames[short] {
				return n.errorf("option short name '%s' already in use", short)
			}
			n.optnames[i.shortName] = true
			i.shortName = short
		}
		// log.Printf("define option: %s %s", name, sf.Type)
		n.flags = append(n.flags, i)
	case "arg":
		//TODO allow other types in 'arg' fields
		if sf.Type.Kind() != reflect.String {
			return n.errorf("arg '%s' type must be a string", i.name)
		}
		n.args = append(n.args, i)
	default:
		return n.errorf("Invalid optype: %s", t)
	}
	// log.Printf("add '%s' (%s)", i.name, sf.Type)
	return nil
}

func (n *node) setHelpFlag() error {
	g := reflect.ValueOf(&n.internalOpts).Elem()
	t, _ := g.Type().FieldByName("Help")
	v := g.FieldByName("Help")
	if err := n.addOptArg(t, v); err != nil {
		return n.errorf("error adding internal --help flag: %s - please report issue", err)
	}
	return nil
}

func (n *node) setPkgDefaults() {
	//attempt to infer package name, repo, author
	pkgpath := n.item.val.Elem().Type().PkgPath()
	parts := strings.Split(pkgpath, "/")
	if len(parts) >= 3 {
		if n.pkgauthor == "" {
			n.pkgauthor = parts[1]
		}
		if n.name == "" {
			n.Name(parts[2])
		}
		if n.pkgrepo == "" {
			switch parts[0] {
			case "github.com", "bitbucket.org":
				n.pkgrepo = "https://" + strings.Join(parts[0:3], "/")
			}
		}
	}
}
