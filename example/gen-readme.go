package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jpillora/md-tmpl/mdtmpl"
)

func main() {
	egs, err := ioutil.ReadDir(".")
	check(err)
	for _, s := range egs {
		eg := s.Name()
		if !s.IsDir() || !strings.HasPrefix(eg, "eg-") {
			continue
		}
		f := ""
		if len(os.Args) >= 2 {
			f = os.Args[1]
		}
		if !strings.Contains(eg, f) {
			continue
		}
		check(err)
		processGo(eg)
		processReadme(eg)
	}
}

func processGo(eg string) {
	b, err := ioutil.ReadFile(filepath.Join(eg, "main.go"))
	if err != nil {
		log.Printf("example '%s' has no main.go file", eg)
		return
	}
	if len(b) == 0 {
		log.Fatalf("example '%s' has empty main.go file", eg)
	}
}

func processReadme(eg string) {
	fp := filepath.Join(eg, "README.md")
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		log.Printf("example '%s' has no README.md file", eg)
		return
	}
	b = bytes.ReplaceAll(b,
		[]byte("go run main.go --help"),
		[]byte(fmt.Sprintf("go build -o %s && ./%s --help && rm %s", eg, eg, eg)),
	)
	b = mdtmpl.ExecuteIn(b, eg)
	check(ioutil.WriteFile(fp, b, 0655))
	log.Printf("executed templates and rewrote '%s'", eg)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
