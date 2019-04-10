package opts

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

//AddCommand to this opts instance
func (n *node) AddCommand(cmd Opts) Opts {
	sub, ok := cmd.(*node)
	if !ok {
		panic("another implementation of opts???")
	}
	//for subcommands, default name is package name
	pkgPath := sub.item.val.Elem().Type().PkgPath()
	if sub.name == "" && pkgPath != "" {
		parts := strings.Split(pkgPath, "/")
		sub.name = parts[len(parts)-1]
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
	n.cmds[sub.name] = sub
	sub.parent = n
	return n
}

func (n *node) EmbedStruct(str interface{}) Opts {
	val := reflect.ValueOf(str)
	err := n.addStructFields(val)
	if err != nil {
		panic(err)
	}
	return n
}

func (n *node) matchedCommand() *node {
	if n.cmd != nil {
		return n.cmd.matchedCommand()
	}
	return n
}

type runner1 interface {
	Run() error
}

type runner2 interface {
	Run()
}

//IsRunnable
func (n *node) IsRunnable() bool {
	m := n.matchedCommand()
	v := m.val.Interface()
	_, ok1 := v.(runner1)
	_, ok2 := v.(runner2)
	return ok1 || ok2
}

//Run the parsed configuration
func (n *node) Run() error {
	m := n.matchedCommand()
	v := m.val.Interface()
	r1, ok1 := v.(runner1)
	if ok1 {
		return r1.Run()
	}
	r2, ok2 := v.(runner2)
	if ok2 {
		r2.Run()
		return nil
	}
	return fmt.Errorf("command '%s' is not runnable", m.name)
}

//Run the parsed configuration
func (n *node) RunFatal() {
	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
