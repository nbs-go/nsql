package opt

import "github.com/nbs-go/nsql/schema"

// Built in option keys

const (
	SchemaKey = "schema"
	AliasKey  = "as"
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

func (o *Options) GetSchema() *schema.Schema {
	v, ok := o.KV[SchemaKey]
	if !ok {
		return nil
	}

	return v.(*schema.Schema)
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
		o.KV[AliasKey] = as
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
