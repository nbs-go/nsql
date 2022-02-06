package opt

import (
	"github.com/nbs-go/nsql/query"
	"github.com/nbs-go/nsql/query/op"
	"github.com/nbs-go/nsql/schema"
)

// Built in option keys

const (
	SchemaKey        = "schema"
	AsKey            = "as"
	ColumnsKey       = "columns"
	SortDirectionKey = "sortDirection"
	CountKey         = "count"
	JoinMethodKey    = "joinMethod"
	VariableKey      = "variable"
)

type Options struct {
	KV map[string]interface{}
}

func (o *Options) GetString(k string) (string, bool) {
	v, ok := o.KV[k]
	if !ok {
		return "", false
	}

	s, ok := v.(string)
	return s, ok
}

func (o *Options) GetStringArray(k string) []string {
	v, ok := o.KV[k]
	if !ok {
		return nil
	}
	return v.([]string)
}

func (o *Options) GetSchema() *schema.Schema {
	v, ok := o.KV[SchemaKey]
	if !ok {
		return nil
	}

	return v.(*schema.Schema)
}

func (o *Options) GetSortDirection() op.SortDirection {
	v, ok := o.KV[SortDirectionKey]
	if !ok {
		return op.Ascending
	}
	return v.(op.SortDirection)
}

func (o *Options) GetJoinMethod() op.JoinMethod {
	v, ok := o.KV[JoinMethodKey]
	if !ok {
		return op.InnerJoin
	}
	return v.(op.JoinMethod)
}

func (o *Options) GetVariable(key string) query.VariableWriter {
	v, ok := o.KV[key]
	if !ok {
		return nil
	}
	return v.(query.VariableWriter)
}

type SetOptionFn = func(*Options)

// Constructors

func NewOptions() *Options {
	return &Options{KV: make(map[string]interface{})}
}

// Option setters

func Schema(s *schema.Schema) SetOptionFn {
	return func(o *Options) {
		o.KV[SchemaKey] = s
	}
}

func As(as string) SetOptionFn {
	return func(o *Options) {
		o.KV[AsKey] = as
	}
}

func Count(col string, args ...interface{}) SetOptionFn {
	return func(o *Options) {
		// Evaluate arguments
		optCopy := NewOptions()
		for _, v := range args {
			switch cv := v.(type) {
			case SetOptionFn:
				cv(optCopy)
			}
		}

		// Copy value to kv
		for k, v := range optCopy.KV {
			o.KV[k] = v
		}

		// Set count column
		o.KV[CountKey] = col
	}
}

func Columns(args ...interface{}) SetOptionFn {
	return func(o *Options) {
		// Init columns containers
		var cols []string

		// Evaluate arguments
		optCopy := NewOptions()
		for _, v := range args {
			switch cv := v.(type) {
			case SetOptionFn:
				cv(optCopy)
			case string:
				cols = append(cols, cv)
			}
		}

		// If no columns, then skip
		if len(cols) == 0 {
			return
		}

		// Copy value to kv
		for k, v := range optCopy.KV {
			o.KV[k] = v
		}

		// Set columns value
		o.KV[ColumnsKey] = cols
	}
}

func SortDirection(direction op.SortDirection) SetOptionFn {
	return func(o *Options) {
		o.KV[SortDirectionKey] = direction
	}
}

func JoinMethod(m op.JoinMethod) SetOptionFn {
	return func(o *Options) {
		o.KV[JoinMethodKey] = m
	}
}

// Evaluator

func EvaluateOptions(args []interface{}) *Options {
	optCopy := NewOptions()
	for _, v := range args {
		fn, ok := v.(SetOptionFn)
		if !ok {
			// Skipping
			continue
		}
		fn(optCopy)
	}
	return optCopy
}
