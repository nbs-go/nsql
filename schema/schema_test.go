package schema

import (
	"github.com/nbs-go/nsql/test_utils"
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
	test_utils.CompareString(t, "SAME TABLE NAME", sPtr.tableName, sStruct.TableName())

	// Test #2
	test_utils.CompareStringArray(t, "SAME COLUMNS", sPtr.Columns(), sStruct.Columns())
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
	test_utils.CompareString(t, "SAME TABLE NAME", sEmb.TableName(), sNoEmb.TableName())

	// Test #2
	test_utils.CompareString(t, "SAME TABLE NAME (POINTER)", sEmbPtr.TableName(), sNoEmb.TableName())

	// Test #3
	test_utils.CompareString(t, "SAME TABLE NAME (EMBEDDED POINTER)", sEmbPtr.TableName(), sEmb.TableName())

	// Test #4
	test_utils.CompareStringArray(t, "SAME COLUMNS", sEmb.Columns(), sNoEmb.Columns())

	// Test #5
	test_utils.CompareStringArray(t, "SAME COLUMNS (POINTER)", sEmbPtr.Columns(), sNoEmb.Columns())

	// Test #6
	test_utils.CompareStringArray(t, "SAME COLUMNS (EMBEDDED POINTER)", sEmbPtr.Columns(), sEmb.Columns())
}

func TestManual(t *testing.T) {
	// Init case
	sModelRef := New(FromModelRef(Person{}), AutoIncrement(false))
	sManual := New(TableName("Person"),
		Columns("createdAt", "updatedAt", "id", "fullName", "birthDate", "NickName"),
		AutoIncrement(false))

	// Test #1
	test_utils.CompareString(t, "SAME TABLE NAME", sManual.TableName(), sModelRef.TableName())

	// Test #2
	test_utils.CompareStringArray(t, "SAME COLUMNS", sManual.Columns(), sModelRef.Columns())

	// Test #3
	test_utils.CompareBoolean(t, "SAME AUTO INCREMENT", sManual.AutoIncrement(), sModelRef.AutoIncrement())

	// Test #4
	test_utils.CompareBoolean(t, "DID NOT HAVE COLUMN", sManual.IsColumnExist("age"), false)

	// Test #5
	test_utils.CompareInt(t, "HAVE SAME COLUMN COUNT", sManual.CountColumns(), sModelRef.CountColumns())
}

func TestCustomPK(t *testing.T) {
	type Log struct {
		LogId   string `db:"logId"`
		Message string `db:"message"`
	}

	s := New(FromModelRef(Log{}), PrimaryKey("logId"))

	// Test #1
	test_utils.CompareString(t, "CUSTOM PK", s.PrimaryKey(), "logId")
}

func TestGetColumns(t *testing.T) {
	// Init schema
	person := New(FromModelRef(new(Person)))

	// Test #1
	test_utils.CompareStringArray(t, "GET COLUMNS", person.Columns(),
		[]string{"createdAt", "updatedAt", "id", "fullName", "birthDate", "NickName"})

	// Test #2
	test_utils.CompareStringArray(t, "GET INSERT COLUMNS", person.InsertColumns(),
		[]string{"createdAt", "updatedAt", "fullName", "birthDate", "NickName"})

	// Test #3
	test_utils.CompareStringArray(t, "GET UPDATE COLUMNS", person.UpdateColumns(),
		[]string{"createdAt", "updatedAt", "fullName", "birthDate", "NickName"})

	// Test #4
	personNoAI := New(FromModelRef(new(Person)), AutoIncrement(false))
	test_utils.CompareStringArray(t, "GET INSERT COLUMNS (NO AUTO INCREMENT)", personNoAI.InsertColumns(),
		[]string{"createdAt", "updatedAt", "id", "fullName", "birthDate", "NickName"})
}

func TestPanicModelRef(t *testing.T) {
	defer test_utils.RecoverPanic(t, "INVALID MODEL REF", "modelRef must be a struct or pointer. Got string")()
	New(FromModelRef(""))
}

func TestPanicNoColumns(t *testing.T) {
	defer test_utils.RecoverPanic(t, "NO COLUMNS", "schema has no columns")()
	New(TableName("Customer"))
}

func TestPanicNoTableName(t *testing.T) {
	defer test_utils.RecoverPanic(t, "NO TABLE NAME", "schema has no table name")()
	New(Columns("id", "name"))
}

func TestPanicNoPrimaryKey(t *testing.T) {
	defer test_utils.RecoverPanic(t, "NO PRIMARY KEY", "primary key is not defined in columns")()
	New(TableName("Customer"), Columns("createdAt", "name"))
}
