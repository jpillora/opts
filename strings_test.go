package opts

import (
	"reflect"
	"strings"
	"testing"
)

func TestKVMap(t *testing.T) {
	for _, testcase := range []struct {
		input  string
		output map[string]string
	}{
		{
			"a=b,c=d",
			map[string]string{"a": "b", "c": "d"},
		},
		{
			"foo,,bar,,",
			map[string]string{"foo": "", "bar": ""},
		},
		{
			"ping=,,pong==,,",
			map[string]string{"ping": "", "pong": "="},
		},
		{
			"nospace=,,  leadingspace==,  trailingspace  ,",
			map[string]string{"nospace": "", "leadingspace": "=", "trailingspace  ": ""},
		},
	} {
		kv := newKV(testcase.input)
		m := kv.m
		if !reflect.DeepEqual(m, testcase.output) {
			t.Fatalf("input: %s\n  expected: %s\n      got: %s",
				testcase.input,
				testcase.output,
				m,
			)
		}
	}
}

func TestConstrain(t *testing.T) {
	for _, testcase := range []struct {
		input  string
		width  int
		output string
	}{
		//exactly the width is not wrapped
		{"hello worl", 10, "hello worl"},
		//one over is
		{"hello world", 10, "hello\nworld"},
		//words longer than the width overflow rather than being broken
		{"a supercalifragilistic word", 10, "a\nsupercalifragilistic\nword"},
		//existing newlines are preserved
		{"one\ntwo three four", 8, "one\ntwo\nthree\nfour"},
		//consecutive spaces are preserved
		{"a  b", 10, "a  b"},
		//a zero width disables wrapping
		{"hello world", 0, "hello world"},
		{
			"the address to listen on. it may be a host or a port or both",
			24,
			"the address to listen\non. it may be a host or\na port or both",
		},
	} {
		got := constrain(testcase.input, testcase.width)
		if got != testcase.output {
			t.Fatalf("input: %q (width %d)\n  expected: %q\n       got: %q",
				testcase.input,
				testcase.width,
				testcase.output,
				got,
			)
		}
		//no line may exceed the width, unless it is a single long word
		for _, line := range strings.Split(got, "\n") {
			if testcase.width > 0 && len(line) > testcase.width && strings.Contains(line, " ") {
				t.Fatalf("input: %q (width %d)\n  line too long: %q",
					testcase.input, testcase.width, line)
			}
		}
	}
}

func TestCamel2Dash(t *testing.T) {
	for _, testcase := range []struct {
		input  string
		output string
	}{
		{
			"fooBar",
			"foo-bar",
		},
		{
			"WordACRONYMAnotherWord",
			"word-acronym-another-word",
		},
		{
			"IDs",
			"ids",
		},
		{
			"URLs",
			"urls",
		},
		{
			"HTTPServer",
			"http-server",
		},
		{
			"GetURLPath",
			"get-url-path",
		},
		{
			"APIKey",
			"api-key",
		},
		{
			"GetUIElement",
			"get-ui-element",
		},
		{
			"ID",
			"id",
		},
		{
			"MyDBs",
			"my-dbs",
		},
	} {
		got := camel2dash(testcase.input)
		if testcase.output != got {
			t.Fatalf("input: %s\n  expected: %s\n       got: %s",
				testcase.input,
				testcase.output,
				got,
			)
		}
	}

}
