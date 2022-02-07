package schema

import (
	"strings"
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

func TestStructAndPointer(t *testing.T) {
	// Init Test Cases
	sPtr := New(FromModelRef(new(Person)))
	sStruct := New(FromModelRef(Person{}))

	// Test #1
	compareString(t, "SAME TABLE NAME", sPtr.tableName, sStruct.TableName())

	// Test #2
	compareStringArray(t, "SAME COLUMNS", sPtr.Columns(), sStruct.Columns())
}

func TestEmbeddedFields(t *testing.T) {
	// Init Test Cases
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

	// Test #1
	compareString(t, "SAME TABLE NAME", sEmb.TableName(), sNoEmb.TableName())

	// Test #2
	compareString(t, "SAME TABLE NAME (POINTER)", sEmbPtr.TableName(), sNoEmb.TableName())

	// Test #3
	compareString(t, "SAME TABLE NAME (EMBEDDED POINTER)", sEmbPtr.TableName(), sEmb.TableName())

	// Test #4
	compareStringArray(t, "SAME COLUMNS", sEmb.Columns(), sNoEmb.Columns())

	// Test #5
	compareStringArray(t, "SAME COLUMNS (POINTER)", sEmbPtr.Columns(), sNoEmb.Columns())

	// Test #6
	compareStringArray(t, "SAME COLUMNS (EMBEDDED POINTER)", sEmbPtr.Columns(), sEmb.Columns())
}

func TestManual(t *testing.T) {
	// Init case
	sModelRef := New(FromModelRef(Person{}), AutoIncrement(false))
	sManual := New(TableName("Person"),
		Columns("createdAt", "updatedAt", "id", "fullName", "birthDate", "NickName"),
		AutoIncrement(false))

	// Test #1
	compareString(t, "SAME TABLE NAME", sManual.TableName(), sModelRef.TableName())

	// Test #2
	compareStringArray(t, "SAME COLUMNS", sManual.Columns(), sModelRef.Columns())

	// Test #3
	compareBoolean(t, "SAME AUTO INCREMENT", sManual.AutoIncrement(), sModelRef.AutoIncrement())
}

func TestCustomPK(t *testing.T) {
	type Log struct {
		LogId   string `db:"logId"`
		Message string `db:"message"`
	}

	s := New(FromModelRef(Log{}), PrimaryKey("logId"))

	// Test #1
	compareString(t, "CUSTOM PK", s.PrimaryKey(), "logId")
}

func TestGetColumns(t *testing.T) {
	// Init schema
	person := New(FromModelRef(new(Person)))

	// Test #1
	compareStringArray(t, "GET COLUMNS", person.Columns(),
		[]string{"createdAt", "updatedAt", "id", "fullName", "birthDate", "NickName"})

	// Test #2
	compareStringArray(t, "GET INSERT COLUMNS", person.InsertColumns(),
		[]string{"createdAt", "updatedAt", "fullName", "birthDate", "NickName"})

	// Test #3
	compareStringArray(t, "GET UPDATE COLUMNS", person.UpdateColumns(),
		[]string{"createdAt", "updatedAt", "fullName", "birthDate", "NickName"})

	// Test #4
	personNoAI := New(FromModelRef(new(Person)), AutoIncrement(false))
	compareStringArray(t, "GET INSERT COLUMNS (NO AUTO INCREMENT)", personNoAI.InsertColumns(),
		[]string{"createdAt", "updatedAt", "id", "fullName", "birthDate", "NickName"})
}

func compareStringArray(t *testing.T, expectation string, actual, expected []string) {
	s1 := strings.Join(actual, ", ")
	s2 := strings.Join(expected, ", ")

	if s1 != s2 {
		t.Errorf("%s: FAILED\n  > got different values: %s", expectation, s1)
	} else {
		t.Logf("%s: PASSED", expectation)
	}
}

func compareString(t *testing.T, expectation string, actual, expected string) {
	if actual != expected {
		t.Errorf("%s: FAILED\n  > got different values: %s", expectation, actual)
	} else {
		t.Logf("%s: PASSED", expectation)
	}
}

func compareBoolean(t *testing.T, expectation string, actual, expected bool) {
	if actual != expected {
		t.Errorf("%s: FAILED\n  > got different values: %t", expectation, actual)
	} else {
		t.Logf("%s: PASSED", expectation)
	}
}
