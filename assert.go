package httpclientmock

import (
	"bytes"
	"testing"
)

// AssertEqual asserts that two values are the same.
//
// It only needs to compare strings and slices of bytes.
//
// This function can be overridden to adjust to other assertion frameworks like testify
var AssertEqual = func(t *testing.T, expected, actual interface{}, format string, args ...interface{}) {
	bexpected, eok := expected.([]byte)
	bactual, aok := actual.([]byte)

	if eok && aok {
		res := bytes.Compare(bexpected, bactual)
		if res != 0 {
			Errorf(t, format, args...)
		}
		return
	}

	if expected != actual {
		Errorf(t, format, args...)
	}
}

// Errorf marks an error in the test
//
// This function can be overridden to make use of another assertion framework like testify
var Errorf = func(t *testing.T, format string, args ...interface{}) {
	t.Errorf(format, args...)
}
