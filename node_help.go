package opts

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strings"
	"text/template"
)

//data is only used for templating below
type data struct {
	datum        //data is also a datum
	ArgList      *datum
	FlagGroups   []*datumGroup
	Args         []*datum
	Cmds         []*datum
	Order        []string
	Parents      string
	Version      string
	Summary      string
	Repo, Author string
	ErrMsg       string
}

type datum struct {
	Name, Help, Pad string //Pad is Opt.padWidth many spaces
}

type datumGroup struct {
	Name  string
	Flags []*datum
}

var DefaultOrder = []string{
	"usage",
	"summary",
	"args",
	"arglist",
	"flaggroups",
	"cmds",
	"author",
	"version",
	"repo",
	"errmsg",
}

func defaultOrder() []string {
	order := make([]string, len(DefaultOrder))
	copy(order, DefaultOrder)
	return order
}

var DefaultTemplates = map[string]string{
	//The root template simply loops through,
	//the 'order' and renders each template by name
	"help": `{{ $root := . }}` +
		`{{range $t := .Order}}{{ templ $t $root }}{{end}}`,
	//sections, from top to bottom
	"usage": `Usage: {{.Name }}` +
		`{{template "usageoptions" .}}` +
		`{{template "usageargs" .}}` +
		`{{template "usagearglist" .}}` +
		`{{template "usagecmd" .}}` + "\n",
	"usageoptions": ` [options]`,
	"usageargs":    `{{range .Args}} {{.Name}}{{end}}`,
	"usagearglist": `{{if .ArgList}} {{.ArgList.Name}}{{end}}`,
	"usagecmd":     `{{if .Cmds}} <command>{{end}}`,
	//extra help text gets appended to option.Help
	"extradefault":  `{{if .}}default {{.}}{{end}}`,
	"extraenv":      `{{if .}}env {{.}}{{end}}`,
	"extramultiple": `{{if .}}allows multiple{{end}}`,
	//summary
	"summary": "{{if .Summary}}\n{{ .Summary }}\n{{end}}", //Test
	//args and arg section
	"args":    `{{range .Args}}{{template "arg" .}}{{end}}`,
	"arg":     "{{if .Help}}\n{{.Help}}\n{{end}}",
	"arglist": "{{if .ArgList}}{{ if .ArgList.Help}}\n{{.ArgList.Help}}\n{{end}}{{end}}",
	//flags
	"flaggroups": `{{ range $opt := .FlagGroups}}{{template "flaggroup" $opt}}{{end}}`,
	"flaggroup": `{{if .Flags}}` + "\n" +
		"{{if .Name}}{{.Name}} options{{else}}Options{{end}}:\n" +
		`{{ range $opt := .Flags}}{{template "flag" $opt}}{{end}}{{end}}`,
	"flag": `{{.Name}}{{if .Help}}{{.Pad}}{{.Help}}{{end}}` + "\n",
	//cmds
	"cmds": "{{if .Cmds}}\nCommands:\n" +
		`{{ range $sub := .Cmds}}{{template "cmd" $sub}}{{end}}{{end}}`,
	"cmd": "· {{ .Name }}{{if .Help}} - {{ .Help }}{{end}}\n",
	//extras
	"version": "{{if .Version}}\nVersion:\n{{.Pad}}{{.Version}}\n{{end}}",
	"repo":    "{{if .Repo}}\nRead more:\n{{.Pad}}{{.Repo}}\n{{end}}",
	"author":  "{{if .Author}}\nAuthor:\n{{.Pad}}{{.Author}}\n{{end}}",
	"errmsg":  "{{if .ErrMsg}}\nError:\n{{.Pad}}{{.ErrMsg}}\n{{end}}",
}

var trailingSpaces = regexp.MustCompile(`(?m)\ +$`)

//Help renders the help text as a string
func (o *node) Help() string {
	h, err := renderHelp(o)
	if err != nil {
		log.Fatalf("render help failed: %s", err)
	}
	return h
}

func renderHelp(o *node) (string, error) {
	var err error
	//add default templates
	for name, str := range DefaultTemplates {
		if _, ok := o.templates[name]; !ok {
			o.templates[name] = str
		}
	}
	//prepare templates
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
	//parse all templates and "define" themselves as nested templates
	for name, str := range o.templates {
		t, err = t.Parse(fmt.Sprintf(`{{define "%s"}}%s{{end}}`, name, str))
		if err != nil {
			return "", fmt.Errorf("template '%s': %s", name, err)
		}
	}
	//convert node into template data
	tf, err := convert(o)
	if err != nil {
		return "", fmt.Errorf("node convert: %s", err)
	}
	//execute all templates
	b := &bytes.Buffer{}
	err = t.ExecuteTemplate(b, "help", tf)
	if err != nil {
		return "", fmt.Errorf("template execute: %s", err)
	}
	out := b.String()
	if o.padAll {
		/*
			"foo
			bar"
			becomes
			"
			  foo
			  bar
			"
		*/
		lines := strings.Split(out, "\n")
		for i, l := range lines {
			lines[i] = tf.Pad + l
		}
		out = "\n" + strings.Join(lines, "\n") + "\n"
	}
	out = trailingSpaces.ReplaceAllString(out, "")
	return out, nil
}

func convert(o *node) (*data, error) {
	names := []string{}
	curr := o
	for curr != nil {
		names = append([]string{curr.name}, names...)
		curr = curr.parent
	}
	name := strings.Join(names, " ")
	//get item help, with optional default values and env names and
	//constrain to a specific line width
	keys := []string{"default", "env", "multiple"}
	extras := make([]*template.Template, 3)
	for i, k := range keys {
		t, err := template.New("").Parse(o.templates["extra"+k])
		if err != nil {
			return nil, fmt.Errorf("template extra%s: %s", k, err)
		}
		extras[i] = t
	}
	itemHelp := func(i *item, width int) string {
		vals := []interface{}{i.defstr, i.envName, i.slice}
		outs := []string{}
		for i, v := range vals {
			b := strings.Builder{}
			if err := extras[i].Execute(&b, v); err != nil {
				log.Printf(">>> %s: %s", keys[i], err)
			}
			if b.Len() > 0 {
				outs = append(outs, b.String())
			}
		}
		help := i.help
		extra := strings.Join(outs, ", ")
		if help == "" {
			help = extra
		} else if extra != "" {
			help += " (" + extra + ")"
		}
		return constrain(help, width)
	}
	args := make([]*datum, len(o.args))
	for i, arg := range o.args {
		//mark argument as required
		n := "<" + arg.name + ">"
		if arg.defstr != "" { //or optional
			n = "[" + arg.name + "]"
		}
		args[i] = &datum{
			Name: n,
			Help: itemHelp(arg, o.lineWidth),
		}
	}
	//TODO:
	// var arglist *datum
	// if o.arglist != nil {
	// 	n := o.arglist.name + "..."
	// 	if o.arglist.min == 0 { //optional
	// 		n = "[" + n + "]"
	// 	}
	// 	arglist = &datum{
	// 		Name: n,
	// 		Help: itemHelp(&o.arglist.item, o.lineWidth),
	// 	}
	// }
	flagGroups := make([]*datumGroup, len(o.flagGroups))
	//initialise and calculate padding
	max := 0
	pad := nletters(' ', o.padWidth)
	for i, g := range o.flagGroups {
		dg := &datumGroup{
			Name:  g.name,
			Flags: make([]*datum, len(g.flags)),
		}
		flagGroups[i] = dg
		for i, item := range g.flags {
			to := &datum{Pad: pad}
			to.Name = "--" + item.name
			if item.shortName != "" {
				to.Name += ", -" + item.shortName
			}
			l := len(to.Name)
			//max shared across ALL groups
			if l > max {
				max = l
			}
			dg.Flags[i] = to
		}
	}
	padsInOption := o.padWidth
	optionNameWidth := max + padsInOption
	spaces := nletters(' ', optionNameWidth)
	helpWidth := o.lineWidth - optionNameWidth
	//go back and render each option using calculated values
	for i, dg := range flagGroups {
		for j, to := range dg.Flags {
			//pad all option names to be the same length
			to.Name += spaces[:max-len(to.Name)]
			//constrain help text
			item := o.flagGroups[i].flags[j]
			help := itemHelp(item, helpWidth)
			//add a margin
			lines := strings.Split(help, "\n")
			for i, l := range lines {
				if i > 0 {
					lines[i] = spaces + l
				}
			}
			to.Help = strings.Join(lines, "\n")
		}
	}
	//commands
	subs := make([]*datum, len(o.cmds))
	i := 0
	for _, s := range o.cmds {
		subs[i] = &datum{
			Name: s.name,
			Help: s.help,
			Pad:  pad,
		}
		i++
	}
	//convert error to string
	err := ""
	if o.err != nil {
		err = o.err.Error()
	}
	return &data{
		datum: datum{
			Name: name,
			Help: o.help,
			Pad:  pad,
		},
		Args:       args,
		FlagGroups: flagGroups,
		Cmds:       subs,
		Order:      o.order,
		Version:    o.version,
		Summary:    constrain(o.summary, o.lineWidth),
		Repo:       o.repo,
		Author:     o.author,
		ErrMsg:     err,
	}, nil
}
