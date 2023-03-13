package schema

type options struct {
	tableName     string
	columns       []string
	primaryKey    string
	autoIncrement bool
	modelRef      interface{}
	as            string
}

var defaultOptions = &options{
	tableName:     "",
	columns:       nil,
	primaryKey:    "id",
	autoIncrement: true,
	modelRef:      nil,
	as:            "",
}

type OptionSetterFn func(*options)

func evaluateSchemaOptions(opts []OptionSetterFn) *options {
	optCopy := &options{}
	*optCopy = *defaultOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// AutoIncrement set auto increment options that will affect insert query generator
func AutoIncrement(ai bool) OptionSetterFn {
	return func(o *options) {
		o.autoIncrement = ai
	}
}

// PrimaryKey set columns that will be primary key, otherwise it will use "id" column
func PrimaryKey(pk string) OptionSetterFn {
	return func(o *options) {
		o.primaryKey = pk
	}
}

// FromModelRef reflect referenced model as Schema
func FromModelRef(m interface{}) OptionSetterFn {
	return func(o *options) {
		o.modelRef = m
	}
}

// TableName set table name of schema
func TableName(s string) OptionSetterFn {
	return func(o *options) {
		o.tableName = s
	}
}

// Columns set columns in a schema
func Columns(cols ...string) OptionSetterFn {
	return func(o *options) {
		o.columns = cols
	}
}

// As set table alias to schema, will be use as reference if set
func As(as string) OptionSetterFn {
	return func(o *options) {
		o.as = as
	}
}
