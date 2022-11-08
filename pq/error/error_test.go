package pqErr_test

import (
	"errors"
	"github.com/lib/pq"
	pqErr "github.com/nbs-go/nsql/pq/error"
	"testing"
)

func TestUnknownError(t *testing.T) {
	// Prepare test case
	sErr := errors.New("not a pq.Error")
	expected := pqErr.UnknownError
	// Do test
	actual, _ := pqErr.Parse(sErr)
	if actual != expected {
		t.Errorf("unexpected test result. Expected = %d, Actual = %d", expected, actual)
	}
}

func TestUniqueError(t *testing.T) {
	// Prepare test case
	sErr := &pq.Error{
		Code:       "23505",
		Message:    "This is a unique error",
		Constraint: "idx_Vehicle_unique",
	}
	expected := pqErr.UniqueError
	// Do test
	actual, _ := pqErr.Parse(sErr)
	if actual != expected {
		t.Errorf("unexpected test result. Expected = %d, Actual = %d", expected, actual)
	}
}

func TestFkViolationError(t *testing.T) {
	// Prepare test case
	sErr := &pq.Error{
		Code:       "23503",
		Message:    "This is an fk violations error",
		Constraint: "fk_Vehicle_Company",
	}
	expected := pqErr.FkViolationError
	// Do test
	actual, _ := pqErr.Parse(sErr)
	if actual != expected {
		t.Errorf("unexpected test result. Expected = %d, Actual = %d", expected, actual)
	}
}

func TestUnhandledError(t *testing.T) {
	// Prepare test case
	sErr := &pq.Error{
		Code:       "0",
		Message:    "This is an unhandled error",
		Constraint: "Constraint",
	}
	expected := pqErr.UnhandledError
	// Do test
	actual, meta := pqErr.Parse(sErr)
	if actual != expected {
		t.Errorf("unexpected test result. Expected = %d, Actual = %d", expected, actual)
	}

	// Do test on metadata
	if meta.Constraint == "" {
		t.Errorf("unexpected test result. metadata.Constraint is empty")
	}
	if meta.Message == "" {
		t.Errorf("unexpected test result. metadata.Message is empty")
	}
}
