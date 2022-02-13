package query

import (
	"github.com/nbs-go/nsql"
)

type FilterBuilder struct {
	conditions []nsql.WhereWriter
	args       []interface{}
}

func (b *FilterBuilder) Conditions() nsql.WhereWriter {
	return And(b.conditions...)
}

func (b *FilterBuilder) Args() []interface{} {
	return b.args
}

// NewFilter create a FilterBuilder that convert querystring to WHERE conditions
func NewFilter(qs map[string]string, funcMap map[string]nsql.FilterParser) *FilterBuilder {
	b := FilterBuilder{
		conditions: make([]nsql.WhereWriter, 0),
		args:       make([]interface{}, 0),
	}

	// Get value from query string
	for k, v := range qs {
		// If value is empty string, then skip
		if v == "" {
			continue
		}

		// Get mapper function
		fn, ok := funcMap[k]

		// If function mapper is not set, then skip
		if !ok {
			continue
		}

		w, args := fn(v)

		// If writer is empty or arguments 0, then skip
		if w == nil || len(args) == 0 {
			continue
		}

		// Set writers and arguments
		b.conditions = append(b.conditions, w)

		// Merge arguments
		b.args = append(b.args, args...)
	}

	return &b
}
