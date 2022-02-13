package test_utils

import (
	"strings"
	"testing"
)

func CompareStringArray(t *testing.T, expectation string, actual, expected []string) {
	s1 := strings.Join(actual, ", ")
	s2 := strings.Join(expected, ", ")

	if s1 != s2 {
		t.Errorf("%s: FAILED\n  > got different values: %s", expectation, s1)
	} else {
		t.Logf("%s: PASSED", expectation)
	}
}
func CompareInterfaceArray(t *testing.T, expectation string, actual, expected []interface{}) {
	if len(actual) != len(expected) {
		t.Errorf("%s: FAILED\n  > got different values: %v", expectation, actual)
		return
	}

	for i, v := range actual {
		if expected[i] != v {
			t.Errorf("%s: FAILED\n  > got different values: %v", expectation, actual)
			return
		}
	}

	t.Logf("%s: PASSED", expectation)
}

func CompareStringIn(t *testing.T, expectation string, actual string, expected []string) {
	for _, exp := range expected {
		if exp == actual {
			t.Logf("%s: PASSED", expectation)
			return
		}
	}
	t.Errorf("%s: FAILED\n  > got different values: %s", expectation, actual)
}

func CompareString(t *testing.T, expectation string, actual, expected string) {
	if actual != expected {
		t.Errorf("%s: FAILED\n  > got different values: %s", expectation, actual)
	} else {
		t.Logf("%s: PASSED", expectation)
	}
}

func CompareBoolean(t *testing.T, expectation string, actual, expected bool) {
	if actual != expected {
		t.Errorf("%s: FAILED\n  > got different values: %t", expectation, actual)
	} else {
		t.Logf("%s: PASSED", expectation)
	}
}

func CompareInt(t *testing.T, expectation string, actual, expected int) {
	if actual != expected {
		t.Errorf("%s: FAILED\n  > got different values: %d", expectation, actual)
	} else {
		t.Logf("%s: PASSED", expectation)
	}
}

func RecoverPanic(t *testing.T, expectation string, errStr string) func() {
	return func() {
		r := recover()
		if r == nil {
			t.Errorf("%s: FAILED\n  > code did not panic", expectation)
			return
		}

		err, ok := r.(error)
		if !ok {
			t.Errorf("%s: FAILED\n  > unknown recovered value: %v", expectation, r)
			return
		}

		if err.Error() != errStr {
			t.Errorf("%s: FAILED\n  > got different error: %v", expectation, err)
		} else {
			t.Logf("%s: PASSED", expectation)
		}
	}
}
