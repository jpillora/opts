package opts

import (
	"strings"

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
	for _, flag := range n.flags {
		//pick a predictor
		p := predictorFromItem(flag)
		// v := flag.val.Interface()
		//apply to flags
		complete.Log("flag completion %s %#+v\n", flag.name, p)
		c.Flags["--"+flag.name] = p
		if flag.shortName != "" {
			c.Flags["-"+flag.shortName] = p
		}
	}
	//prepare args
	var p complete.Predictor
	for _, arg := range n.args {
		if p == nil {
			c.Args = predictorFromItem(arg)
		} else {
			c.Args = complete.PredictOr(c.Args, predictorFromItem(arg))
		}
	}
	//prepare sub-commands
	for name, subn := range n.cmds {
		c.Sub[name] = subn.nodeCompletion() //recurse
	}
	return c
}

func predictorFromItem(i *item) (p complete.Predictor) {
	if pv, ok := i.val.Interface().(complete.Predictor); ok {
		p = pv
	} else if i.val.CanAddr() {
		if pv, ok := i.val.Addr().Interface().(complete.Predictor); ok {
			p = pv
		}
	} else if _, ok := i.val.Interface().(bool); ok {
		p = complete.PredictNothing
	}
	if p != nil && i.predict != "" {
		panic("predict tag set on Predictable or bool field " + i.name)
	}
	if p == nil {
		switch {
		case i.predict == "":
			p = complete.PredictAnything
		case i.predict == "any":
			p = complete.PredictAnything
		case i.predict == "none":
			p = complete.PredictNothing
		case i.predict == "dirs":
			p = complete.PredictDirs("*")
		case i.predict == "files":
			p = complete.PredictFiles("*")
		case strings.HasPrefix(i.predict, "dirs:"):
			p = complete.PredictDirs(i.predict[len("dirs:"):])
		case strings.HasPrefix(i.predict, "files:"):
			p = complete.PredictFiles(i.predict[len("files:"):])
		default:
			panic("bad predict tag '" + i.predict + "' on field " + i.name)
		}
	}
	return
}
