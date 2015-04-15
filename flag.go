package flag

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
)

//Flag is ...
type Flag struct {
	//privates
	config reflect.Value
	parent *Flag
	Subs   map[string]*Flag
	Opts   []*Option
	//public settings
	Name, Version string
	Repo, Author  string
	Args          []string //os.Args
	LineWidth     int      //42
	PadAll        bool     //true
	PadWidth      int      //2
	Padding       string   //calculated
	Templates     map[string]string
}

type Option struct {
	val reflect.Value
	//"publics" for templating
	Name        string
	DisplayName string //calculated
	TypeName    string
	Help        string
}

//NewFlag creates a new Flag
func New(config interface{}) *Flag {

	v := reflect.ValueOf(config)

	//nil parent -> root command
	f := fork(nil, v)

	//attempt to infer package name and author
	parts := strings.Split(v.Type().PkgPath(), "/")
	if len(parts) >= 2 {
		f.Name = parts[len(parts)-1]
		f.Author = parts[len(parts)-2]
	}

	return f
}

func fork(parent *Flag, c reflect.Value) *Flag {

	t := c.Type()
	k := t.Kind()

	//must be pointer (meaningless to modify a copy of the struct)
	if k != reflect.Ptr {
		log.Fatalf("flag.New(config): config should be a pointer (%s) to a struct", k)
	}

	c = c.Elem()
	t = c.Type()
	k = t.Kind()

	if k != reflect.Struct {
		log.Fatalf("flag.New(config): config should be a pointer to a struct (%s)", k)
	}

	//instantiate
	f := &Flag{
		config: c,
		parent: parent,
		Subs:   map[string]*Flag{},
		Opts:   []*Option{},
		//public defaults
		Name:      "",
		Version:   "",
		Author:    "",
		Args:      os.Args[1:],
		LineWidth: 42,
		PadAll:    true,
		PadWidth:  2,
		Templates: map[string]string{},
	}

	//parse struct fields
	for i := 0; i < c.NumField(); i++ {
		val := c.Field(i)
		sf := t.Field(i)
		switch sf.Type.Kind() {
		case reflect.Ptr, reflect.Struct:
			f.addSubcmd(sf, val)
		case reflect.Bool, reflect.String, reflect.Int:
			f.addOption(sf, val)
		default:
			log.Fatalf("Field type not allowed: %s", sf.Type.Kind().String())
		}
	}

	// fmt.Println(field.Tag.Get("color"), field.Tag.Get("species"))
	return f
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
	sub.Name = subname
	f.Subs[subname] = sub
}

func (f *Flag) addOption(sf reflect.StructField, val reflect.Value) {

	n := camel2dash(sf.Name)
	// log.Printf("define Option: %s %s", n, sf.Type)
	// fmt.Printf("\thelp:%s\n", sf.Tag.Get("help"))
	// fmt.Printf("\tflag:%s\n", sf.Tag.Get("flag"))

	f.Opts = append(f.Opts, &Option{
		val:      val,
		Name:     n,
		TypeName: sf.Type.Name(),
		Help:     sf.Tag.Get("help"),
	})
}

func (f *Flag) Parse() *Flag {

	//peek at args, maybe use subcommand
	if len(f.Args) > 0 {
		a := f.Args[0]
		//matching subcommand, use it
		if sub, exists := f.Subs[a]; exists {
			sub.Args = f.Args[1:]
			return sub.Parse()
		}
	}

	//use this command
	flagset := flag.NewFlagSet(f.Name, flag.ContinueOnError)
	flagset.Usage = func() {
		fmt.Fprint(os.Stdout, f.Help())
	}

	for _, opt := range f.Opts {
		// log.Printf("parse prepare Option: %s", opt.name)
		//take address, not value
		addr := opt.val.Addr().Interface()
		switch addr := addr.(type) {
		case *string:
			flagset.StringVar(addr, opt.Name, "", "")
		case *int:
			flagset.IntVar(addr, opt.Name, 0, "")
		}
	}

	// log.Printf("parse %+v", f.Args)
	flagset.Parse(f.Args)
	//user config is now set
	return f
}
