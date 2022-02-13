package query

import (
	"fmt"
	"github.com/nbs-go/nsql"
	"strings"
)

func newTableWriter(tableName string, as string) *tableWriter {
	return &tableWriter{
		tableName: tableName,
		as:        as,
		joints:    map[string]nsql.JoinWriter{},
	}
}

type tableWriter struct {
	tableName string
	as        string
	joints    map[string]nsql.JoinWriter
}

func (s *tableWriter) Join(j nsql.JoinWriter) {
	// Set join index
	idx := len(s.joints)
	j.SetIndex(idx)

	// Add joints
	s.joints[j.GetTableName()] = j
}

func (s *tableWriter) GetTableName() string {
	return s.tableName
}

func (s *tableWriter) FromQuery() string {
	q := fmt.Sprintf(`"%s"`, s.tableName)

	if s.as != "" {
		q += fmt.Sprintf(` AS "%s"`, s.as)
	}

	jointCount := len(s.joints)
	if jointCount == 0 {
		return q
	}

	jointQueries := make([]string, jointCount)
	for _, jw := range s.joints {
		jointQueries[jw.GetIndex()] = jw.JoinQuery()
	}
	join := strings.Join(jointQueries, " ")

	return q + " " + join
}
