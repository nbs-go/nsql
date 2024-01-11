package nsql

import (
	"encoding/json"
	"fmt"
)

// ScanJSON is a generic scanner function that can be added to a struct that implements sql.Scanner.
//
//	func (o *Object) Scan(src interface{}) error {
//			 return nsql.ScanJSON(src, o)
//	}
func ScanJSON(src interface{}, target interface{}) error {
	// If source is nil, set target to nil
	if src == nil {
		return nil
	}

	// Assert source to byte
	var b []byte
	switch v := src.(type) {
	case []byte:
		b = v
	case json.RawMessage:
		b = v
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("nsql: type assertion to []byte failed, got type %T", v)
	}

	// Unmarshal to target
	err := json.Unmarshal(b, target)
	if err != nil {
		return err
	}
	return nil
}
