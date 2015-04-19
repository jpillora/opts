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
	"version",
	"options",
}

var DefaultTemplates = map[string]string{
	"help":       `{{ $r := . }}{{range $o := .Order}}{{ templ $o $r }}{{end}}`,
	"usage":      `Usage: {{.Name }}{{template "hasoptions" .}}` + "\n",
	"hasoptions": `{{with len .Opts}}{{if ge . 0}} [options]{{end}}{{end}}`,
	"version":    "{{if .Version}}\nVersion: {{.Version}}\n{{end}}",
	"options": `{{if .Opts}}` + "\nOptions:\n" +
		`{{range .Opts}}{{template "option" .}}{{end}}{{end}}`,
	"option": `{{.Name}}{{if .Help}}{{.Pad}}{{.Help}}.{{end}}` + "\n",
	"repo":   "\nRead more:\n{{.Pad}}{{.Repo}}\n",
	"author": "\nAuthor:\n{{.Pad}}{{.Author}}\n",
}

type tOpts struct {
	Opts                        []*tOption
	Order                       []string
	Name, Version, Repo, Author string
	Pad                         string
}

type tOption struct {
	Name string
	Help string
	Pad  string
}

var anyspace = regexp.MustCompile(`[\s]+`)

func (f *Opts) Help() string {

	//last ditch effort at finding a name
	if f.name == "" {
		if exe, err := osext.Executable(); err == nil {
			_, f.name = path.Split(exe)
		} else {
			f.name = "main"
		}
	}

	var err error
	opts := make([]*tOption, len(f.opts))
	tf := &tOpts{
		Order:   f.order,
		Name:    f.name,
		Version: f.version,
		Repo:    f.repo,
		Author:  f.author,
		Opts:    opts,
	}

	// log.Printf("order %+v", tf.Order)

	//calculate padding etc.
	max := 0
	shorts := map[string]bool{}
	tf.Pad = nletters(' ', f.PadWidth)

	for i, opt := range f.opts {
		to := &tOption{Pad: tf.Pad}
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

	padsInOption := f.PadWidth
	optionNameWidth := max + padsInOption
	spaces := nletters(' ', optionNameWidth)
	helpWidth := f.LineWidth - optionNameWidth

	//render each option
	for i, to := range opts {
		//pad all names to be the same length
		to.Name += spaces[:max-len(to.Name)]
		//constrain help text
		words := anyspace.Split(f.opts[i].help, -1)
		n := 0
		for i, w := range words {
			d := helpWidth - n
			wn := len(w) + 1 //+space
			n += wn
			if n > helpWidth && n-helpWidth > d {
				n = wn
				w = "\n" + string(spaces) + w
			}
			words[i] = w
		}
		to.Help = strings.Join(words, " ")
	}

	//root
	t := template.New(f.name)

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
		if s, ok := f.templates[name]; ok {
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

	if f.PadAll {
		lines := strings.Split(out, "\n")
		for i, l := range lines {
			lines[i] = tf.Pad + l
		}
		out = "\n" + strings.Join(lines, "\n") + "\n"
	}

	return out
}
