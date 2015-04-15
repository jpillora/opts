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
		`{{template "repo" .}}` +
		`{{template "author" .}}`,
	"version": "{{if .Version}}\nVersion: {{.Version}}\n{{end}}",
	"optionslist": `{{if .Opts}}` + "\nOptions:\n" +
		`{{range .Opts}}{{template "option" .}}{{end}}{{end}}`,
	"option": `  {{.Name}}{{if .Help}}   {{.Help}}{{end}}` + "\n",
	"repo":   "{{if .Repo}}\nRead more:\n  {{.Repo}}\n{{end}}",
	"author": "{{if .Author}}\nAuthor:\n  {{.Author}}\n{{end}}",
}

type tFlag struct {
	*Flag
	Opts []*tOption
}

type tOption struct {
	*Option
	Name string
	Help string
}

var anyspace = regexp.MustCompile(`[\s]+`)

func (f *Flag) Help() string {

	var err error

	numOpts := len(f.Opts)
	opts := make([]*tOption, numOpts)
	tf := &tFlag{
		Flag: f,
		Opts: opts,
	}

	//calculate padding etc.
	max := 0
	letters := map[string]bool{}
	f.Padding = nletters(' ', f.PadWidth)

	for i, opt := range f.Opts {
		to := &tOption{Option: opt}
		to.Name = "--" + opt.Name
		n := opt.Name[0:1]
		if _, ok := letters[n]; !ok {
			to.Name += ", -" + n
			letters[n] = true
		}
		l := len(to.Name)
		if l > max {
			max = l
		}
		opts[i] = to
	}

	spaces := nletters(' ', max+5) //extra spaces
	helpWidth := f.LineWidth - max

	for _, to := range opts {
		//pad
		to.Name += spaces[:max-numOpts]
		//constrain help text
		words := anyspace.Split(to.Option.Help, -1)
		n := 0
		for i, w := range words {
			n += len(w)
			if n > helpWidth {
				n = 0
				w = "\n" + string(spaces) + w
			}
			words[i] = w
		}
		to.Help = strings.Join(words, " ")
	}

	//last ditch effort at finding a name
	if f.Name == "" {
		if exe, err := osext.Executable(); err == nil {
			_, f.Name = path.Split(exe)
		} else {
			f.Name = "main"
		}
	}

	//
	t := template.New(f.Name)
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

	b := &bytes.Buffer{}
	err = t.ExecuteTemplate(b, "help", tf)
	if err != nil {
		log.Fatalf("Template execute: %s", err)
	}

	out := b.String()

	if f.PadAll {
		lines := strings.Split(out, "\n")
		for i, l := range lines {
			lines[i] = f.Padding + l
		}
		out = "\n" + strings.Join(lines, "\n") + "\n"
	}

	return out
}
