package op

type LikeVariable uint8

const (
	LikeExact LikeVariable = iota
	LikeSubString
	LikePrefix
	LikeSuffix
)
