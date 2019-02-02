package gen

import (
	"gencli/root"
)

type codeGen struct {
	Name string
}

func init() {
	root.Singleton().AddSubCmd("gen-cmd", &codeGen{})
}
