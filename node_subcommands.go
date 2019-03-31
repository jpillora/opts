package opts

func (n *node) AddCommand(cmd Opts) Opts {
	ncmd, ok := cmd.(*node)
	if !ok {
		panic("another implementation of opts???")
	}
	if ncmd.name == "" {
		n.errorf("cannot add command, please set a Name()")
		return n
	}
	if _, exists := n.cmds[ncmd.name]; exists {
		n.errorf("cannot add command, '%s' already exists", ncmd.name)
		return n
	}
	n.cmds[ncmd.name] = ncmd
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
