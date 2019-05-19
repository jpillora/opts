package initopts

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/jpillora/opts"
)

type initOpts struct {
	SrcControlHost string
	Owner          string `opts:"mode=arg"`
	Name           string `opts:"mode=arg"`
	Package        string
	Command        string
	Directory      string `opts:"help=output directory"`
}

func Register(parent opts.Opts) opts.Opts {
	in := initOpts{
		SrcControlHost: "github.com",
		Directory:      ".",
	}
	parent.AddCommand(opts.New(&in).Name("init"))
	return parent
}

func (in *initOpts) Run() error {
	dir, err := os.OpenFile(in.Directory, os.O_APPEND, 0755)
	if err != nil {
		err = os.MkdirAll(in.Directory, 0755)
		if err != nil {
			return err
		}
	} else {
		names, err := dir.Readdirnames(1)
		if len(names) > 0 {
			return errors.New("output directory not empty")
		}
		if err != io.EOF {
			return err
		}
	}
	// if err := f.Close(); err != nil {
	// 	log.Fatal(err)
	// }

	data := struct {
		Module  string
		Command string
		Name    string
		Owner   string
	}{
		Module:  in.SrcControlHost + "/" + in.Owner + "/" + in.Name,
		Command: in.Name,
		Name:    in.Name,
		Owner:   in.Owner,
	}
	if in.Package != "" {
		data.Module = data.Module + "/" + in.Package
	}
	if in.Command != "" {
		data.Command = in.Command
	} else if in.Package != "" {
		data.Command = in.Package
	}
	fmt.Printf("#init %+v\n", data)
	for _, fi := range files {
		tmpl, err := template.New(fi.Path).Parse(fi.Tmpl)
		if err != nil {
			fmt.Printf("tmpl parse error : %v\n", err)
			continue
		}
		fmt.Printf("#%v\n", fi.Path)
		pa := filepath.Join(in.Directory, path.Dir(fi.Path))
		_ = os.MkdirAll(pa, 0755)
		// if err != nil {
		// 	return err
		// }
		pa = filepath.Join(pa, path.Base(fi.Path))
		ofi, err := os.OpenFile(pa, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			fmt.Printf("new file error: %v", err)
			continue
		}
		err = tmpl.Execute(ofi, data)
		if err != nil {
			fmt.Printf("tmpl exec error : %v\n", err)
			continue
		}
	}
	return nil
}

type file struct {
	Path string
	Tmpl string
}

var files = []file{
	{
		Path: "go.mod",
		Tmpl: `module {{.Module}}

go 1.12

require github.com/jpillora/opts v1.0.0
`,
	},
	{
		Path: "main.go",
		Tmpl: `package main

import (
	"fmt"
	"os"

	"github.com/jpillora/opts"
	"{{.Module}}/internal/initopts"
)

var (
	Version string = "dev"
	Date    string = "na"
	Commit  string = "na"
)

type root struct {
	parsedOpts opts.ParsedOpts
}

func main() {
	r := root{}
	ro := opts.New(&r).Name("{{.Name}}").
		EmbedGlobalFlagSet().
		Complete().
		Version(Version)

	initopts.Register(ro)

	r.parsedOpts = ro.Parse()
	err := r.parsedOpts.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "run error %v\n", err)
		os.Exit(2)
	}
}

func (rt *root) Run() {
	fmt.Printf("Version: %s\nDate: %s\nCommit: %s\n", Version, Date, Commit)
	fmt.Println(rt.parsedOpts.Help())
}
`,
	},
	{
		Path: "internal/initopts/init.go",
		Tmpl: `package initopts

import (
	"fmt"

	"github.com/jpillora/opts"
)

type initOpts struct {
}

func Register(parent opts.Opts) opts.Opts {
	in := initOpts{	}
	parent.AddCommand(opts.New(&in).Name("init"))
	return parent
}

func (in *initOpts) Run() error {
	fmt.Printf("#init %+v\n", in)
	return nil
}

`,
	},
	{
		Path: ".goreleaser.yml",
		Tmpl: `# This is an example goreleaser.yaml file with some sane defaults.
# Make sure to check the documentation at http://goreleaser.com
project_name: {{.Name}}
release:
  github:
    owner: {{.Owner}} 
    name: {{.Name}}
  name_template: '{{"{{"}}.Tag}}'
  # disable: true

builds:
- 
  binary: {{.Command}}
  env:
  - CGO_ENABLED=0
  goos:
  - linux
  - darwin
  - windows
  goarch:
  - amd64
  - "386"
  ignore:
  - goos: darwin
    goarch: 386
  main: .
  ldflags:
  - -s -w -X main.Version={{"{{"}}.Version}} -X main.Commit={{"{{"}}.Commit}} -X main.Date={{"{{"}}.Date}}

archive:
  replacements:
    386: i386
    amd64: x86_64
  format_overrides:
  - goos: windows
    format: zip
  files:
  - licence*
  - LICENCE*
  - license*
  - LICENSE*
  - readme*
  - README*
  - changelog*
  - CHANGELOG*

checksum:
  name_template: 'checksums.txt'

snapshot:
  name_template: "{{"{{"}} .Tag {{"}}"}}-next"

changelog:
  sort: asc
  filters:
    exclude:
    - '^docs:'
    - '^test:'
`,
	},
	{
		Path: ".gitignore",
		Tmpl: `{{.Command}}
dist/`,
	},
}
