package flag

import "testing"

func check(t *testing.T, a, b interface{}) {
	if a != b {
		t.Fatalf("got '%v', expected '%v'", a, b)
	}
}

func TestSimple(t *testing.T) {
	//application config
	type AppConfig struct {
		Foo string
		Bar string
	}

	c := &AppConfig{}

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

	//application config
	type AppConfig struct {
		Cmd string `cmd:"!"`
		//subcommand (external struct)
		Foo FooConfig
		//subcommand (inline struct)
		Bar struct {
			Zip string
			Zap string
		}
	}

	c := &AppConfig{}

	New(c).ParseArgs([]string{"bar", "--zip", "hello", "--zap", "world"})

	//check config is filled
	check(t, c.Cmd, "bar")
	check(t, c.Foo.Ping, "")
	check(t, c.Foo.Pong, "")
	check(t, c.Bar.Zip, "hello")
	check(t, c.Bar.Zap, "world")
}

func TestUnsupportedType(t *testing.T) {
	//application config
	type AppConfig struct {
		Foo string
		Bar interface{}
	}

	c := &AppConfig{}

	//flag example parse
	err := New(c).Process([]string{"--foo", "hello", "--bar", "world"})

	if err == nil {
		t.Fatal("Expected error")
	}
	if err.Error() != "Struct field 'Bar' has unsupported type: interface" {
		t.Fatal("Unexpected error type")
	}
}

func TestHelp(t *testing.T) {
	//application config
	type AppConfig struct {
		Foo string
		Bar string `help:"some help text"`
	}

	c := &AppConfig{}

	//flag example parse
	New(c).Name("zoop")

	//check config is filled
	// 	check(t, f.Help(), `Usage: zoop [options]

	// Options:
	// --foo
	// --bar some help text
	// `)
}
