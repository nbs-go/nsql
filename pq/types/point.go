package pgTypes

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
)

// Point represents an x,y coordinate in EPSG:4326 for PostGIS.
// Use type `geography(Point, 4326)` when declaring column
//
// Reference: https://github.com/go-pg/pg/issues/829#issuecomment-505882885
type Point [2]float64

func NewPoint(lat, lng float64) Point {
	if lat < -90 || lat > 90 {
		panic(errors.New("nsql/pq/types: lat out of range"))
	}

	if lng < -180 || lng > 180 {
		panic(errors.New("nsql/pq/types: lng out of range"))
	}

	return Point{lng, lat}
}

func (p Point) String() string {
	return fmt.Sprintf("SRID=4326;POINT(%v %v)", p[0], p[1])
}

// Scan implements the sql.Scanner interface.
func (p *Point) Scan(val interface{}) error {
	b, err := hex.DecodeString(string(val.([]uint8)))
	if err != nil {
		return err
	}
	r := bytes.NewReader(b)
	var wkbByteOrder uint8
	if err = binary.Read(r, binary.LittleEndian, &wkbByteOrder); err != nil {
		return err
	}

	var byteOrder binary.ByteOrder
	switch wkbByteOrder {
	case 0:
		byteOrder = binary.BigEndian
	case 1:
		byteOrder = binary.LittleEndian
	default:
		return fmt.Errorf("invalid byte order %d", wkbByteOrder)
	}

	var wkbGeometryType uint64
	if err = binary.Read(r, byteOrder, &wkbGeometryType); err != nil {
		return err
	}

	if err = binary.Read(r, byteOrder, p); err != nil {
		return err
	}

	return nil
}

// Value implements driver.Valuer interface
func (p Point) Value() (driver.Value, error) {
	return p.String(), nil
}

func (p Point) Lat() float64 {
	return p[1]
}

func (p Point) Lng() float64 {
	return p[0]
}
