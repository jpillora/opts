package opts

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

//Name sets the name of the program
func (n *node) Name(name string) Opts {
	n.name = name
	return n
}

//Version sets the version of the program
//and renders the 'version' template in the help text
func (n *node) Version(version string) Opts {
	//add version option
	g := reflect.ValueOf(&n.internalOpts).Elem()
	n.addOptArg(g.Type().Field(1), g.Field(1))
	n.version = version
	return n
}

//Repo sets the repository link of the program
//and renders the 'repo' template in the help text
func (n *node) Repo(repo string) Opts {
	n.repo = repo
	return n
}

//PkgRepo infers the repository link of the program
//from the package import path of the struct (So note,
//this will not work for 'main' packages)
func (n *node) PkgRepo() Opts {
	if n.pkgrepo == "" {
		return n.errorf("Package repository could not be infered")
	}
	n.Repo(n.pkgrepo)
	return n
}

//Author sets the author of the program
//and renders the 'author' template in the help text
func (n *node) Author(author string) Opts {
	n.author = author
	return n
}

//PkgRepo infers the repository link of the program
//from the package import path of the struct (So note,
//this will not work for 'main' packages)
func (n *node) PkgAuthor() Opts {
	if n.pkgrepo == "" {
		return n.errorf("Package author could not be infered")
	}
	n.Author(n.pkgauthor)
	return n
}

//Set the padding width
func (n *node) SetPadWidth(p int) Opts {
	n.PadWidth = p
	return n
}

//Disable auto-padding
func (n *node) DisablePadAll() Opts {
	n.PadAll = false
	return n
}

//Set the line width (defaults to 72)
func (n *node) SetLineWidth(l int) Opts {
	n.LineWidth = l
	return n
}

//ConfigPath defines a path to a JSON file which matches
//the structure of the provided config. Environment variables
//override JSON Config variables.
func (n *node) ConfigPath(path string) Opts {
	n.cfgPath = path
	return n
}

//UseEnv enables an implicit "env" struct tag option on
//all struct fields, the name of the field is converted
//into an environment variable with the transform
//`FooBar` -> `FOO_BAR`.
func (n *node) UseEnv() Opts {
	n.useEnv = true
	return n
}

//DocBefore inserts a text block before the specified template
func (n *node) DocBefore(target, newID, template string) Opts {
	return n.docOffset(0, target, newID, template)
}

//DocAfter inserts a text block after the specified template
func (n *node) DocAfter(target, newID, template string) Opts {
	return n.docOffset(1, target, newID, template)
}

//DecSet replaces the specified template
func (n *node) DocSet(id, template string) Opts {
	if _, ok := DefaultTemplates[id]; !ok {
		if _, ok := n.templates[id]; !ok {
			n.errorf("template does not exist: %s", id)
			return n
		}
	}
	n.templates[id] = template
	return n
}

//=================================

func (n *node) addFields(c reflect.Value) *node {
	t := c.Type()
	k := t.Kind()
	//deref pointer
	if k == reflect.Ptr {
		c = c.Elem()
		t = c.Type()
		k = t.Kind()
	}
	if k != reflect.Struct {
		n.errorf("opts: %s should be a pointer to a struct (got %s)", t.Name(), k)
		return n
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
		if sf.Type.Implements(flagValueType) {
			n.addOptArg(sf, val)
			continue
		}
		switch k {
		case reflect.Ptr, reflect.Struct:
			if sf.Tag.Get("type") == "embedded" {
				n.addFields(val)
			} else {
				n.addCmd(sf, val)
			}
		case reflect.Slice:
			if sf.Type.Elem().Kind() != reflect.String {
				n.errorf("slice (list) types must be []string")
				return n
			}
			if sf.Tag.Get("type") == "commalist" {
				n.addOptArg(sf, val)
			} else if sf.Tag.Get("type") == "spacelist" {
				n.addOptArg(sf, val)
			} else {
				n.addArgList(sf, val)
			}
		case reflect.Bool, reflect.String, reflect.Int, reflect.Int64:
			if sf.Tag.Get("type") == "cmdname" {
				if k != reflect.String {
					n.errorf("cmdname field '%s' must be a string", sf.Name)
					return n
				}
				n.cmdname = &val
			} else {
				n.addOptArg(sf, val)
			}
		case reflect.Interface:
			n.errorf("Struct field '%s' interface type must implement flag.Value", sf.Name)
			return n
		default:
			n.errorf("Struct field '%s' has unsupported type: %s", sf.Name, k)
			return n
		}
	}
	return n
}

func (n *node) errorf(format string, args ...interface{}) *node {
	//only store the first error
	if n.erred == nil {
		n.erred = fmt.Errorf(format, args...)
	}
	return n
}

func (n *node) addCmd(sf reflect.StructField, val reflect.Value) {
	if n.arglist != nil {
		n.errorf("argslists and commands cannot be used together")
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
	// log.Printf("define cmd: %s =====", subname)
	sub := fork(n, val)
	sub.name = name
	sub.help = sf.Tag.Get("help")
	n.cmds[name] = sub
}

func (n *node) addArgList(sf reflect.StructField, val reflect.Value) {
	if len(n.cmds) > 0 {
		n.errorf("argslists and commands cannot be used together")
		return
	}
	if n.arglist != nil {
		n.errorf("only 1 arglist field is allowed ('%s' already defined)", n.arglist.name)
		return
	}
	name := sf.Tag.Get("name")
	if name == "" || name == "!" {
		name = camel2dash(sf.Name) //default to struct field name
	}
	if val.Len() != 0 {
		n.errorf("arglist '%s' is required so it should not be set. "+
			"If you'd like to set a default, consider using an option instead.", name)
		return
	}
	min, _ := strconv.Atoi(sf.Tag.Get("min"))
	//insert
	n.arglist = &argumentlist{
		item: item{
			val:  val,
			name: name,
			help: sf.Tag.Get("help"),
		},
		min: min,
	}
}

var durationType = reflect.TypeOf(time.Second)

func (n *node) addOptArg(sf reflect.StructField, val reflect.Value) {
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
			n.errorf("option env name '%s' already in use", i.name)
			return
		}
		n.envnames[i.envName] = true
		i.useEnv = true
	}
	//opt names cannot clash with each other
	if n.optnames[i.name] {
		n.errorf("option name '%s' already in use", i.name)
		return
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
				n.errorf("option short name '%s' already in use", short)
				return
			}
			n.optnames[i.shortName] = true
			i.shortName = short
		}
		// log.Printf("define option: %s %s", name, sf.Type)
		n.flags = append(n.flags, i)
	case "arg":
		//TODO allow other types in 'arg' fields
		if sf.Type.Kind() != reflect.String {
			n.errorf("arg '%s' type must be a string", i.name)
			return
		}
		n.args = append(n.args, i)
	default:
		n.errorf("Invalid optype: %s", t)
	}
}

func (n *node) docOffset(offset int, target, newID, template string) *node {
	if _, ok := n.templates[newID]; ok {
		n.errorf("new template already exists: %s", newID)
		return n
	}
	for i, id := range n.order {
		if id == target {
			n.templates[newID] = template
			index := i + offset
			rest := []string{newID}
			if index < len(n.order) {
				rest = append(rest, n.order[index:]...)
			}
			n.order = append(n.order[:index], rest...)
			return n
		}
	}
	n.errorf("target template not found: %s", target)
	return n
}
