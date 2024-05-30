package pgTypes_test

import (
	"encoding/json"
	pgTypes "github.com/nbs-go/nsql/pq/types"
	"testing"
)

func TestNewNullPoint(t *testing.T) {
	p := pgTypes.NewNullPoint(-0.8459788815501532, 119.79155128721658)
	if !p.Valid {
		t.Errorf("unexpected NullPoint must be valid. Actual=%v", p.Valid)
		return
	}

	if p.Lat() != -0.8459788815501532 {
		t.Errorf("unexpected lat set. Actual=%v", p.Lat())
	}

	if p.Lng() != 119.79155128721658 {
		t.Errorf("unexpected lng set. Actual=%v", p.Lng())
	}

	po := p.Point()

	if po.Lat() != -0.8459788815501532 {
		t.Errorf("unexpected lat set. Actual=%v", po.Lat())
	}

	if po.Lng() != 119.79155128721658 {
		t.Errorf("unexpected lng set. Actual=%v", po.Lng())
	}
}

func TestNullPointInvalidLat(t *testing.T) {
	p := pgTypes.NewNullPoint(90.1, 0)
	if p.Valid {
		t.Errorf("unexpected NullPoint must not Valid")
		return
	}

}
func TestNullPointInvalidLatMinus(t *testing.T) {
	p := pgTypes.NewNullPoint(-90.1, 0)
	if p.Valid {
		t.Errorf("unexpected NullPoint must not Valid")
		return
	}
}

func TestNullPointInvalidLng(t *testing.T) {
	p := pgTypes.NewNullPoint(0, 180.1)
	if p.Valid {
		t.Errorf("unexpected NullPoint must not Valid")
		return
	}
}
func TestNullPointInvalidLngMinus(t *testing.T) {
	p := pgTypes.NewNullPoint(0, -180.1)
	if p.Valid {
		t.Errorf("unexpected NullPoint must not Valid")
		return
	}
}

func TestNullPointValue_Null(t *testing.T) {
	p := pgTypes.NullPoint{}

	actual, _ := p.Value()

	if actual != nil {
		t.Errorf("Unexpected value must be nil. Actual=%s", actual)
	}
}

func TestNullPointScan_Null(t *testing.T) {
	var scanValue interface{}

	p := pgTypes.NullPoint{}
	err := p.Scan(scanValue)
	if err != nil {
		t.Errorf("Failed to scan value. Error=%s", err)
		return
	}

	if p.Valid {
		t.Errorf("NullPoint must be invalid. Actual=%v", p.Valid)
	}
}

func TestNullPointValue(t *testing.T) {
	p := pgTypes.NewNullPoint(-0.8459788815501532, 119.79155128721658)

	actual, _ := p.Value()
	expected := "SRID=4326;POINT(119.79155128721658 -0.8459788815501532)"

	if actual != expected {
		t.Errorf("Unexpected value generated. Expected=%s Actual=%s", expected, actual)
	}
}

func TestNullPointScan(t *testing.T) {
	scanValue := []uint8("0101000020E6100000ECECBAC6A8F25D4086AAAB4D4212EBBF")

	p := pgTypes.NullPoint{}
	err := p.Scan(scanValue)
	if err != nil {
		t.Errorf("Failed to scan value. Error=%s", err)
		return
	}

	if !p.Valid {
		t.Errorf("NullPoint must be valid")
		return
	}

	if p.Lat() != -0.8459788815501532 {
		t.Errorf("unexpected lat parsed. Actual=%v", p.Lat())
	}

	if p.Lng() != 119.79155128721658 {
		t.Errorf("unexpected lng parsed. Actual=%v", p.Lng())
	}
}

func TestNullPoint_Null_String(t *testing.T) {
	p := pgTypes.NullPoint{}
	if p.String() != "NULL" {
		t.Errorf("Unexpected value given. Expected=NULL Actual=%s", p.String())
	}
}

func TestNullPointScan_Error(t *testing.T) {
	scanValue := []uint8("INVALID VALUE")

	p := pgTypes.NullPoint{}
	err := p.Scan(scanValue)
	if err == nil {
		t.Errorf("an error is supposed to be occurred")
		return
	}

	if p.Valid {
		t.Errorf("NullPoint must be invalid")
	}
}

func TestNullPoint_MarshalJson(t *testing.T) {
	p := pgTypes.NewNullPoint(-0.9491934392039606, 119.81513150854423)
	m, err := p.MarshalJSON()
	if err != nil {
		t.Errorf("failed to marshall NullPoint. Error=%s", err)
		return
	}

	if string(m) != `{"lat":-0.9491934392039606,"lng":119.81513150854423}` {
		t.Errorf("Unexpected value given. Actual=%s", string(m))
	}
}

func TestNullPoint_MarshalJsonNull(t *testing.T) {
	p := pgTypes.NullPoint{}
	m, err := p.MarshalJSON()
	if err != nil {
		t.Errorf("failed to marshall NullPoint. Error=%s", err)
		return
	}

	if m != nil {
		t.Errorf("Unexpected value given. Actual=%s", string(m))
	}
}

func TestNullPoint_UnmarshalJson(t *testing.T) {
	var p pgTypes.NullPoint
	err := json.Unmarshal([]byte(`{"lat":-0.9491934392039606,"lng":119.81513150854423}`), &p)
	if err != nil {
		t.Errorf("failed to unmarshal NullPoint. Error=%s", err)
		return
	}

	if !p.Valid {
		t.Errorf("NullPoint must be valid")
		return
	}

	if p.Lat() != -0.9491934392039606 {
		t.Errorf("unexpected lat. Actual=%v", p.Lat())
	}

	if p.Lng() != 119.81513150854423 {
		t.Errorf("unexpected lng. Actual=%v", p.Lng())
	}
}

func TestNullPoint_UnmarshalNil(t *testing.T) {
	var p pgTypes.NullPoint
	err := json.Unmarshal([]byte(`null`), &p)
	if err != nil {
		t.Errorf("failed to unmarshal NullPoint. Error=%s", err)
		return
	}

	if p.Valid {
		t.Errorf("NullPoint must be not valid")
	}
}
