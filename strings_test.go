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
