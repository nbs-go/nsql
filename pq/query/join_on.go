package query

import (
	opt "github.com/nbs-go/nsql/query/option"
)

func On(col string, args ...interface{}) opt.SetOptionFn {
	return func(o *opt.Options) {
		// Evaluate options
		opts := opt.EvaluateOptions(args)
		schema := opts.GetSchema()

		// Get table name
		var tableName string
		if schema == nil {
			tableName = joinTableFlag
		}

		// Set variable
		o.KV[opt.VariableKey] = &columnWriter{
			name:      col,
			tableName: tableName,
		}
	}
}
