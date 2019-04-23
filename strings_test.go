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
			"",
			nil,
		},
	} {
		m := newKV(testcase.input)
		if !reflect.DeepEqual(m, testcase.output) {
			t.Fatalf("input: %s\n  expected: %s\n  got: %s",
				testcase.input,
				testcase.output,
				m,
			)
		}
	}

}
