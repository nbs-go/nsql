package nsql

type FilterParser = func(queryValue string) (WhereWriter, []interface{})
