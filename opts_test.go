package opts

import (
	"fmt"
	"os"
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
	err := New(c).parse([]string{"--foo", "hello", "--bar", "world"})
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
	err := New(c).parse([]string{"--foo", "hello", "--foo", "world", "--bar", "ping", "--bar", "pong"})
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
		Cmd string `type:"cmdname"`
		//command (external struct)
		Foo FooConfig `type:"cmd"`
		//command (inline struct)
		Bar struct {
			Zip string
			Zap string
		} `type:"cmd"`
	}
	c := &Config{}
	err := New(c).parse([]string{"bar", "--zip", "hello", "--zap", "world"})
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
	err := New(c).parse([]string{"--zip", "hello", "--pong", "world"})
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
	err := New(&c).parse([]string{"--foo", "hello", "--bar", "world"})
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
	err := New(&c).parse([]string{"--foo", "hello", "--bar", "world"})
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
	if err := New(c).UseEnv().parse([]string{}); err != nil {
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

func TestArg(t *testing.T) {
	//config
	type Config struct {
		Foo string `type:"arg"`
		Zip string `type:"arg"`
		Bar string
	}
	c := &Config{}
	//flag example parse
	if err := New(c).UseEnv().parse([]string{"-b", "wld", "hel", "lo"}); err != nil {
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
		Zip string   `type:"arg"`
		Foo []string `type:"arg"`
		Bar string
	}
	c := &Config{}
	//flag example parse
	if err := New(c).UseEnv().parse([]string{"-b", "wld", "!!!", "hel", "lo"}); err != nil {
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
	err := New(c).parse([]string{"-f", "1", "-b", "2"})
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
  --help, -h

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
  --help, -h

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
	o := New(c).Name("groups").ParseArgs([]string{})
	check(t, o.Help(), `
  Usage: groups [options]

  Options:
  --fizz, -f
  --buzz, -b
  --help, -h

  More options:
  --ping, -p
  --pong

`)
}

func TestDocArgList(t *testing.T) {
	//config
	type Config struct {
		Foo string   `type:"arg"`
		Bar []string `type:"arg"`
	}
	c := &Config{}
	//flag example parse
	o := New(c).Name("").ParseArgs([]string{"zzz"})
	check(t, o.Help(), `TODO`)
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
