package query

import (
	"fmt"
	"github.com/nbs-go/nsql/query"
	"strings"
)

func newTableWriter(tableName string, as string) *tableWriter {
	return &tableWriter{
		tableName: tableName,
		as:        as,
		joints:    map[string]query.JoinWriter{},
	}
}

type tableWriter struct {
	tableName string
	as        string
	joints    map[string]query.JoinWriter
}

func (s *tableWriter) Join(j query.JoinWriter) {
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
	i := 0
	for _, jw := range s.joints {
		jointQueries[i] = jw.JoinQuery()
		i++
	}
	join := strings.Join(jointQueries, " ")

	return q + " " + join
}
