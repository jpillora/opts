package opts

import (
	"fmt"
	"log"
	"strings"
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
		parts := strings.Split(pkgPath, "/")
		sub.name = parts[len(parts)-1]
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

//IsRunnable
func (n *node) IsRunnable() bool {
	ok, _ := n.run(true)
	return ok
}

//Run the parsed configuration
func (n *node) Run() error {
	_, err := n.run(false)
	return err
}

type runner1 interface {
	Run() error
}

type runner2 interface {
	Run()
}

func (n *node) run(test bool) (bool, error) {
	m := n.matchedCommand()
	v := m.val.Addr().Interface()
	r1, ok1 := v.(runner1)
	r2, ok2 := v.(runner2)
	if test {
		return ok1 || ok2, nil
	}
	if ok1 {
		return true, r1.Run()
	}
	if ok2 {
		r2.Run()
		return true, nil
	}
	return false, fmt.Errorf("command '%s' is not runnable", m.name)
}

//Run the parsed configuration
func (n *node) RunFatal() {
	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
