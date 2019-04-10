package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jpillora/md-tmpl/modtmpl"
)

func main() {
	egs, err := ioutil.ReadDir(".")
	check(err)
	wg := sync.WaitGroup{}
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

		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			processGo(f)
			processReadme(f)
		}(eg)
	}
	wg.Wait()
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
	b2 := bytes.ReplaceAll(b,
		[]byte("go run main.go --help"),
		[]byte(fmt.Sprintf("go build -o %s && ./%s --help && rm %s", eg, eg, eg)),
	)
	if !bytes.Equal(b, b2) {
		check(ioutil.WriteFile(fp, b2, 0655))
		log.Printf("edited %s", fp)
	}
	proc := modtmpl.NewProcessor()
	proc.Write = true
	proc.ProcessFile(eg, path.Join(eg, "README.md"))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
