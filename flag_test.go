package flag

import (
	"fmt"
	"testing"
)

func TestSimple(t *testing.T) {

	//application config
	type AppConfig struct {
		Foo string `help:"Foo does stuff"`
		Bar string `help:"Bar does stuff"`
	}

	c := &AppConfig{}

	//flag example parse
	f := /*flag.*/ New(c)
	f.Version = "1.1.0"
	f.Args = []string{"--foo", "hello", "--bar", "world"} //replace os.Args
	f.Parse()

	fmt.Println(c.Foo)
	fmt.Println(c.Bar)
	//check config is filled
	if c.Foo != "hello" {
		t.Fatal("Foo should be 'hello'")
	}
	if c.Bar != "world" {
		t.Fatal("Bar should be 'world'")
	}
}

func TestSubCommand(t *testing.T) {

	type FooConfig struct {
		Ping string `opt`
		Pong string `opt`
	}

	// type BarConfig

	//application config
	type AppConfig struct {
		Foo FooConfig
		Bar struct {
			Zip string `opt`
			Zap string `opt`
		}
		// Bazz struct{}
	}

	New(&AppConfig{})

	// c := AppConfig{}
}
