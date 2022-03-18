package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestCompiler(t *testing.T) {
	var tests = []struct {
		s   string
		err string
	}{
		{
			s:   "Allow group blaha to use vcns in bar",
			err: "",
		},
		{
			s:   "Foo",
			err: "1:1: syntax error: unexpected VARIABLE, expecting ALLOW",
		},
		{
			s:   "",
			err: "1:0: syntax error: unexpected $end, expecting ALLOW",
		},
		{
			s:   "Allow group blaha to use vcns in bar where foo != 42",
			err: "",
		},
		{
			s:   "Allow group blaha to use vcns in bar where foo > 4",
			err: "1:49: syntax error: unexpected VALUE, expecting EQ or NE",
		},
	}
	comp := NewCompiler()
	for i, tt := range tests {
		err := comp.Parse(strings.NewReader(tt.s))
		if !reflect.DeepEqual(tt.err, errstring(err)) {
			t.Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, err)
		}
	}
}

// errstring returns the string representation of an error.
func errstring(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
