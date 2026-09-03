package opts

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/template"
)

// data is only used for templating below
type data struct {
	datum        //data is also a datum
	FlagGroups   []*datumGroup
	Args         []*datum
	CmdGroups    []*datumGroup
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

// DefaultOrder defines which templates get rendered in which order.
// This list is referenced in the "help" template below.
var DefaultOrder = []string{
	"usage",
	"summary",
	"args",
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

// DefaultTemplates define a set of individual templates
// that get rendered in DefaultOrder. You can replace templates or insert templates before or after existing
// templates using the DocSet, DocBefore and DocAfter methods. For example, you can insert a string after the
// usage text with:
//
//	DocAfter("usage", "this is a string, and if it is very long, it will be wrapped")
//
// The entire help text is simply the "help" template listed below, which renders a set of these templates in
// the order defined above. All templates can be referenced using the keys in this map:
var DefaultTemplates = map[string]string{
	"help":          `{{ $root := . }}{{range $t := .Order}}{{ templ $t $root }}{{end}}`,
	"usage":         `{{bold "Usage:"}} {{accentBold .Name }} {{accent "[options]"}}{{template "usageargs" .}}{{template "usagecmd" .}}` + "\n",
	"usageargs":     `{{range .Args}} {{accent .Name}}{{end}}`,
	"usagecmd":      `{{if .CmdGroups}} {{accent "<command>"}}{{end}}`,
	"extradefault":  `{{if .}}default {{.}}{{end}}`,
	"extraenv":      `{{if .}}env {{.}}{{end}}`,
	"extramultiple": `{{if .}}allows multiple{{end}}`,
	"summary":       "{{if .Summary}}\n{{ .Summary }}\n{{end}}",
	"args":          `{{range .Args}}{{template "arg" .}}{{end}}`,
	"arg":           "{{if .Help}}\n{{.Help}}\n{{end}}",
	"flaggroups":    `{{ range $g := .FlagGroups}}{{template "flaggroup" $g}}{{end}}`,
	"flaggroup": "{{if .Flags}}\n{{if .Name}}{{bold (printf \"%s options:\" .Name)}}{{else}}{{bold \"Options:\"}}{{end}}\n" +
		`{{ range $f := .Flags}}{{template "flag" $f}}{{end}}{{end}}`,
	"flag": `{{accentPadded .Name}}{{if .Help}}{{.Pad}}{{.Help}}{{end}}` + "\n",
	"cmds": `{{ range $g := .CmdGroups}}{{template "cmdgroup" $g}}{{end}}`,
	"cmdgroup": "{{if .Flags}}\n{{if .Name}}{{bold (printf \"%s commands:\" .Name)}}{{else}}{{bold \"Commands:\"}}{{end}}\n" +
		`{{ range $sub := .Flags}}{{template "cmd" $sub}}{{end}}{{end}}`,
	"cmd":     "· {{accent .Name}}{{if .Help}}{{.Pad}}{{ .Help }}{{end}}\n",
	"version": "{{if .Version}}\n{{bold \"Version:\"}}\n{{.Pad}}{{accent .Version}}\n{{end}}",
	"repo":    "{{if .Repo}}\n{{bold \"Read more:\"}}\n{{.Pad}}{{accent .Repo}}\n{{end}}",
	"author":  "{{if .Author}}\n{{bold \"Author:\"}}\n{{.Pad}}{{accent .Author}}\n{{end}}",
	"errmsg":  "{{if .ErrMsg}}\n{{dangerBold \"Error:\"}}\n{{.Pad}}{{danger .ErrMsg}}\n{{end}}",
}

var (
	//Keep ANSI resets at the end of a line while removing spaces before them.
	trailingSpaces   = regexp.MustCompile(`(?m) +((?:\x1b\[[0-9;]*m)*)$`)
	trailingBrackets = regexp.MustCompile(`^(.+)\(([^\)]+)\)$`)
)

const (
	ansiReset    = "\x1b[0m"
	ansiBold     = "\x1b[1m"
	ansiCyan     = "\x1b[36m"
	ansiBoldCyan = "\x1b[1;36m"
	ansiRed      = "\x1b[31m"
	ansiBoldRed  = "\x1b[1;31m"
)

type helpStyler struct {
	enabled bool
}

func (s helpStyler) wrap(code, text string) string {
	if !s.enabled || text == "" {
		return text
	}
	return code + text + ansiReset
}

func (s helpStyler) bold(text string) string {
	return s.wrap(ansiBold, text)
}

func (s helpStyler) accent(text string) string {
	return s.wrap(ansiCyan, text)
}

func (s helpStyler) accentBold(text string) string {
	return s.wrap(ansiBoldCyan, text)
}

// accentPadded leaves alignment spaces outside the ANSI sequence so color is
// applied only to the option name, not the gap between the two help columns.
func (s helpStyler) accentPadded(text string) string {
	name := strings.TrimRight(text, " ")
	return s.accent(name) + text[len(name):]
}

func (s helpStyler) danger(text string) string {
	return s.wrap(ansiRed, text)
}

func (s helpStyler) dangerBold(text string) string {
	return s.wrap(ansiBoldRed, text)
}

// highlightHelp reports whether ANSI styling should be used for a terminal.
// NO_COLOR and TERM=dumb follow the conventions used by other command-line
// tools to explicitly request plain output.
func highlightHelp(supportsANSI bool) bool {
	if !supportsANSI || os.Getenv("NO_COLOR") != "" {
		return false
	}
	return !strings.EqualFold(os.Getenv("TERM"), "dumb")
}

const (
	//defaultLineWidth is used when the terminal dimensions are unknown
	defaultLineWidth = 96
	//minHelpWidth is the narrowest second column still worth rendering.
	//Below this, help text is stacked underneath its name instead.
	minHelpWidth = 24
	//cmdPrefixWidth is the display width of the "· " command bullet
	cmdPrefixWidth = 2
)

// renderWidth returns the number of characters available between the left and
// right padAll margins. The configured, detected, or fallback width always
// describes the complete rendered line, so both margins are removed after
// choosing that width.
func (o *node) renderWidth(detected int) int {
	w := 0
	explicit := false
	for n := o; n != nil; n = n.parent {
		if n.lineWidth > 0 {
			w = n.lineWidth
			explicit = true
			break
		}
	}
	if !explicit {
		w = detected
		if w <= 0 {
			w = defaultLineWidth
		}
		//wide terminals are capped, very long lines are hard to read
		if w > defaultLineWidth {
			w = defaultLineWidth
		}
	}
	//padAll leaves the configured margin at both edges
	if o.padAll {
		w -= 2 * o.padWidth
	}
	//an absurdly narrow terminal is still honoured, rendering
	//wider than the terminal only makes it harder to read
	if w < 1 {
		w = 1
	}
	return w
}

// column computes the two column layout used by the option and command lists.
// The first column is as wide as the longest name (plus padding) and the second
// takes the remaining width. When the second column would be too narrow to be
// useful, stacked is returned true, and the caller should instead render the
// help text on its own line, indented by the (much smaller) first column.
func column(nameWidth, pad, total int) (indent int, help int, stacked bool) {
	indent = nameWidth + pad
	help = total - indent
	if help < minHelpWidth {
		indent = pad
		help = total - indent
		stacked = true
	}
	if help < 1 {
		help = 1
	}
	return
}

// wrap constrains help text to the given width, indenting every line after
// the first, such that the text forms a column beginning at indent.
func wrap(help string, width int, indent string) string {
	help = constrain(help, width)
	lines := strings.Split(help, "\n")
	for i, l := range lines {
		if i > 0 {
			lines[i] = indent + l
		}
	}
	return strings.Join(lines, "\n")
}

// Help renders the help text as a string
func (o *node) Help() string {
	h, err := renderHelp(o)
	if err != nil {
		log.Fatalf("render help failed: %s", err)
	}
	return h
}

func renderHelp(o *node) (string, error) {
	var err error
	terminal := termInfo()
	styler := helpStyler{enabled: highlightHelp(terminal.isTTY && terminal.supportsANSI)}
	//add default templates
	for name, str := range DefaultTemplates {
		if _, ok := o.templates[name]; !ok {
			o.templates[name] = str
		}
	}
	//prepare templates
	t := template.New(o.name)
	t = t.Funcs(map[string]interface{}{
		"bold":         styler.bold,
		"accent":       styler.accent,
		"accentBold":   styler.accentBold,
		"accentPadded": styler.accentPadded,
		"danger":       styler.danger,
		"dangerBold":   styler.dangerBold,
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
	tf, err := convert(o, terminal.width)
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
	out = trailingSpaces.ReplaceAllString(out, "$1")
	return out, nil
}

func convert(o *node, detectedWidth int) (*data, error) {
	//the total width available for a single line of help text
	width := o.renderWidth(detectedWidth)
	names := []string{}
	curr := o
	for curr != nil {
		names = append([]string{curr.name}, names...)
		curr = curr.parent
	}
	name := strings.Join(names, " ")
	args := make([]*datum, len(o.args))
	for i, arg := range o.args {
		//arguments are required
		n := "<" + arg.name + ">"
		//unless...
		if arg.slice {
			p := []string{arg.name, arg.name}
			for i, n := range p {
				if i < arg.min {
					//still required
					n = "<" + n + ">"
				} else {
					//optional!
					n = "[" + n + "]"
				}
				p[i] = n
			}
			n = strings.Join(p, " ") + " ..."
		}
		args[i] = &datum{
			Name: n,
			Help: constrain(arg.help, width),
		}
	}
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
			if item.shortName != "" && !o.flagSkipShort[item.name] {
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
	//get item help, with optional default values and env names and
	//constrain to a specific line width
	extras := make([]*template.Template, 3)
	keys := []string{"default", "env", "multiple"}
	for i, k := range keys {
		t, err := template.New("").Parse(o.templates["extra"+k])
		if err != nil {
			return nil, fmt.Errorf("template extra%s: %s", k, err)
		}
		extras[i] = t
	}
	//calculate the two column layout
	optionNameWidth, helpWidth, stacked := column(max, o.padWidth, width)
	indent := nletters(' ', optionNameWidth)
	//go back and render each option using calculated values
	for i, dg := range flagGroups {
		for j, to := range dg.Flags {
			if stacked {
				//help text drops onto its own line
				to.Pad = "\n" + indent
			} else {
				//pad all option names to be the same length
				to.Name += nletters(' ', max-len(to.Name))
			}
			//constrain help text
			item := o.flagGroups[i].flags[j]
			//render flag help string
			vals := []interface{}{item.defstr, item.envName, item.slice}
			outs := []string{}
			for i, v := range vals {
				b := strings.Builder{}
				if err := extras[i].Execute(&b, v); err != nil {
					return nil, err
				}
				if b.Len() > 0 {
					outs = append(outs, b.String())
				}
			}
			help := item.help
			extra := strings.Join(outs, ", ")
			if extra != "" {
				if help == "" {
					help = extra
				} else if trailingBrackets.MatchString(help) {
					m := trailingBrackets.FindStringSubmatch(help)
					help = m[1] + "(" + m[2] + ", " + extra + ")"
				} else {
					help += " (" + extra + ")"
				}
			}
			//align each row after the flag
			to.Help = wrap(help, helpWidth, indent)
		}
	}
	//commands - find max name length across all groups
	max = 0
	for _, s := range o.cmds {
		if l := len(s.name); l > max {
			max = l
		}
	}
	//commands are prefixed with a bullet, which shifts their columns across
	cmdNameWidth, cmdHelpWidth, cmdStacked := column(cmdPrefixWidth+max, o.padWidth, width)
	cmdIndent := nletters(' ', cmdNameWidth)
	//build command groups from o.cmdGroups (ordered)
	cmdGroups := make([]*datumGroup, len(o.cmdGroups))
	for gi, cg := range o.cmdGroups {
		//sort commands within each group alphabetically
		sorted := make([]*node, len(cg.cmds))
		copy(sorted, cg.cmds)
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].name < sorted[j].name
		})
		dg := &datumGroup{
			Name:  cg.name,
			Flags: make([]*datum, len(sorted)),
		}
		for i, s := range sorted {
			h := s.help
			if h == "" {
				h = s.summary
			}
			explicitMatch := o.cmdname != nil && *o.cmdname == s.name
			envMatch := o.cmdnameEnv != "" && os.Getenv(o.cmdnameEnv) == s.name
			if explicitMatch || envMatch {
				if h == "" {
					h = "default"
				} else {
					h += " (default)"
				}
			}
			d := &datum{
				Name: s.name,
				Help: wrap(h, cmdHelpWidth, cmdIndent),
				Pad:  nletters(' ', max-len(s.name)+o.padWidth),
			}
			if cmdStacked {
				d.Pad = "\n" + cmdIndent
			}
			dg.Flags[i] = d
		}
		cmdGroups[gi] = dg
	}
	//convert error to string
	err := ""
	if o.err != nil {
		err = wrap(o.err.Error(), width-o.padWidth, pad)
	}
	return &data{
		datum: datum{
			Name: name,
			Help: o.help,
			Pad:  pad,
		},
		Args:       args,
		FlagGroups: flagGroups,
		CmdGroups:  cmdGroups,
		Order:      o.order,
		Version:    o.version,
		Summary:    constrain(o.summary, width),
		Repo:       o.repo,
		Author:     o.author,
		ErrMsg:     err,
	}, nil
}
