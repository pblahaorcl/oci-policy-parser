//go:build fuzzy
// +build fuzzy

package main

import (
	"strings"
	"testing"
)

func FuzzParserTest(f *testing.F) {
	f.Add("Allow group blaha to use vcns in bar where foo != 42")
	f.Add("Allow group blaha to")

	f.Fuzz(func(t *testing.T, s string) {
		c := NewCompiler()
		err := c.Parse(strings.NewReader(s))
		if err != nil {
			switch e := err.(type) {
			case *SyntaxError:
				// Ignore as the correct syntax error
			default:
				t.Errorf("%v", e)
			}
		}
	})
}
