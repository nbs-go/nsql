package query

import (
	"fmt"
)

type tableWriter struct {
	tableName string
	as        string
}

func (s *tableWriter) GetTableName() string {
	return s.tableName
}

func (s *tableWriter) FromQuery() string {
	q := fmt.Sprintf(`"%s"`, s.tableName)

	if s.as != "" {
		q += fmt.Sprintf(` AS "%s"`, s.as)
	}

	return q
}
