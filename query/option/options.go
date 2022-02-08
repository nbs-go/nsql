package opt

import (
	"github.com/nbs-go/nsql/query"
	"github.com/nbs-go/nsql/query/op"
	"github.com/nbs-go/nsql/schema"
)

// Built in option keys

const (
	SchemaKey         = "schema"
	AsKey             = "as"
	ColumnsKey        = "columns"
	SortDirectionKey  = "sortDirection"
	CountKey          = "count"
	JoinMethodKey     = "joinMethod"
	VariableKey       = "variable"
	VariableFormatKey = "varFmt"
	ColumnFormatKey   = "columnFmt"
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

func (o *Options) GetVariableFormat() (query.VariableFormat, bool) {
	v, ok := o.KV[VariableFormatKey]
	if !ok {
		return 0, false
	}

	vf, fOk := v.(query.VariableFormat)
	return vf, fOk
}

func (o *Options) GetColumnFormat() (query.ColumnFormat, bool) {
	v, ok := o.KV[ColumnFormatKey]
	if !ok {
		return 0, false
	}

	f, fOk := v.(query.ColumnFormat)
	return f, fOk
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

func VariableFormat(vf query.VariableFormat) SetOptionFn {
	return func(o *Options) {
		o.KV[VariableFormatKey] = vf
	}
}

func ColumnFormat(f query.ColumnFormat) SetOptionFn {
	return func(o *Options) {
		o.KV[ColumnFormatKey] = f
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
