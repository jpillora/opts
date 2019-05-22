package opts

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestSimple(t *testing.T) {
	//config
	type Config struct {
		Foo string
		Bar string
	}
	c := &Config{}
	//flag example parse
	err := testNew(c).parse([]string{"/bin/prog", "--foo", "hello", "--bar", "world"})
	if err != nil {
		t.Fatal(err)
	}
	//check config is filled
	check(t, c.Foo, "hello")
	check(t, c.Bar, "world")
}

func TestList(t *testing.T) {
	//config
	type Config struct {
		Foo []string
		Bar []string
	}
	c := &Config{}
	//flag example parse
	err := testNew(c).parse([]string{"/bin/prog", "--foo", "hello", "--foo", "world", "--bar", "ping", "--bar", "pong"})
	if err != nil {
		t.Fatal(err)
	}
	//check config is filled
	check(t, c.Foo, []string{"hello", "world"})
	check(t, c.Bar, []string{"ping", "pong"})
}

func TestSubCommand(t *testing.T) {
	//subconfig
	type FooConfig struct {
		Ping string
		Pong string
	}
	//config
	type Config struct {
		Cmd string `opts:"mode=cmdname"`
		//command (external struct)
		Foo FooConfig `opts:"mode=cmd"`
		//command (inline struct)
		Bar struct {
			Zip string
			Zap string
		} `opts:"mode=cmd"`
	}
	c := &Config{}
	err := testNew(c).parse([]string{"/bin/prog", "bar", "--zip", "hello", "--zap", "world"})
	if err != nil {
		t.Fatal(err)
	}
	//check config is filled
	check(t, c.Cmd, "bar")
	check(t, c.Foo.Ping, "")
	check(t, c.Foo.Pong, "")
	check(t, c.Bar.Zip, "hello")
	check(t, c.Bar.Zap, "world")
}

func TestEmbed(t *testing.T) {
	type Foo struct {
		Ping string
		Pong string
	}
	type Bar struct {
		Zip string
		Zap string
	}
	//config
	type Config struct {
		Foo
		Bar
	}
	c := &Config{}
	err := testNew(c).parse([]string{"/bin/prog", "--zip", "hello", "--pong", "world"})
	if err != nil {
		t.Fatal(err)
	}
	//check config is filled
	check(t, c.Bar.Zap, "")
	check(t, c.Bar.Zip, "hello")
	check(t, c.Foo.Ping, "")
	check(t, c.Foo.Pong, "world")
}

func TestUnsupportedType(t *testing.T) {
	//config
	type Config struct {
		Foo string
		Bar map[string]bool
	}
	c := Config{}
	//flag example parse
	err := testNew(&c).parse([]string{"/bin/prog", "--foo", "hello", "--bar", "world"})
	if err == nil {
		t.Fatal("Expected error")
	}
	if !strings.Contains(err.Error(), "field type not supported: map") {
		t.Fatalf("Expected unsupported map, got: %s", err)
	}
}

func TestUnsupportedInterfaceType(t *testing.T) {
	//config
	type Config struct {
		Foo string
		Bar interface{}
	}
	c := Config{}
	//flag example parse
	err := testNew(&c).parse([]string{"/bin/prog", "--foo", "hello", "--bar", "world"})
	if err == nil {
		t.Fatal("Expected error")
	}
	if !strings.Contains(err.Error(), "field type not supported: interface") {
		t.Fatalf("Expected unsupported interface, got: %s", err)
	}
}

func TestEnv(t *testing.T) {
	os.Setenv("STR", "helloworld")
	os.Setenv("NUM", "42")
	os.Setenv("BOOL", "true")
	//config
	type Config struct {
		Str  string
		Num  int
		Bool bool
	}
	c := &Config{}
	//flag example parse
	n := testNew(c)
	n.UseEnv()
	if err := n.parse([]string{}); err != nil {
		t.Fatal(err)
	}
	os.Unsetenv("STR")
	os.Unsetenv("NUM")
	os.Unsetenv("BOOL")
	//check config is filled
	check(t, c.Str, `helloworld`)
	check(t, c.Num, 42)
	check(t, c.Bool, true)
}

func TestLongClash(t *testing.T) {
	type Config struct {
		Foo string
		Fee string `opts:"name=foo"`
	}
	c := &Config{}
	//flag example parse
	n := testNew(c)
	if err := n.parse([]string{}); err == nil {
		t.Fatal("expected error")
	} else if !strings.Contains(err.Error(), "already exists") {
		t.Fatal("expected already exists error")
	}
}

func TestShortClash(t *testing.T) {
	type Config struct {
		Foo string `opts:"short=f"`
		Fee string `opts:"short=f"`
	}
	c := &Config{}
	//flag example parse
	n := testNew(c)
	if err := n.parse([]string{}); err == nil {
		t.Fatal("expected error")
	} else if !strings.Contains(err.Error(), "already exists") {
		t.Fatal("expected already exists error")
	}
}

func TestJSON(t *testing.T) {
	//insert a config file
	p := filepath.Join(os.TempDir(), "opts.json")
	b := []byte(`{"foo":"hello", "bar":7}`)
	if err := ioutil.WriteFile(p, b, 0755); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(p)
	//parse flags
	type Config struct {
		Foo string
		Bar int
	}
	c := &Config{}
	//flag example parse
	n := testNew(c)
	n.ConfigPath(p)
	if err := n.parse([]string{"/bin/prog", "--bar", "8"}); err != nil {
		t.Fatal(err)
	}
	check(t, c.Foo, `hello`)
	check(t, c.Bar, 7) //currently uses JSON value... might change...
}

func TestArg(t *testing.T) {
	//config
	type Config struct {
		Foo string `opts:"mode=arg"`
		Zip string `opts:"mode=arg"`
		Bar string
	}
	c := &Config{}
	//flag example parse
	if err := testNew(c).parse([]string{"/bin/prog", "-b", "wld", "hel", "lo"}); err != nil {
		t.Fatal(err)
	}
	//check config is filled
	check(t, c.Foo, `hel`)
	check(t, c.Zip, `lo`)
	check(t, c.Bar, `wld`)
}

func TestArgs(t *testing.T) {
	//config
	type Config struct {
		Zip string   `opts:"mode=arg"`
		Foo []string `opts:"mode=arg"`
		Bar string
	}
	c := &Config{}
	//flag example parse
	if err := testNew(c).parse([]string{"/bin/prog", "-b", "wld", "!!!", "hel", "lo"}); err != nil {
		t.Fatal(err)
	}
	//check config is filled
	check(t, c.Zip, `!!!`)
	check(t, c.Foo, []string{`hel`, `lo`})
	check(t, c.Bar, `wld`)
}

func TestIgnoreUnexported(t *testing.T) {
	//config
	type Config struct {
		Foo string
		bar string
	}
	c := &Config{}
	//flag example parse
	err := testNew(c).parse([]string{"/bin/prog", "-f", "1", "-b", "2"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDocBefore(t *testing.T) {
	//config
	type Config struct {
		Foo string
		bar string
	}
	c := &Config{}
	//flag example parse
	o := New(c).Name("doc-before")
	n := o.(*node)
	l := len(n.order)
	o.DocBefore("usage", "mypara", "hello world this some text\n\n")
	op := o.ParseArgs(nil)
	check(t, len(n.order), l+1)
	check(t, op.Help(), `
  hello world this some text

  Usage: doc-before [options]

  Options:
  --foo, -f
  --help, -h  display help

`)
}

func TestDocAfter(t *testing.T) {
	//config
	type Config struct {
		Foo string
		bar string
	}
	c := &Config{}
	//flag example parse
	o := New(c).Name("doc-after")
	n := o.(*node)
	l := len(n.order)
	o.DocAfter("usage", "mypara", "\nhello world this some text\n")
	op := o.ParseArgs(nil)
	check(t, len(n.order), l+1)
	check(t, op.Help(), `
  Usage: doc-after [options]

  hello world this some text

  Options:
  --foo, -f
  --help, -h  display help

`)
}

func TestDocGroups(t *testing.T) {
	//config
	type Config struct {
		Fizz       string
		Buzz       bool
		Ping, Pong int `opts:"group=More"`
	}
	c := &Config{}
	//flag example parse
	o := New(c).Name("groups").ParseArgs(nil)
	check(t, o.Help(), `
  Usage: groups [options]

  Options:
  --fizz, -f
  --buzz, -b
  --help, -h  display help

  More options:
  --ping, -p
  --pong

`)
}

func TestDocArgList(t *testing.T) {
	//config
	type Config struct {
		Foo string   `opts:"mode=arg"`
		Bar []string `opts:"mode=arg"`
	}
	c := &Config{}
	//flag example parse
	o := New(c).Name("docargs").ParseArgs([]string{"/bin/prog", "zzz"})
	check(t, o.Help(), `
  Usage: docargs [options] <foo> [bar] [bar] ...

  Options:
  --help, -h  display help

`)
}

func TestSubCommandMap(t *testing.T) {
	//config
	type Config struct {
		Foo string
		bar string
	}
	c := Config{
		Foo: "foo",
	}
	New(&c)
}

var spaces = regexp.MustCompile(`\ `)
var newlines = regexp.MustCompile(`\n`)

func readable(s string) string {
	s = spaces.ReplaceAllString(s, "•")
	s = newlines.ReplaceAllString(s, "⏎\n")
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = fmt.Sprintf("%5d: %s", i+1, l)
	}
	s = strings.Join(lines, "\n")
	return s
}

func check(t *testing.T, a, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		stra := readable(fmt.Sprintf("%v", a))
		strb := readable(fmt.Sprintf("%v", b))
		typea := reflect.ValueOf(a)
		typeb := reflect.ValueOf(b)
		extra := ""
		if out, ok := diffstr(stra, strb); ok {
			extra = "\n\n" + out
			stra = "\n" + stra + "\n"
			strb = "\n" + strb + "\n"
		} else {
			stra = "'" + stra + "'"
			strb = "'" + strb + "'"
		}
		t.Fatalf("got %s (%s), expected %s (%s)%s", stra, typea.Kind(), strb, typeb.Kind(), extra)
	}
}

func diffstr(a, b interface{}) (string, bool) {
	stra, oka := a.(string)
	strb, okb := b.(string)
	if !oka || !okb {
		return "", false
	}
	ra := []rune(stra)
	rb := []rune(strb)
	line := 1
	char := 1
	var diff rune
	for i, a := range ra {
		if a == '\n' {
			line++
			char = 1
		} else {
			char++
		}
		var b rune
		if i < len(rb) {
			b = rb[i]
		}
		if a != b {
			a = diff
			break
		}
	}
	return fmt.Sprintf("Diff on line %d char %d (%d)", line, char, diff), true
}

func testNew(config interface{}) *node {
	o := New(config)
	n := o.(*node)
	return n
}
