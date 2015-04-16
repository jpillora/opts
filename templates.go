package flag

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

var DefaultTemplates = map[string]string{
	"hasoptions": `{{with len .Opts}}{{if ge . 0}} [options]{{end}}{{end}}`,
	"help": `Usage: {{.Name }}{{template "hasoptions" .}}` +
		"\n" +
		`{{template "version" .}}` +
		`{{template "optionslist" .}}` +
		`{{template "author" .}}` +
		`{{template "repo" .}}`,
	"version": "{{if .Version}}\nVersion: {{.Version}}\n{{end}}",
	"optionslist": `{{if .Opts}}` + "\nOptions:\n" +
		`{{range .Opts}}{{template "option" .}}{{end}}{{end}}`,
	"option": `{{.Pad}}{{.Name}}{{if .Help}}{{.Pad}}{{.Help}}{{end}}` + "\n",
	"repo":   "{{if .Repo}}\nRead more:\n{{.Pad}}{{.Repo}}\n{{end}}",
	"author": "{{if .Author}}\nAuthor:\n{{.Pad}}{{.Author}}\n{{end}}",
}

type tFlag struct {
	Opts                        []*tOption
	Name, Version, Repo, Author string
	Pad                         string
}

type tOption struct {
	Name string
	Help string
	Pad  string
}

var anyspace = regexp.MustCompile(`[\s]+`)

func (f *Flag) Help() string {

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
	tf := &tFlag{
		Name:    f.name,
		Version: f.version,
		Repo:    f.repo,
		Author:  f.author,
		Opts:    opts,
	}

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

	padsInOption := f.PadWidth * 2
	spaces := nletters(' ', max+padsInOption) //extra spaces
	helpWidth := f.LineWidth - max

	for i, to := range opts {
		//pad
		to.Name += spaces[:max-len(to.Name)]
		//constrain help text
		words := anyspace.Split(f.opts[i].help, -1)
		n := 0
		for i, w := range words {
			d := helpWidth - n
			n += len(w)
			if n > helpWidth && d < n-helpWidth {
				n = 0
				w = "\n" + string(spaces) + w
			}
			words[i] = w
		}
		to.Help = strings.Join(words, " ")
	}

	//parse each template
	t := template.New(f.name)
	for name, str := range DefaultTemplates {
		//check for user template
		if s, ok := f.Templates[name]; ok {
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
