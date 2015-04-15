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
	subs   map[string]*Flag
	opts   []*option
	//publics
	Name      string
	Version   string
	Author    string
	Args      []string
	LineWidth int
	Templates map[string]string
}

type option struct {
	name string
	val  reflect.Value
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

	//meaningless to modify an incorrect copy of the struct
	if k != reflect.Ptr {
		log.Fatalf("flag.New(config): config should be a pointer (%s) to a struct", k)
	}

	c = c.Elem()
	t = c.Type()
	k = t.Kind()

	if k != reflect.Struct {
		log.Fatalf("flag.New(config): config should be a pointer to a struct (%s)", k)
	}

	//copy defaults
	tmpls := map[string]string{}
	for k, v := range defaultTemplates {
		defaultTemplates[k] = v
	}

	//instantiate
	f := &Flag{
		config: c,
		parent: parent,
		subs:   map[string]*Flag{},
		opts:   []*option{},
		//public defaults
		Name:      "",
		Version:   "",
		Author:    "",
		Args:      os.Args[1:],
		LineWidth: 72,
		Templates: tmpls,
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
	if sf.Type.Kind() == reflect.Struct {
		val = val.Addr()
	}
	subname := camel2dash(sf.Name)
	log.Printf("define subcmd: %s =====", subname)
	sub := fork(f, val)
	sub.Name = subname
	f.subs[subname] = sub
}

func (f *Flag) addOption(sf reflect.StructField, val reflect.Value) {

	n := camel2dash(sf.Name)
	log.Printf("define option: %s %s", n, sf.Type)

	f.opts = append(f.opts, &option{
		name: n,
		val:  val,
	})
}

func (f *Flag) Parse() *Flag {

	flagset := flag.NewFlagSet("tmp", flag.ContinueOnError)
	flagset.Usage = func() {
		fmt.Fprint(os.Stdout, f.Help())
	}

	for _, opt := range f.opts {
		log.Printf("parse prepare option: %s", opt.name)
		//take address, not value
		addr := opt.val.Addr().Interface()
		switch addr := addr.(type) {
		case *string:
			flagset.StringVar(addr, opt.name, "", "")
		case *int:
			flagset.IntVar(addr, opt.name, 0, "")
		}
	}

	log.Printf("parse %+v", f.Args)
	flagset.Parse(f.Args)
	//user config is now set
	return f
}
