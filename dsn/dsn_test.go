package dsn_test

import (
	"github.com/nbs-go/nsql/dsn"
	"testing"
)

func TestNormalizeDriver_Postgres(t *testing.T) {
	actual := dsn.NormalizeDriver("pg")
	expected := dsn.DriverPostgres
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different normalized driver value. Actual = %s", expected, actual)
	}

	actual = dsn.NormalizeDriver("postgresql")
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different normalized driver value. Actual = %s", expected, actual)
	}
}

func TestNormalizeDriver_Mysql(t *testing.T) {
	actual := dsn.NormalizeDriver("mysql")
	expected := dsn.DriverMysql
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different normalized driver value. Actual = %s", expected, actual)
	}
}

func TestFormat_Postgres_Default(t *testing.T) {
	actual, err := dsn.Format(dsn.DriverPostgres, "user", "pass", "localhost", 5432, "test_nsql")
	if err != nil {
		t.Errorf("Unable to generate DSN for driver. Error=%s", err)
		return
	}
	expected := "postgres://user:pass@localhost:5432/test_nsql?parseTime=true&sslmode=false"
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different value. Actual = %s", expected, actual)
	}
}

func TestFormat_Postgres_WithOptions(t *testing.T) {
	actual, err := dsn.Format(dsn.DriverPostgres, "user", "pass", "localhost", 5432, "test_nsql", dsn.ParseTime(false), dsn.SearchPath("app"))
	if err != nil {
		t.Errorf("Unable to generate DSN for driver. Error=%s", err)
		return
	}
	expected := "postgres://user:pass@localhost:5432/test_nsql?parseTime=false&search_path=app&sslmode=false"
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different value. Actual = %s", expected, actual)
	}
}

func TestFormat_Mysql_Default(t *testing.T) {
	actual, err := dsn.Format(dsn.DriverMysql, "user", "pass", "localhost", 3306, "test_nsql")
	if err != nil {
		t.Errorf("Unable to generate DSN for driver. Error=%s", err)
		return
	}
	expected := "mysql://user:pass@localhost:3306/test_nsql?parseTime=true&sslmode=false"
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different value. Actual = %s", expected, actual)
	}
}

func TestFormat_UnescapedSpecialCharPassword(t *testing.T) {
	actual, err := dsn.Format(dsn.DriverPostgres, "user", "pass!@#$%^&*()", "localhost", 5432, "test_nsql")
	if err != nil {
		t.Errorf("Unable to generate DSN for driver. Error=%s", err)
		return
	}
	expected := "postgres://user:pass%21%40%23$%25%5E&%2A%28%29@localhost:5432/test_nsql?parseTime=true&sslmode=false"
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different value. Actual = %s", expected, actual)
	}
}

func TestFormat_UnsupportedDriver_Error(t *testing.T) {
	_, err := dsn.Format("mssql", "user", "pass", "localhost", 1433, "test_nsql")
	if err == nil {
		t.Errorf("Unexpected result. Function must return error")
		return
	}
	actual := err.Error()
	expected := "nsql: Unsupported driver mssql"
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different value. Actual = %s", expected, actual)
	}
}

func TestFormat_UnescapedSpecialCharHostname_Error(t *testing.T) {
	_, err := dsn.Format(dsn.DriverPostgres, "user", "pass", "local%host", 5432, "test_nsql")
	if err == nil {
		t.Errorf("Unexpected result. Function must return error")
		return
	}
	actual := err.Error()
	expected := `parse "postgres://user:pass@local%host:5432/test_nsql": invalid URL escape "%ho"`
	if actual != expected {
		t.Errorf("Expected = %s\n  > got different value. Actual = %s", expected, actual)
	}
}
