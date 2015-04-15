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
	f := New(c)
	f.Args = []string{"--foo", "hello", "--bar", "world"} //replace os.Args
	f.Parse()

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
		//struct ptr (automatically new'd)
		Foo *FooConfig
		//inline struct
		Bar struct {
			Zip string
			Zap string
		}
	}

	c := &AppConfig{}

	f := New(c)
	f.Args = []string{"bar", "--zip", "hello", "--zap", "world"}
	f.Parse()

	//check config is filled
	check(t, c.Foo.Ping, "")
	check(t, c.Foo.Pong, "")
	check(t, c.Bar.Zip, "hello")
	check(t, c.Bar.Zap, "world")
}

func TestHelp(t *testing.T) {
	//application config
	type AppConfig struct {
		Foo string
		Bar string `help:"some help text"`
	}

	c := &AppConfig{}

	//flag example parse
	f := New(c)
	f.Name = "zoop"

	//check config is filled
	check(t, f.Help(), `Usage: zoop [options]

Options:
--foo
--bar some help text
`)
}
