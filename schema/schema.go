package schema

import (
	"fmt"
	"reflect"
)

type Schema struct {
	tableName     string
	autoIncrement bool
	primaryKey    string
	columns       map[string]int
}

func (s *Schema) TableName() string {
	return s.tableName
}

func (s *Schema) AutoIncrement() bool {
	return s.autoIncrement
}

func (s *Schema) PrimaryKey() string {
	return s.primaryKey
}

func (s *Schema) IsColumnExist(col string) bool {
	_, ok := s.columns[col]
	return ok
}

func (s *Schema) Columns() []string {
	cols := make([]string, len(s.columns))
	for k, i := range s.columns {
		cols[i] = k
	}
	return cols
}

func (s *Schema) InsertColumns() []string {
	if !s.autoIncrement {
		return s.Columns()
	}
	return s.UpdateColumns()
}

func (s *Schema) UpdateColumns() []string {
	// Get columns
	cols := s.Columns()

	// Filter out pk
	pk := s.primaryKey
	for i, c := range cols {
		if c == pk {
			cols = append(cols[:i], cols[i+1:]...)
			break
		}
	}
	return cols
}

func (s *Schema) CountColumns() int {
	return len(s.columns)
}

func New(args ...OptionSetterFn) *Schema {
	// Evaluate options
	o := evaluateSchemaOptions(args)

	// Init schema
	s := Schema{}

	// If option has referenced model, then evaluate referenced model
	if m := o.modelRef; m != nil {
		s = evaluateModelRef(m)
	} else {
		// If no columns set, then panic
		if len(o.columns) == 0 {
			panic(fmt.Errorf("schema has no columns"))
		}

		// Set columns
		s.columns = map[string]int{}
		for i, c := range o.columns {
			s.columns[c] = i
		}
	}

	// If table name is not set, then panic
	if s.tableName == "" && o.tableName == "" {
		panic(fmt.Errorf("schema has no table name"))
	}

	// Set table name or override table name if already evaluated from model reference
	if o.tableName != "" {
		s.tableName = o.tableName
	}

	// Set other options
	s.autoIncrement = o.autoIncrement

	// Check if primary key is defined in columns
	if _, ok := s.columns[o.primaryKey]; !ok {
		panic(fmt.Errorf("primary key is not defined in columns"))
	}
	s.primaryKey = o.primaryKey

	return &s
}

// evaluateModelRef returns Schema by evaluating struct
func evaluateModelRef(m interface{}) Schema {
	// Init schema
	s := Schema{}

	// Reflect type
	t := reflect.TypeOf(m)

	// Validate kind
	switch t.Kind() {
	case reflect.Ptr:
		t = t.Elem()
	case reflect.Struct:
		break
	default:
		panic(fmt.Errorf("modelRef must be a struct or pointer. Got %s", t.Name()))
	}

	s.tableName = evaluateTableName(t)
	s.columns = evaluateColumns(t)

	return s
}

// evaluateTableName returns Table Name from struct name
func evaluateTableName(t reflect.Type) string {
	return t.Name()
}

// evaluateColumns evaluate table Columns from struct fields
func evaluateColumns(t reflect.Type) map[string]int {
	// Init columns
	columns := make(map[string]int)
	pos := 0

	// Get columns from tagged fields
	for i := 0; i < t.NumField(); i++ {
		// Get field
		f := t.Field(i)

		// If field is unexported / private, then skip
		if !f.IsExported() {
			continue
		}

		// If field is an embedded field
		if f.Anonymous {
			// Get type of embedded field
			et := f.Type

			// Check kind
			switch et.Kind() {
			case reflect.Struct:
				break
			case reflect.Ptr:
				// Get struct type
				et = et.Elem()
			}

			// Get columns from embedded struct
			embColumns := evaluateColumns(et)

			// Merge map
			for k, ePos := range embColumns {
				columns[k] = pos + ePos
			}
			pos += len(embColumns)

			continue
		}

		// Get config from tag
		col := f.Tag.Get(DbTag)

		// If skipped, then move to next field
		if col == SkipField {
			continue
		}

		// If empty, use field name
		if col == "" {
			col = f.Name
		}

		// Append columns
		columns[col] = pos
		pos += 1
	}

	return columns
}
