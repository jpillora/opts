package opts

import (
	"bytes"
	"fmt"
	"log"
	"path"
	"regexp"
	"strings"
	"text/template"

	"github.com/kardianos/osext"
)

var DefaultOrder = []string{
	"usage",
	"args",
	"arglist",
	"options",
}

var DefaultTemplates = map[string]string{
	//loop through the default order and render all templates
	"help": `{{ $root := . }}{{range $t := .Order}}{{ templ $t $root }}{{end}}`,
	//sections, from top to bottom
	"usage":        `Usage: {{.Name }}{{template "usageoptions" .}}{{template "usageargs" .}}{{template "usagearglist" .}}` + "\n",
	"usageoptions": `{{ $nopts := len .Opts}}{{if gt $nopts 0}} [options]{{end}}`,
	"usageargs":    `{{ range .Args}} {{.Name}}{{end}}`,
	"usagearglist": `{{ if .ArgList}} {{.ArgList.Name}}{{end}}`,
	"args":         `{{ range .Args}}{{template "arg" .}}{{end}}`,
	"arg":          "{{ if .Help}}\n{{.Help}}\n{{end}}",
	"arglist":      "{{ if .ArgList}}{{ if .ArgList.Help}}\n{{.ArgList.Help}}\n{{end}}{{end}}",
	"options": `{{if .Opts}}` + "\nOptions:\n" +
		`{{ range $opt := .Opts}}{{template "option" $opt}}{{end}}{{end}}`,
	"option":  `{{.Name}}{{if .Help}}{{.Pad}}{{.Help}}{{end}}` + "\n",
	"version": "\nVersion:\n{{.Pad}}{{.Version}}\n",
	"repo":    "\nRead more:\n{{.Pad}}{{.Repo}}\n",
	"author":  "\nAuthor:\n{{.Pad}}{{.Author}}\n",
}

type tOpts struct {
	Args                        []*targument
	ArgList                     *targlist
	Opts                        []*toption
	Order                       []string
	Name, Version, Repo, Author string
	Pad                         string
}

type toption struct {
	Name string
	Help string
	Pad  string
}

type targument struct {
	Name string
	Help string
}

type targlist struct {
	Name string
	Help string
}

var anyspace = regexp.MustCompile(`[\s]+`)

func (o *Opts) Help() string {

	//last ditch effort at finding a name
	if o.name == "" {
		if exe, err := osext.Executable(); err == nil {
			_, o.name = path.Split(exe)
		} else {
			o.name = "main"
		}
	}

	var err error
	args := make([]*targument, len(o.args))
	for i, arg := range o.args {
		//mark argument as required
		n := "<" + arg.name + ">"
		if arg.hasdef { //or optional
			n = "[" + arg.name + "]"
		}
		args[i] = &targument{
			Name: n,
			Help: constrain(arg.help, o.LineWidth),
		}
	}

	var arglist *targlist = nil
	if o.arglist != nil {
		n := o.arglist.name + "..."
		if o.arglist.min == 0 { //optional
			n = "[" + n + "]"
		}
		arglist = &targlist{
			Name: n,
			Help: o.arglist.help,
		}
	}

	opts := make([]*toption, len(o.opts))
	tf := &tOpts{
		Args:    args,
		ArgList: arglist,
		Opts:    opts,
		Order:   o.order,
		Name:    o.name,
		Version: o.version,
		Repo:    o.repo,
		Author:  o.author,
	}

	//calculate padding etc.
	max := 0
	shorts := map[string]bool{}
	tf.Pad = nletters(' ', o.PadWidth)

	for i, opt := range o.opts {
		to := &toption{Pad: tf.Pad}
		to.Name = "--" + opt.name
		n := opt.name[0:1]
		if _, ok := shorts[n]; !ok {
			to.Name += ", -" + n
			shorts[n] = true
		}
		l := len(to.Name)
		if l > max {
			max = l
		}
		opts[i] = to
	}

	padsInOption := o.PadWidth
	optionNameWidth := max + padsInOption
	spaces := nletters(' ', optionNameWidth)
	helpWidth := o.LineWidth - optionNameWidth

	//render each option
	for i, to := range opts {
		//pad all names to be the same length
		to.Name += spaces[:max-len(to.Name)]
		//constrain help text
		h := constrain(o.opts[i].help, helpWidth)
		lines := strings.Split(h, "\n")
		for i, l := range lines {
			if i > 0 {
				lines[i] = spaces + l
			}
		}
		to.Help = strings.Join(lines, "\n")
	}

	//begin
	t := template.New(o.name)

	t = t.Funcs(map[string]interface{}{
		//reimplementation of "template" except with dynamic name
		"templ": func(name string, data interface{}) (string, error) {
			b := &bytes.Buffer{}
			err = t.ExecuteTemplate(b, name, data)
			if err != nil {
				return "", err
			}
			return b.String(), nil
		},
	})

	//parse each template
	for name, str := range DefaultTemplates {
		//check for user template
		if s, ok := o.templates[name]; ok {
			str = s
		}
		t, err = t.Parse(fmt.Sprintf(`{{define "%s"}}%s{{end}}`, name, str))
		if err != nil {
			log.Fatalf("Template error: %s: %s", name, err)
		}
	}

	//execute all templates
	b := &bytes.Buffer{}
	err = t.ExecuteTemplate(b, "help", tf)
	if err != nil {
		log.Fatalf("Template execute: %s", err)
	}

	out := b.String()

	if o.PadAll {
		lines := strings.Split(out, "\n")
		for i, l := range lines {
			lines[i] = tf.Pad + l
		}
		out = "\n" + strings.Join(lines, "\n") + "\n"
	}

	return out
}
