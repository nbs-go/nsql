package op

type Operator uint8

const (
	// --------------------
	// Comparison Operators
	// --------------------

	Equal Operator = iota
	NotEqual
	GreaterThan
	GreaterThanEqual
	LessThan
	LessThanEqual
	Like
	NotLike
	ILike
	NotILike
	Between
	NotBetween
	In
	NotIn

	// -----------------
	// Logical Operators
	// -----------------

	And
	Or
)
