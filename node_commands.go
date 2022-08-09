package opts

import (
	"fmt"
	"os"
	"path"
)

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
	if _, exists := n.cmds[sub.name]; exists {
		n.errorf("cannot add command, '%s' already exists", sub.name)
		return n
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

func (n *node) Children() map[string]Opts {
	result := make(map[string]Opts)
	for name, cmd := range n.cmds {
		result[name] = cmd
	}
	return result
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
