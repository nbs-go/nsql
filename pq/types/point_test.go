package pgTypes_test

import (
	pgTypes "github.com/nbs-go/nsql/pq/types"
	"github.com/nbs-go/nsql/test_utils"
	"testing"
)

func TestNewPoint(t *testing.T) {
	p := pgTypes.NewPoint(-0.8459788815501532, 119.79155128721658)
	if p.Lat() != -0.8459788815501532 {
		t.Errorf("unexpected lat set. Actual=%v", p.Lat())
	}

	if p.Lng() != 119.79155128721658 {
		t.Errorf("unexpected lng set. Actual=%v", p.Lng())
	}
}

func TestPointInvalidLat(t *testing.T) {
	defer test_utils.RecoverPanic(t, "Panic when lat > 90", "nsql/pq/types: lat out of range")()
	_ = pgTypes.NewPoint(90.1, 0)

}
func TestPointInvalidLatMinus(t *testing.T) {
	defer test_utils.RecoverPanic(t, "Panic when lat < -90", "nsql/pq/types: lat out of range")()
	_ = pgTypes.NewPoint(-90.1, 0)
}

func TestPointInvalidLng(t *testing.T) {
	defer test_utils.RecoverPanic(t, "Panic when lng > 180", "nsql/pq/types: lng out of range")()
	_ = pgTypes.NewPoint(0, 180.1)

}
func TestPointInvalidLngMinus(t *testing.T) {
	defer test_utils.RecoverPanic(t, "Panic when lng < -180", "nsql/pq/types: lng out of range")()
	_ = pgTypes.NewPoint(0, -180.1)
}

func TestPointValue(t *testing.T) {
	p := pgTypes.NewPoint(-0.8459788815501532, 119.79155128721658)

	actual, _ := p.Value()
	expected := "SRID=4326;POINT(119.79155128721658 -0.8459788815501532)"

	if actual != expected {
		t.Errorf("Unexpected value generated. Expected=%s Actual=%s", expected, actual)
	}
}

func TestPointScan(t *testing.T) {
	scanValue := []uint8("0101000020E6100000ECECBAC6A8F25D4086AAAB4D4212EBBF")

	p := pgTypes.Point{}
	err := p.Scan(scanValue)
	if err != nil {
		t.Errorf("Failed to scan value. Error=%s", err)
		return
	}

	if p.Lat() != -0.8459788815501532 {
		t.Errorf("unexpected lat parsed. Actual=%v", p.Lat())
	}

	if p.Lng() != 119.79155128721658 {
		t.Errorf("unexpected lng parsed. Actual=%v", p.Lng())
	}
}

func TestPointScan_Error_ByteOrder(t *testing.T) {
	scanValue := []uint8("4E554C4C") // NULL

	p := pgTypes.Point{}
	err := p.Scan(scanValue)
	if err == nil {
		t.Errorf("an error is supposed to be occurred")
		return
	}

	t.Logf("Error=%s", err)
}
