package opts

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
)

//Parse with os.Args
func (n *node) Parse() ParsedOpts {
	return n.ParseArgs(os.Args[1:])
}

//ParseArgs with the provided arguments
func (n *node) ParseArgs(args []string) ParsedOpts {
	//shell-completion?
	if n.complete {
		if cl := os.Getenv("COMP_LINE"); cl != "" {
			args := strings.Split(cl, " ")
			n.parse(args[1:]) //ignore error
			if ok := n.doCompletion(); !ok {
				os.Exit(1)
			}
			os.Exit(0)
		}
	}
	//use built state to perform parse
	if err := n.parse(args); err != nil {
		//expected exit (0) print message as-is
		if ee, ok := err.(*exitError); ok {
			fmt.Fprintf(os.Stderr, ee.msg)
			os.Exit(0)
		}
		//unexpected exit (1) print message to programmer
		if ae, ok := err.(*authorError); ok {
			fmt.Fprintf(os.Stderr, "opts usage error: %s\n", ae.err)
			os.Exit(1)
		}
		//unexpected exit (1) embed message in help to user
		n.err = err
		fmt.Fprintf(os.Stderr, n.Help())
		os.Exit(1)
	}
	//success
	return n
}

//parse is the same as ParseArgs except
//it returns an error on failure
func (n *node) parse(args []string) error {
	//return the stored error
	if n.err != nil {
		return n.err
	}
	//find default name for root-node
	if n.name == "" && n.parent == nil {
		if exe, err := os.Executable(); err == nil {
			_, n.name = path.Split(exe)
		}
	}
	//root value
	v := n.item.val
	//all fields from val
	if v.Type().Kind() != reflect.Ptr && v.Type().Elem().Kind() != reflect.Struct {
		return n.errorf("%s should be a pointer to a struct", v.Type().Name())
	}
	//add this node and its fields (recurses if has sub-commands)
	if err := n.addStructFields(v.Elem()); err != nil {
		return err
	}
	//add help flag
	n.addInternalFlags()
	//add user provided flagsets, will error if there is a naming collision
	if err := n.addFlagsets(); err != nil {
		return err
	}
	//find defaults from config's package
	n.setPkgDefaults()
	//first, set config via JSON file, unmarshal it into the struct
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
	//add shortnames where possible
	for _, flag := range n.flags {
		if flag.shortName == "" && len(flag.name) >= 3 {
			if s := flag.name[0:1]; !n.flagNames[s] {
				flag.shortName = s
				n.flagNames[s] = true
			}
		}
	}
	//link each flag to fields in the underlying struct
	flagset := flag.NewFlagSet(n.name, flag.ContinueOnError)
	flagset.SetOutput(ioutil.Discard)
	for _, item := range n.flags {
		flagset.Var(item, item.name, "")
		if sn := item.shortName; sn != "" {
			flagset.Var(item, sn, "")
		}
	}
	if err := flagset.Parse(args); err != nil {
		//insert flag errors into help text
		n.err = err
		n.internalOpts.Help = true
	}
	//internal opts
	if n.internalOpts.Help {
		return &exitError{n.Help()}
	} else if n.internalOpts.Version {
		return &exitError{n.version}
	} else if n.internalOpts.Install {
		return n.manageCompletion(false)
	} else if n.internalOpts.Uninstall {
		return n.manageCompletion(true)
	}
	//fill each individual arg
	args = flagset.Args()
	for i, item := range n.args {
		s := ""
		if len(args) > 0 {
			s = args[0]
			args = args[1:]
		}
		if s == "" {
			return fmt.Errorf("Argument '%s' (#%d) is missing", item.name, i+1)
		}
		if err := item.Set(s); err != nil {
			return fmt.Errorf("Argument '%s' (#%d) is invalid: %s", item.name, i+1, err)
		}
	}
	//use command? peek at args
	if len(n.cmds) > 0 && len(args) > 0 {
		a := args[0]
		//matching command, use it
		if sub, exists := n.cmds[a]; exists {
			//store matched command
			n.cmd = sub
			//user wants command name to be set on their struct?
			if n.cmdname != nil {
				*n.cmdname = a
			}
			//recurse!
			return sub.parse(args[1:])
		}
	}
	//fill arglist? assign remaining as slice
	// if n.arglist != nil {
	// 	if len(args) < n.arglist.min {
	// 		return fmt.Errorf("Too few arguments (expected %d, got %d)", n.arglist.min, len(args))
	// 	}
	// 	n.arglist.val.Set(reflect.ValueOf(args))
	// 	args = nil
	// }
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
		sf := t.Field(i)
		val := c.Field(i)
		//add one field
		if err := n.addStructField(sf, val); err != nil {
			return fmt.Errorf("field '%s': %s", sf.Name, err)
		}
	}
	return nil
}

func (n *node) addStructField(sf reflect.StructField, val reflect.Value) error {
	kv := newKV(sf.Tag.Get("opts"))
	help := sf.Tag.Get("help")
	typeName := sf.Tag.Get("type") 
	log.Printf("add field: %s", sf.Name)
	if err := n.addKVField(kv, help, sf.Name, typeName, val); err != nil {
		return err
	}
	if ks := kv.keys(); len(ks) > 0 {
		return fmt.Errorf("unused opts keys: %v", ks)
	}
	return nil
}

func (n *node) addKVField(kv *kv, fName, help, typeName string, val reflect.Value) error {
	//ignore unaddressed unexported fields
	if !val.CanSet() {
		return nil
	}
	//deref pointer
	if val.Type().Kind() == reflect.Ptr {
		val = val.Elem()
	}
	//parse key-values
	//ignore `opts:"-"`
	if _, ok := kv.take("-"); ok {
		return nil
	}
	//get field name and type
	name, _ := kv.take("name")
	if name == "" {
		//default to struct field name
		name = camel2dash(fName)
		//slice? use singular, since usage is --file foo --file bar
		if val.Type().Kind() == reflect.Slice {
			name = strings.TrimSuffix(name, "s")
		}
	}
	//new kv type defs supercede legacy defs
	if t, ok := kv.take("type"); ok {
		typeName = t
	}
	//default opts type from go type
	if typeName == "" {
		switch val.Type().Kind() {
		case reflect.Struct:
			typeName = "embedded"
		default:
			typeName = "flag"
		}
	}
	if typeName == "embedded" {
		//recurse!
		return n.addStructFields(val)
	}
	if typeName == "cmdname" {
		return n.setCmdName(val)
	}
	//new kv help defs supercede legacy defs
	if h, ok := kv.take("help"); ok {
		help = h
	}
	//inline sub-command
	if typeName == "cmd" {
		return n.addInlineCmd(name, help, val)
	}
	//from this point, we must have a flag or an arg
	i, err := newItem(val)
	if err != nil {
		return err
	}
	i.typeName = typeName
	i.name = name
	i.help = help
	//set default text
	if d, ok := kv.take("default"); ok {
		i.defstr = d
	} else if !i.slice {
		v := val.Interface()
		zero := v == reflect.Zero(val.Type()).Interface()
		if !zero {
			i.defstr = fmt.Sprintf("%v", v)
		}
	}
	//set env var name to use
	if e, ok := kv.take("env"); ok || n.useEnv {
		explicit := true
		if e == "" {
			explicit = false
			e = camel2const(i.name)
		}
		_, set := n.envNames[e]
		if set && explicit {
			return n.errorf("env name '%s' already in use", e)
		}
		if !set {
			n.envNames[e] = true
			i.envName = e
			i.useEnv = true
		}
	}
	//minimum number of items
	if i.slice {
		if m, ok := kv.take("min"); ok {
			min, err := strconv.Atoi(m)
			if err != nil {
				return n.errorf("min not an integer")
			}
			i.min = min
		}
	}
	//insert either as flag or as argument
	switch typeName {
	case "flag":
		//cannot have duplicates
		if n.flagNames[name] {
			return n.errorf("flag '%s' already exists", name)
		}
		//flags can also set short names
		if short, ok := kv.take("short"); ok {
			if n.flagNames[short] {
				return n.errorf("flag '%s' (%s) already exists", short, name)
			}
			i.shortName = short
		}
		//add to this command's flags
		n.flags = append(n.flags, i)
	case "arg":
		//add to this command's arguments
		n.args = append(n.args, i)
	default:
		return fmt.Errorf("invalid opts type '%s'", typeName)
	}
	return nil
}

func (n *node) setCmdName(val reflect.Value) error {
	if n.cmdname != nil {
		return n.errorf("cmdname set twice")
	} else if val.Type().Kind() != reflect.String {
		return n.errorf("cmdname type must be string")
	} else if !val.CanAddr() {
		return n.errorf("cannot address cmdname string")
	}
	n.cmdname = val.Addr().Interface().(*string)
	return nil
}

func (n *node) addInlineCmd(name, help string, val reflect.Value) error {
	v := val
	if v.Type().Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Type().Kind() != reflect.Struct {
		return errors.New("inline commands 'type=cmd' must be structs")
	} else if !v.CanAddr() {
		return errors.New("cannot address inline command")
	}
	v = v.Addr()
	//if nil ptr, auto-create new struct
	if v.IsNil() {
		ptr := reflect.New(v.Type().Elem())
		v.Set(ptr)
	}
	//ready!
	if _, ok := n.cmds[name]; ok {
		return n.errorf("command already exists: %s", name)
	}
	sub := newNode(v)
	sub.Name(name)
	sub.help = help
	sub.Description(help)
	sub.parent = n
	n.cmds[name] = sub
	return nil
}

// func (n *node) addArgList(kv map[string]string, val reflect.Value) error {
// 	if len(n.cmds) > 0 {
// 		return n.errorf("argslists and commands cannot be used together")
// 	}
// 	if n.arglist != nil {
// 		return n.errorf("only 1 arglist field is allowed ('%s' already defined)", n.arglist.name)
// 	}
// 	name := kv["name"]
// 	if val.Len() != 0 {
// 		return n.errorf("arglist '%s' is required so it should not be set. "+
// 			"If you'd like to set a default, consider using an option instead.", name)
// 	}
// 	min, _ := strconv.Atoi(kv["min"])
// 	//insert (there can only be one)
// 	n.arglist = &argumentlist{
// 		item: item{
// 			val:  val,
// 			name: name,
// 			help: kv["help"],
// 		},
// 		min: min,
// 	}
// 	return nil
// }

func (n *node) addInternalFlags() error {
	flags := []string{"Help"}
	if n.version != "" {
		flags = append(flags, "Version")
	}
	if n.complete {
		flags = append(flags, "Install", "Uninstall")
	}
	g := reflect.ValueOf(&n.internalOpts).Elem()
	for _, flag := range flags {
		sf, _ := g.Type().FieldByName(flag)
		v := g.FieldByName(flag)
		if err := n.addStructField(sf, v); err != nil {
			return n.errorf("error adding internal %s flag: %s - please report issue", flag, err)
		}
	}
	return nil
}

func (n *node) addFlagsets() error {
	//add provided flag sets
	for _, fs := range n.flagsets {
		var err error
		//add all flags in each set
		fs.VisitAll(func(f *flag.Flag) {
			//convert into item
			val := reflect.ValueOf(f.Value)
			i, er := newItem(val)
			if er != nil {
				err = n.errorf("imported flag '%s': %s", f.Name, er)
				return
			}
			i.name = f.Name
			i.defstr = f.DefValue
			i.help = f.Usage
			//cannot have duplicates
			if n.flagNames[i.name] {
				err = n.errorf("imported flag '%s' already exists", i.name)
				return
			}
			//ready!
			n.flags = append(n.flags, i)
			n.flagNames[i.name] = true
			//convert f into a black hole
			f.Value = noopValue
		})
		//
		if err != nil {
			return err
		}
		fs.Init(fs.Name(), flag.ContinueOnError)
		fs.SetOutput(ioutil.Discard)
		fs.Parse([]string{}) //ensure this flagset returns Parsed() => true
	}
	return nil
}

func (n *node) setPkgDefaults() {
	//attempt to infer package name, repo, author
	configStruct := n.item.val.Elem().Type()
	pkgPath := configStruct.PkgPath()
	parts := strings.Split(pkgPath, "/")
	if len(parts) >= 3 {
		if n.authorInfer && n.author == "" {
			n.author = parts[1]
		}
		if n.repoInfer && n.repo == "" {
			switch parts[0] {
			case "github.com", "bitbucket.org":
				n.repo = "https://" + strings.Join(parts[0:3], "/")
			}
		}
	}
}
