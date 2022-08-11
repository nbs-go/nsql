package qs_test

import (
	qs "github.com/nbs-go/nsql/parse/querystring"
	"testing"
)

func TestParseInt_EmptyString(t *testing.T) {
	_, actual := qs.ParseInt("")

	// Assert
	if actual != false {
		t.Errorf("FAILED\n  > got different values: %t", actual)
	}
}

func TestParseInt(t *testing.T) {
	actual, ok := qs.ParseInt("666")

	// Assert
	expected := int64(666)
	if !ok || actual != expected {
		t.Errorf("FAILED\n  > got different values: %d, expected: %d", actual, expected)
	}
}

func TestParseInt_Invalid(t *testing.T) {
	_, actual := qs.ParseInt("NOT AN INT")

	// Assert
	if actual != false {
		t.Errorf("FAILED\n  > got different values: %t", actual)
	}
}

func TestParseFloat_EmptyString(t *testing.T) {
	_, actual := qs.ParseFloat("")

	// Assert
	if actual != false {
		t.Errorf("FAILED\n  > got different values: %t", actual)
	}
}

func TestParseFloat(t *testing.T) {
	actual, ok := qs.ParseFloat("10.5")

	// Assert
	expected := 10.5
	if !ok || actual != expected {
		t.Errorf("FAILED\n  > got different values: %f, expected: %f", actual, expected)
	}
}

func TestParseFloat_Invalid(t *testing.T) {
	_, actual := qs.ParseFloat("NOT A FLOAT")

	// Assert
	if actual != false {
		t.Errorf("FAILED\n  > got different values: %t", actual)
	}
}
