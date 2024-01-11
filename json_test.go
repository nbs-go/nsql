package nsql_test

import (
	"encoding/json"
	"github.com/nbs-go/nsql"
	"testing"
)

func TestScanJSON_Object(t *testing.T) {
	type KeyValue struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	j := json.RawMessage(`{"key": "message", "value": "hello"}`)
	actual := KeyValue{}
	err := nsql.ScanJSON(j, &actual)
	if err != nil {
		t.Errorf("Failed to scan json object. Error=%s", err)
		return
	}

	if actual.Key != "message" || actual.Value != "hello" {
		t.Errorf("Unexpected scanned json value. JsonInput=%s Actual=%+v", j, actual)
	}
}

func TestScanJSON_ObjectArray(t *testing.T) {
	type User struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}

	j := []byte(`[{"id": 1, "name": "John Doe"}, {"id": 2, "name": "Jane Doe"}]`)
	var actual []User
	err := nsql.ScanJSON(j, &actual)
	if err != nil {
		t.Errorf("Failed to scan json array of object. Error=%s", err)
		return
	}

	if len(actual) != 2 {
		t.Errorf("Unexpected scanned json array length. JsonInput=%s Actual=%+v", j, actual)
		return
	}

	if user := actual[0]; user.Id != 1 || user.Name != "John Doe" {
		t.Errorf("Unexpected scanned json array values. JsonInput=%s Actual=%+v", j, actual)
		return
	}

	if user := actual[1]; user.Id != 2 || user.Name != "Jane Doe" {
		t.Errorf("Unexpected scanned json array values. JsonInput=%s Actual=%+v", j, actual)
		return
	}
}

func TestScanJSON_StringArray(t *testing.T) {
	j := `["alpha", "beta"]`
	var actual []string
	err := nsql.ScanJSON(j, &actual)
	if err != nil {
		t.Errorf("Failed to scan json array of string. Error=%s", err)
		return
	}

	if len(actual) != 2 {
		t.Errorf("Unexpected scanned json array length. JsonInput=%s Actual=%+v", j, actual)
		return
	}

	if actual[0] != "alpha" || actual[1] != "beta" {
		t.Errorf("Unexpected scanned json array values. JsonInput=%s Actual=%+v", j, actual)
		return
	}
}

func TestScanJSON_Int64Array(t *testing.T) {
	j := []byte(`[1, 2, 3]`)
	var actual []int64
	err := nsql.ScanJSON(j, &actual)
	if err != nil {
		t.Errorf("Failed to scan json array of int. Error=%s", err)
		return
	}

	if len(actual) != 3 {
		t.Errorf("Unexpected scanned json array length. JsonInput=%s Actual=%+v", j, actual)
		return
	}

	if actual[0] != 1 || actual[1] != 2 || actual[2] != 3 {
		t.Errorf("Unexpected scanned json array values. JsonInput=%s Actual=%+v", j, actual)
		return
	}
}

func TestScanJSON_NilSource(t *testing.T) {
	type KeyValue struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	var actual *KeyValue
	err := nsql.ScanJSON(nil, &actual)
	if err != nil {
		t.Errorf("Failed to scan nil source. Error=%s", err)
		return
	}

	var expected *KeyValue
	if actual != expected {
		t.Errorf("Unexpected nil soruce value scanned. Expected=%+v Actual=%+v", expected, actual)
		return
	}
}

func TestScanJSON_NilDest(t *testing.T) {
	type KeyValue struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	var actual *KeyValue
	rawJson := `{"key": "message", "value": "hello"}`
	err := nsql.ScanJSON(rawJson, &actual)
	if err != nil {
		t.Errorf("Failed to scan nil source. Error=%s", err)
		return
	}

	if actual.Key != "message" || actual.Value != "hello" {
		t.Errorf("Unexpected scanned value. RawJson=%s Actual=%+v", rawJson, actual)
		return
	}
}

func TestScanJSON_PrimitiveString_Error(t *testing.T) {
	j := "string"
	var dest string
	err := nsql.ScanJSON(j, &dest)
	if err == nil {
		t.Errorf("Raw Json is supposed to cause an error. Dest=%+v", dest)
		return
	}

	expected := "invalid character 's' looking for beginning of value"
	actual := err.Error()
	if actual != expected {
		t.Errorf("Unexpected error returned from scanner. Expected=%s Actual=%s", expected, actual)
		return
	}
}

func TestScanJSON_IntSource_Error(t *testing.T) {
	j := 10
	var dest int
	err := nsql.ScanJSON(j, &dest)
	if err == nil {
		t.Errorf("Int source is supposed to cause an error")
		return
	}

	expected := "nsql: type assertion to []byte failed, got type int"
	actual := err.Error()
	if actual != expected {
		t.Errorf("Unexpected error returned from scanner. Expected=%s Actual=%s", expected, actual)
		return
	}
}
