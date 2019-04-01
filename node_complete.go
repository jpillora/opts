package opts

import (
	"reflect"

	"github.com/posener/complete"
	"github.com/posener/complete/cmd/install"
)

//Complete enables shell-completion for this command and
//its subcommands
func (n *node) Complete() Opts {
	n.complete = true
	return n
}

func (n *node) manageCompletion(uninstall bool) error {
	e := &exitError{}
	fn := install.Install
	if uninstall {
		fn = install.Uninstall
	}
	if err := fn(n.name); err != nil {
		e.msg = err.Error()
	} else if uninstall {
		e.msg = "Uninstalled"
	} else {
		e.msg = "Installed"
	}
	return e //always exit
}

func (n *node) doCompletion() bool {
	return complete.New(n.name, n.nodeCompletion()).Complete()
}

func (n *node) nodeCompletion() complete.Command {
	//make a command for this node
	c := complete.Command{
		Sub:         complete.Commands{},
		Flags:       complete.Flags{},
		GlobalFlags: nil,
		Args:        nil,
	}
	//prepare flags
	for _, flag := range n.flags {
		p := complete.Predictor(complete.PredictAnything)
		if flag.val.Kind() == reflect.Bool {
			p = complete.PredictNothing
		}
		c.Flags["--"+flag.name] = p
		c.Flags["-"+flag.name] = p
		if flag.shortName != "" {
			c.Flags["-"+flag.shortName] = p
		}
	}
	//prepare args
	if len(n.args) > 0 {
		c.Args = complete.PredictAnything
	}
	//prepare sub-commands
	for name, subn := range n.cmds {
		c.Sub[name] = subn.nodeCompletion() //recurse
	}
	return c
}
