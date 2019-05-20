package opts

import (
	"flag"
	"fmt"
)

//errorf to be stored until parse-time
func (n *node) errorf(format string, args ...interface{}) error {
	err := &authorError{fmt.Sprintf(format, args...)}
	//only store the first error
	if n.err == nil {
		n.err = err
	}
	return err
}

//parseErrorf
func (n *node) parseErrorf(format string, args ...interface{}) error {
	err := &parseError{
		msg: fmt.Sprintf(format, args...),
		n:   n,
	}
	//only store the first error
	if n.err == nil {
		n.err = err
	}
	return err
}

//Name sets the name of the program
func (n *node) Name(name string) Opts {
	n.name = name
	return n
}

//Version sets the version of the program
//and renders the 'version' template in the help text
func (n *node) Version(version string) Opts {
	n.version = version
	return n
}

//Summary sets the text summary of the program,
//which, by default, is inserted below the usage text
func (n *node) Summary(summary string) Opts {
	n.summary = summary
	return n
}

//Repo sets the repository link of the program
//and renders the 'repo' template in the help text
func (n *node) Repo(repo string) Opts {
	n.repo = repo
	return n
}

//PkgRepo infers the repository link of the program
//from the package import path of the struct (Note:
//this will not work for 'main' packages)
func (n *node) PkgRepo() Opts {
	n.repoInfer = true
	return n
}

//Author sets the author of the program
//and renders the 'author' template in the help text
func (n *node) Author(author string) Opts {
	n.author = author
	return n
}

//PkgRepo infers the repository link of the program
//from the package import path of the struct (So note,
//this will not work for 'main' packages)
func (n *node) PkgAuthor() Opts {
	n.authorInfer = true
	return n
}

//Set the padding width, which defines the amount padding
//when rendering help text (defaults to 72)
func (n *node) SetPadWidth(p int) Opts {
	n.padWidth = p
	return n
}

func (n *node) DisablePadAll() Opts {
	n.padAll = false
	return n
}

//Set the line width (defaults to 72),
//which defines where new-lines
//are inserted into the help text
//(defaults to 42)
func (n *node) SetLineWidth(l int) Opts {
	n.lineWidth = l
	return n
}

//ConfigPath defines a path to a JSON file which matches
//the structure of the provided config. Environment variables
//override JSON Config variables.
func (n *node) ConfigPath(path string) Opts {
	n.internalOpts.ConfigPath = path
	return n
}

func (n *node) UserConfigPath() Opts {
	n.userCfgPath = true
	return n
}

//UseEnv enables an implicit "env" struct tag option on
//all struct fields, the name of the field is converted
//into an environment variable with the transform
//`FooBar` -> `FOO_BAR`.
func (n *node) UseEnv() Opts {
	n.useEnv = true
	return n
}

//DocBefore inserts a text block before the specified template
func (n *node) DocBefore(target, newID, template string) Opts {
	return n.docOffset(0, target, newID, template)
}

//DocAfter inserts a text block after the specified template
func (n *node) DocAfter(target, newID, template string) Opts {
	return n.docOffset(1, target, newID, template)
}

//DecSet replaces the specified template
func (n *node) DocSet(id, template string) Opts {
	if _, ok := DefaultTemplates[id]; !ok {
		if _, ok := n.templates[id]; !ok {
			n.errorf("template does not exist: %s", id)
			return n
		}
	}
	n.templates[id] = template
	return n
}

func (n *node) docOffset(offset int, target, newID, template string) *node {
	if _, ok := n.templates[newID]; ok {
		n.errorf("new template already exists: %s", newID)
		return n
	}
	for i, id := range n.order {
		if id == target {
			n.templates[newID] = template
			index := i + offset
			rest := []string{newID}
			if index < len(n.order) {
				rest = append(rest, n.order[index:]...)
			}
			n.order = append(n.order[:index], rest...)
			return n
		}
	}
	n.errorf("target template not found: %s", target)
	return n
}

func (n *node) EmbedFlagSet(fs *flag.FlagSet) Opts {
	n.flagsets = append(n.flagsets, fs)
	return n
}

func (n *node) EmbedGlobalFlagSet() Opts {
	return n.EmbedFlagSet(flag.CommandLine)
}

func (n *node) Call(fn func(o Opts)) Opts {
	fn(n)
	return n
}

func (n *node) flagGroup(name string) *itemGroup {
	//NOTE: the default group is the empty string
	//get existing group
	for _, g := range n.flagGroups {
		if g.name == name {
			return g
		}
	}
	//otherwise, create and append
	g := &itemGroup{name: name}
	n.flagGroups = append(n.flagGroups, g)
	return g
}

func (n *node) flags() []*item {
	flags := []*item{}
	for _, g := range n.flagGroups {
		flags = append(flags, g.flags...)
	}
	return flags
}

type authorError struct {
	err string
}

func (o *authorError) Error() string {
	return o.err
}

type parseError struct {
	msg string
	n   *node
}

func (o *parseError) Error() string {
	return o.msg
}

type exitError struct {
	msg string
}

func (o *exitError) Error() string {
	return o.msg
}
