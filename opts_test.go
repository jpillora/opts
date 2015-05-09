package opts

import (
	"os"
	"testing"
)

func check(t *testing.T, a, b interface{}) {
	if a != b {
		t.Fatalf("got '%v', expected '%v'", a, b)
	}
}

func TestSimple(t *testing.T) {
	//config
	type Config struct {
		Foo string
		Bar string
	}

	c := &Config{}

	//flag example parse
	New(c).ParseArgs([]string{"--foo", "hello", "--bar", "world"})

	//check config is filled
	check(t, c.Foo, "hello")
	check(t, c.Bar, "world")
}

func TestSubCommand(t *testing.T) {

	type FooConfig struct {
		Ping string
		Pong string
	}

	//config
	type Config struct {
		Cmd string `type:"cmdname"`
		//subcommand (external struct)
		Foo FooConfig
		//subcommand (inline struct)
		Bar struct {
			Zip string
			Zap string
		}
	}

	c := &Config{}

	New(c).ParseArgs([]string{"bar", "--zip", "hello", "--zap", "world"})

	//check config is filled
	check(t, c.Cmd, "bar")
	check(t, c.Foo.Ping, "")
	check(t, c.Foo.Pong, "")
	check(t, c.Bar.Zip, "hello")
	check(t, c.Bar.Zap, "world")
}

func TestUnsupportedType(t *testing.T) {
	//config
	type Config struct {
		Foo string
		Bar interface{}
	}

	c := &Config{}

	//flag example parse
	err := New(c).Process([]string{"--foo", "hello", "--bar", "world"})

	if err == nil {
		t.Fatal("Expected error")
	}
	check(t, err.Error(), "Struct field 'Bar' has unsupported type: interface")
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
	New(c).UseEnv().Parse()

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
	New(c).ParseArgs([]string{"-b", "wld", "hel", "lo"})

	//check config is filled
	check(t, c.Foo, `hel`)
	check(t, c.Zip, `lo`)
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
	err := New(c).Process([]string{"-f", "1", "-b", "2"})
	if err == nil {
		t.Fatal("expected error")
	}
}
