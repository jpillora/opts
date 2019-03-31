package opts

import "log"

//AddCommand to this opts instance
func (n *node) AddCommand(cmd Opts) Opts {
	sub, ok := cmd.(*node)
	if !ok {
		panic("another implementation of opts???")
	}
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

func (n *node) matchedCommand() *node {
	if n.cmd != nil {
		return n.cmd.matchedCommand()
	}
	return n
}

type runner interface {
	Run() error
}

//IsRunnable
func (n *node) IsRunnable() bool {
	m := n.matchedCommand()
	v := m.val.Interface()
	_, ok := v.(runner)
	return ok
}

//Run the parsed configuration
func (n *node) Run() error {
	m := n.matchedCommand()
	v := m.val.Interface()
	r, ok := v.(runner)
	if !ok {
		log.Fatalf("command '%s' is not runnable", m.name)
	}
	return r.Run()
}

//Run the parsed configuration
func (n *node) RunFatal() {
	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
