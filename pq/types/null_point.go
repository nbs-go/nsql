package pgTypes

import (
	"database/sql/driver"
	"encoding/json"
)

// NullPoint represents an x,y coordinate in EPSG:4326 for PostGIS that supports null value
type NullPoint struct {
	point Point
	Valid bool
}

func NewNullPoint(lat, lng float64) NullPoint {
	if lat < -90 || lat > 90 || lng < -180 || lng > 180 {
		return NullPoint{Point{0, 0}, false}
	}
	return NullPoint{
		point: Point{lng, lat},
		Valid: true,
	}
}

func (p NullPoint) String() string {
	if !p.Valid {
		return "NULL"
	}
	return p.point.String()
}

// Scan implements the sql.Scanner interface.
func (p *NullPoint) Scan(val interface{}) error {
	if val == nil {
		return nil
	}

	err := p.point.Scan(val)
	if err != nil {
		return err
	}

	p.Valid = true

	return nil
}

// Value implements driver.Valuer interface
func (p NullPoint) Value() (driver.Value, error) {
	if !p.Valid {
		return nil, nil
	}
	return p.String(), nil
}

func (p NullPoint) Lat() float64 {
	return p.point.Lat()
}

func (p NullPoint) Lng() float64 {
	return p.point.Lng()
}

func (p NullPoint) Point() Point {
	return p.point
}

func (p NullPoint) MarshalJSON() ([]byte, error) {
	if !p.Valid {
		return nil, nil
	}

	return json.Marshal(p.point)
}

func (p *NullPoint) UnmarshalJSON(bytes []byte) error {
	if string(bytes) == "null" {
		return nil
	}

	var po Point
	err := json.Unmarshal(bytes, &po)
	if err != nil {
		return nil
	}

	p.Valid = true
	p.point = po
	return nil
}
