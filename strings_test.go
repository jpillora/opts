package opts

import (
	"reflect"
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
