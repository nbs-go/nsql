package pqErr

type Code = int8

const (
	UnknownError = Code(iota)
	UnhandledError
	UniqueError
	FkViolationError
)
