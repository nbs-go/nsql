package query

import (
	"fmt"
	"github.com/nbs-go/nsql"
	"github.com/nbs-go/nsql/op"
	"github.com/nbs-go/nsql/option"
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

func LikeFilter(col string, likeVar op.LikeVariable, args ...interface{}) nsql.FilterParser {
	// Get options
	opts := option.EvaluateOptions(args)
	s := opts.GetSchema()

	return func(qv string) (nsql.WhereWriter, []interface{}) {
		// Trim value
		switch likeVar {
		case op.LikeSubString:
			qv = fmt.Sprintf(`%%%s%%`, qv)
		case op.LikePrefix:
			qv = fmt.Sprintf(`%%%s`, qv)
		case op.LikeSuffix:
			qv = fmt.Sprintf(`%s%%`, qv)
		}

		w := Like(Column(col, option.Schema(s)), qv)

		return w, []interface{}{qv}
	}

}
