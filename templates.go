package flag

var defaultTemplates = map[string]string{
	"help": `{{.Name}} [options]`,
}

// tmpl, err := template.New("test").Parse(help)
// if err != nil { panic(err) }
// err = tmpl.Execute(os.Stdout, sweaters)
// if err != nil { panic(err) }

func (f *Flag) Help() string {
	return ""
}
