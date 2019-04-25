package opts

import (
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
	//make a completion command for this node
	c := complete.Command{
		Sub:         complete.Commands{},
		Flags:       complete.Flags{},
		GlobalFlags: nil,
		Args:        nil,
	}
	//prepare flags
	for _, item := range n.flags {
		//item's predictor
		p := item.predictor
		//default predictor
		if p == nil {
			if item.noarg {
				p = complete.PredictNothing
			} else {
				p = complete.Predictor(complete.PredictAnything)
			}
		}
		//add to completion flags set
		c.Flags["--"+item.name] = p
		if item.shortName != "" {
			c.Flags["-"+item.shortName] = p
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
