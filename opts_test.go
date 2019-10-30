package opts

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestStrings(t *testing.T) {
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

func TestStrings2(t *testing.T) {
	//config
	type Config struct {
		Foo string
		Bar string
	}
	c := &Config{}
	//flag example parse
	err := testNew(c).parse([]string{"/bin/prog", "--foo", "hello", "--bar", "world with spaces"})
	if err != nil {
		t.Fatal(err)
	}
	//check config is filled
	check(t, c.Foo, "hello")
	check(t, c.Bar, "world with spaces")
}

func TestStrings3(t *testing.T) {
	type MyString string
	//config
	type Config struct {
		Foo string
		Bar MyString
	}
	c := &Config{}
	//flag example parse
	err := testNew(c).parse([]string{"/bin/prog", "--foo", "hello", "--bar", "world with spaces"})
	if err != nil {
		t.Fatal(err)
	}
	//check config is filled
	check(t, c.Foo, "hello")
	check(t, c.Bar, MyString("world with spaces"))
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

func TestBool(t *testing.T) {
	//config
	type Config struct {
		Foo string
		Bar bool
	}
	c := &Config{}
	//flag example parse
	err := testNew(c).parse([]string{"/bin/prog", "--foo", "hello", "--bar"})
	if err != nil {
		t.Fatal(err)
	}
	//check config is filled
	check(t, c.Foo, "hello")
	check(t, c.Bar, true)
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
	os.Setenv("ANOTHER_NUM", "21")
	//config
	type Config struct {
		Str        string
		Num        int
		Bool       bool
		AnotherNum int64
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
	check(t, c.AnotherNum, int64(21))
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

func TestShortSkip(t *testing.T) {
	type Config struct {
		Foo    string `opts:"short=f"`
		Bar    string `opts:"short=-"`
		Lalala string
	}
	c := &Config{}
	o, _ := New(c).Version("1.2.3").Name("skipshort").ParseArgsError([]string{"/bin/prog", "--help"})
	check(t, o.Help(), `
  Usage: skipshort [options]

  Options:
  --foo, -f
  --bar
  --lalala, -l
  --version, -v  display version
  --help, -h     display help

  Version:
    1.2.3

`)
}

func TestShortSkipConflictHelp(t *testing.T) {
	type Config struct {
		Foo    string `opts:"short=f"`
		Bar    string `opts:"short=-"`
		Hahaha string
	}
	c := &Config{}
	o, _ := New(c).Version("1.2.3").Name("skipshort").ParseArgsError([]string{"/bin/prog", "--help"})
	check(t, o.Help(), `
  Usage: skipshort [options]

  Options:
  --foo, -f
  --bar
  --hahaha, -h
  --version, -v  display version
  --help         display help

  Version:
    1.2.3

`)
}

func TestShortSkipInternal(t *testing.T) {
	type Config struct {
		Foo    string `opts:"short=f"`
		Bar    string `opts:"short=-"`
		Hahaha string `opts:"short=-"`
	}
	c := &Config{}
	o, _ := New(c).Version("1.2.3").Name("skipshort").ParseArgsError([]string{"/bin/prog", "--help"})
	check(t, o.Help(), `
  Usage: skipshort [options]

  Options:
  --foo, -f
  --bar
  --hahaha
  --version, -v  display version
  --help, -h     display help

  Version:
    1.2.3

`)
}

func TestJSON(t *testing.T) {
	//insert a config file
	p := filepath.Join(os.TempDir(), "opts.json")
	b := []byte(`{"foo":"hello", "bar":7, "faz":2, "boz": "test"}`)
	if err := ioutil.WriteFile(p, b, 0755); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(p)
	//parse flags
	type Config struct {
		Foo string
		Bar int
		Faz int
		Boz string
	}
	c := &Config{}
	//flag example parse
	n := testNew(c)
	n.ConfigPath(p)
	os.Setenv("FAZ", "0")
	n.UseEnv()
	if err := n.parse([]string{"/bin/prog", "--bar", "181", "--boz", ""}); err != nil {
		t.Fatal(err)
	}
	os.Unsetenv("FAZ")
	check(t, c.Foo, `hello`)
	check(t, c.Bar, 181) // JSON value overridden by command-line option parameter
	check(t, c.Faz, 0)
	check(t, c.Boz, ``)
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

func TestDocBrackets(t *testing.T) {
	//config
	type Config struct {
		Foo string `opts:"help=a message (submessage)"`
	}
	c := &Config{
		Foo: "bar",
	}
	//flag example parse
	o, _ := New(c).Name("docbrackets").ParseArgsError([]string{"/bin/prog", "--help"})
	check(t, o.Help(), `
  Usage: docbrackets [options]

  Options:
  --foo, -f   a message (submessage, default bar)
  --help, -h  display help

`)
}

func TestDocUseEnv(t *testing.T) {
	//config
	type Config struct {
		Foo string `opts:"help=a message"`
	}
	c := &Config{}
	//flag example parse
	o, _ := New(c).UseEnv().Version("1.2.3").Name("docuseenv").ParseArgsError([]string{"/bin/prog", "--help"})
	check(t, o.Help(), `
  Usage: docuseenv [options]

  Options:
  --foo, -f      a message (env FOO)
  --version, -v  display version
  --help, -h     display help

  Version:
    1.2.3

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

func TestCustomFlags(t *testing.T) {
	//config
	type Config struct {
		Foo *url.URL `opts:"help=my url"`
		Bar *url.URL `opts:"help=another url"`
	}
	c := Config{
		Foo: &url.URL{},
	}
	//flag example parse
	n := testNew(&c)
	if err := n.parse([]string{"/bin/prog", "-f", "http://foo.com"}); err != nil {
		t.Fatal(err)
	}
	if c.Foo == nil || c.Foo.String() != "http://foo.com" {
		t.Fatalf("incorrect foo: %v", c.Foo)
	}
	if c.Bar != nil {
		t.Fatal("bar should be nil")
	}
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
		t.Fatalf("got %s (%s), expected %s (%s)%s", stra, typea.Type(), strb, typeb.Type(), extra)
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
