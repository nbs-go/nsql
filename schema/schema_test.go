package schema

import (
	"testing"
	"time"
)

type Person struct {
	CreatedAt time.Time `db:"createdAt"`
	UpdatedAt time.Time `db:"updatedAt"`
	Id        int64     `db:"id"`
	FullName  string    `db:"fullName"`
	BirthDate string    `db:"birthDate"`
	Age       int       `db:"-"`
	NickName  string
	lastName  string
}

func TestPointer(t *testing.T) {
	sPtr := New(FromModelRef(new(Person)))
	sStruct := New(FromModelRef(Person{}))

	if sPtr.TableName != sStruct.TableName {
		t.Errorf("got different table name. sPtr = %s, sStruct = %s", sPtr.TableName, sStruct.TableName)
	}

	if sPtr.CountColumns() != sStruct.CountColumns() {
		t.Errorf("got different columns length. sPtr = %d, sStruct = %d", sPtr.CountColumns(), sStruct.CountColumns())
	}
}

func TestEmbeddedFields(t *testing.T) {
	type BaseModel struct {
		CreatedAt time.Time `db:"createdAt"`
		UpdatedAt time.Time `db:"updatedAt"`
		Id        int64     `db:"id"`
	}

	type PersonEmb struct {
		BaseModel
		FullName  string `db:"fullName"`
		BirthDate string `db:"birthDate"`
		NickName  string
	}

	type PersonEmbPtr struct {
		*BaseModel
		FullName  string `db:"fullName"`
		BirthDate string `db:"birthDate"`
		NickName  string
	}

	sEmb := New(FromModelRef(PersonEmb{}), TableName("Person"))
	sEmbPtr := New(FromModelRef(PersonEmbPtr{}), TableName("Person"))
	sNoEmb := New(FromModelRef(Person{}), TableName("Person"))

	if sEmb.TableName != sNoEmb.TableName {
		t.Errorf("got different table name. sEmb = %s, sNoEmb = %s", sEmb.TableName, sNoEmb.TableName)
	}

	if sEmbPtr.TableName != sNoEmb.TableName {
		t.Errorf("got different table name. sEmbPtr = %s, sNoEmb = %s", sEmbPtr.TableName, sNoEmb.TableName)
	}

	if sEmb.CountColumns() != sNoEmb.CountColumns() {
		t.Errorf("got different columns length. sEmb = %d, sNoEmb = %d", sEmb.CountColumns(), sNoEmb.CountColumns())
	}

	if sEmbPtr.CountColumns() != sNoEmb.CountColumns() {
		t.Errorf("got different columns length. sEmbPtr = %d, sNoEmb = %d", sEmbPtr.CountColumns(), sNoEmb.CountColumns())
	}
}

func TestManual(t *testing.T) {
	sModelRef := New(FromModelRef(Person{}), AutoIncrement(false))
	sManual := New(TableName("Person"), Columns("createdAt", "updatedAt", "id", "fullName", "birthDate", "NickName"), AutoIncrement(false))

	if sModelRef.TableName != sManual.TableName {
		t.Errorf("got different table name. sModelRef = %s, sManual = %s", sModelRef.TableName, sManual.TableName)
	}

	if sModelRef.CountColumns() != sManual.CountColumns() {
		t.Errorf("got different columns length. sModelRef = %d, sManual = %d", sModelRef.CountColumns(), sManual.CountColumns())
	}

	if sModelRef.AutoIncrement != sManual.AutoIncrement {
		t.Errorf("got different auto increment value. sModelRef = %t, sManual = %t", sModelRef.AutoIncrement, sManual.AutoIncrement)
	}

	// Check columns
	cols := sManual.GetColumns()

	for _, c := range cols {
		if !sModelRef.IsColumnExist(c) {
			t.Errorf("column not found in model reference. Column: %s", c)
		}
	}
}

func TestCustomPK(t *testing.T) {
	type Log struct {
		LogId   string `db:"logId"`
		Message string `db:"message"`
	}

	s := New(FromModelRef(Log{}), PrimaryKey("logId"))

	if s.PrimaryKey != "logId" {
		t.Errorf("got different custom primary key. PrimaryKey = %s", s.PrimaryKey)
	}
}
