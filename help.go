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

//data is only used for templating below
type data struct {
	ArgList       *datum
	Opts          []*datum
	Args          []*datum
	Subcmds       []*datum
	Order         []string
	Name, Version string
	Repo, Author  string
	Pad           string //Pad is Opt.PadWidth many spaces
	ErrMsg        string
}

type datum struct {
	Name string
	Help string
	Pad  string
}

var DefaultOrder = []string{
	"usage",
	"args",
	"arglist",
	"options",
	"subcmds",
	"author",
	"version",
	"repo",
	"errmsg",
}

var DefaultTemplates = map[string]string{
	//loop through the 'order' and render all templates
	"help": `{{ $root := . }}` +
		`{{range $t := .Order}}{{ templ $t $root }}{{end}}`,
	//sections, from top to bottom
	"usage": `Usage: {{.Name }}` +
		`{{template "usageoptions" .}}` +
		`{{template "usageargs" .}}` +
		`{{template "usagearglist" .}}` +
		`{{template "usagesubcmd" .}}` + "\n",
	"usageoptions": `{{if .Opts}} [options]{{end}}`,
	"usageargs":    `{{range .Args}} {{.Name}}{{end}}`,
	"usagearglist": `{{if .ArgList}} {{.ArgList.Name}}{{end}}`,
	"usagesubcmd":  `{{if .Subcmds}} <subcommand>{{end}}`,
	//args and arg section
	"args":    `{{range .Args}}{{template "arg" .}}{{end}}`,
	"arg":     "{{if .Help}}\n{{.Help}}\n{{end}}",
	"arglist": "{{if .ArgList}}{{ if .ArgList.Help}}\n{{.ArgList.Help}}\n{{end}}{{end}}",
	//options
	"options": `{{if .Opts}}` + "\nOptions:\n" +
		`{{ range $opt := .Opts}}{{template "option" $opt}}{{end}}{{end}}`,
	"option": `{{.Name}}{{if .Help}}{{.Pad}}{{.Help}}{{end}}` + "\n",
	//subcmds
	"subcmds": "{{if .Subcmds}}\nSubcommands:\n" +
		`{{ range $sub := .Subcmds}}{{template "subcmd" $sub}}{{end}}{{end}}`,
	"subcmd": "* {{ .Name }}{{if .Help}} - {{ .Help }}{{end}}\n",
	//extras
	"version": "{{if .Version}}\nVersion:\n{{.Pad}}{{.Version}}\n{{end}}",
	"repo":    "{{if .Repo}}\nRead more:\n{{.Pad}}{{.Repo}}\n{{end}}",
	"author":  "{{if .Author}}\nAuthor:\n{{.Pad}}{{.Author}}\n{{end}}",
	"errmsg":  "{{if .ErrMsg}}\nError:\n{{.Pad}}{{.ErrMsg}}\n{{end}}",
}

var anyspace = regexp.MustCompile(`[\s]+`)

func convert(o *Opts) *data {

	names := []string{}
	curr := o
	for curr != nil {
		names = append([]string{curr.name}, names...)
		curr = curr.parent
	}
	name := strings.Join(names, " ")

	args := make([]*datum, len(o.args))
	for i, arg := range o.args {
		//mark argument as required
		n := "<" + arg.name + ">"
		if arg.hasDef { //or optional
			n = "[" + arg.name + "]"
		}
		args[i] = &datum{
			Name: n,
			Help: constrain(arg.help, o.LineWidth),
		}
	}

	var arglist *datum = nil
	if o.arglist != nil {
		n := o.arglist.name + "..."
		if o.arglist.min == 0 { //optional
			n = "[" + n + "]"
		}
		arglist = &datum{
			Name: n,
			Help: o.arglist.help,
		}
	}

	opts := make([]*datum, len(o.opts))

	//calculate padding etc.
	max := 0
	shorts := map[string]bool{}
	pad := nletters(' ', o.PadWidth)

	for i, opt := range o.opts {
		to := &datum{Pad: pad}
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
		//pad all option names to be the same length
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

	//subcommands
	subs := make([]*datum, len(o.subcmds))
	i := 0
	for _, s := range o.subcmds {
		subs[i] = &datum{
			Name: s.name,
			Help: s.help,
			Pad:  pad,
		}
		i++
	}

	err := ""
	if o.erred != nil {
		err = o.erred.Error()
	}

	return &data{
		Args:    args,
		ArgList: arglist,
		Opts:    opts,
		Subcmds: subs,
		Order:   o.order,
		Name:    name,
		Version: o.version,
		Repo:    o.repo,
		Author:  o.author,
		Pad:     pad,
		ErrMsg:  err,
	}
}

func (o *Opts) Help() string {
	var err error

	//last ditch effort at finding the program name
	root := o
	for root.parent != nil {
		root = root.parent
	}
	if root.name == "" {
		if exe, err := osext.Executable(); err == nil {
			_, root.name = path.Split(exe)
		} else {
			root.name = "main"
		}
	}

	//convert Opts into template data
	tf := convert(o)

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
