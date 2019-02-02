package main

import (
	"gencli/root"

	_ "gencli/gen"
)

func main() {
	root.Singleton().Parse()
}
