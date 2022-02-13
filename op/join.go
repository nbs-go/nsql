package op

type JoinMethod uint8

const (
	InnerJoin JoinMethod = iota
	LeftJoin
	RightJoin
	FullJoin
)
