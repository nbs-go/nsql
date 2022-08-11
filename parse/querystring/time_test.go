package qs_test

import (
	qs "github.com/nbs-go/nsql/parse/querystring"
	"testing"
)

func TestParseTime_EmptyString(t *testing.T) {
	_, actual := qs.ParseTime("")

	// Assert
	if actual != false {
		t.Errorf("FAILED\n  > got different values: %t", actual)
	}
}

func TestParseTime_Now(t *testing.T) {
	_, actual := qs.ParseTime("NOW")

	// Assert
	if actual != true {
		t.Errorf("FAILED\n  > got different values: %t", actual)
	}
}

func TestParseTime_UnixEpoch(t *testing.T) {
	dt, ok := qs.ParseTime("1640995200")

	// Assert
	actual := dt.Unix()
	expected := int64(1640995200)
	if !ok || actual != expected {
		t.Errorf("FAILED\n  > got different values: %d, expected: %d", actual, expected)
	}
}

func TestParseTime_RFC3999(t *testing.T) {
	dt, ok := qs.ParseTime("2022-01-01T00:00:00-00:00")

	// Assert
	actual := dt.Unix()
	expected := int64(1640995200)
	if !ok || actual != expected {
		t.Errorf("FAILED\n  > got different values: %d, expected: %d", actual, expected)
	}
}

func TestParseTime_CustomLayout(t *testing.T) {
	dt, ok := qs.ParseTime("2022-01-01 00:00:00", "2006-01-02 15:04:05")

	// Assert
	actual := dt.Unix()
	expected := int64(1640995200)
	if ok && actual != expected {
		t.Errorf("FAILED\n  > got different values: %d, expected: %d", actual, expected)
	}
}

func TestParseTime_InvalidInput(t *testing.T) {
	_, ok := qs.ParseTime("Not a date", "2006-01-02 15:04:05")

	// Assert
	if ok {
		t.Errorf("FAILED\n  > got different values: %t", ok)
	}
}
