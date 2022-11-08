package pqErr

import (
	"errors"
	"github.com/lib/pq"
)

func Parse(err error) (Code, *Metadata) {
	// Cast as pq error
	pErr := new(pq.Error)
	if !errors.As(err, &pErr) {
		return UnknownError, nil
	}

	// Set metadata value
	meta := &Metadata{
		Constraint: pErr.Constraint,
		Message:    pErr.Message,
	}

	switch pErr.Code {
	case "23505":
		return UniqueError, meta
	case "23503":
		return FkViolationError, meta
	default:
		return UnhandledError, meta
	}
}
