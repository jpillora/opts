package opts

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"strings"
)

type subCmdHolderI interface {
	isSubCmdHolder()
}

type subCmdHolder struct{}

func (subCmdHolder) isSubCmdHolder() {}

//NewPlaceholder creates a placeholder subcommand to be used in AddCommand
func NewPlaceholder(name string) Opts {
	sub := newNode(reflect.ValueOf(&subCmdHolder{}))
	sub.name = name
	return sub
}

//NewStruct adds multiple Opts instances under the sub-command name.
//The name can be a slash separated path.
func NewStruct(path string, subs ...Opts) Opts {
	names := strings.Split(path, "/")
	fn := func(name string) *node {
		cmd := newNode(reflect.ValueOf(&subCmdHolder{}))
		cmd.name = names[0]
		return cmd
	}
	return addStruct(fn, path, subs...)
}

func (n *node) AddStruct(path string, subs ...Opts) Opts {
	names := strings.Split(path, "/")
	fn := func(name string) *node {
		return n.newStruct(names[0])
	}
	return addStruct(fn, path, subs...)
}

func addStruct(fn func(string) *node, name string, subs ...Opts) Opts {
	names := strings.Split(name, "/")
	cmd := fn(names[0])
	cmd0 := cmd
	for _, name0 := range names[1:] {
		cmd1 := cmd0.newStruct(name0)
		cmd0.AddCommand(cmd1)
		cmd0 = cmd1
	}
	for _, sub := range subs {
		cmd0.AddCommand(sub)
	}
	return cmd
}

func (n *node) newStruct(name string) *node {
	if sub, exists := n.cmds[name]; exists {
		return sub
	} else {
		sub := newNode(reflect.ValueOf(&subCmdHolder{}))
		sub.name = name
		sub.parent = n
		n.cmds[name] = sub
		return sub
	}
}

func (n *node) AddCommand(cmd Opts) Opts {
	sub, ok := cmd.(*node)
	if !ok {
		panic("another implementation of opts???")
	}
	//default name should be package name,
	//unless its in the main package, then
	//the default becomes the struct name
	structType := sub.item.val.Type()
	pkgPath := structType.PkgPath()
	if sub.name == "" && pkgPath != "main" && pkgPath != "" {
		_, sub.name = path.Split(pkgPath)
	}
	structName := structType.Name()
	if sub.name == "" && structName != "" {
		sub.name = camel2dash(structName)
	}
	//if still no name, needs to be manually set
	if sub.name == "" {
		n.errorf("cannot add command, please set a Name()")
		return n
	}
	if sub0, exists := n.cmds[sub.name]; exists {
		if _, ok1 := sub0.item.val.Interface().(subCmdHolderI); ok1 {
			if _, ok2 := sub.item.val.Interface().(subCmdHolderI); !ok2 {
				// replace sub cmd holding with provided sub cmd
				for k, v := range sub0.cmds {
					sub.cmds[k] = v
					v.parent = sub
				}
			}
		} else {
			n.errorf("cannot add command, '%s' already exists", sub.name)
			return n
		}
	}
	sub.parent = n
	n.cmds[sub.name] = sub
	return n
}

func (n *node) matchedCommand() *node {
	if n.cmd != nil {
		return n.cmd.matchedCommand()
	}
	return n
}

//IsRunnable
func (n *node) IsRunnable() bool {
	_, ok, _ := n.run(true)
	return ok
}

//Run the parsed configuration
func (n *node) Run() error {
	_, _, err := n.run(false)
	return err
}

//Selected returns the subcommand picked when parsing the command line
func (n *node) Selected() ParsedOpts {
	m := n.matchedCommand()
	return m
}

type runner1 interface {
	Run() error
}

type runner2 interface {
	Run()
}

func (n *node) run(test bool) (ParsedOpts, bool, error) {
	m := n.matchedCommand()
	v := m.val.Addr().Interface()
	r1, ok1 := v.(runner1)
	r2, ok2 := v.(runner2)
	if test {
		return m, ok1 || ok2, nil
	}
	if ok1 {
		return m, true, r1.Run()
	}
	if ok2 {
		r2.Run()
		return m, true, nil
	}
	if len(m.cmds) > 0 {
		//if matched command has no run,
		//but has commands, show help instead
		return m, false, fmt.Errorf("sub command '%s' is not runnable", m.name)
	}
	return m, false, fmt.Errorf("command '%s' is not runnable", m.name)
}

//Run the parsed configuration
func (n *node) RunFatal() {
	if err := n.Run(); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}
