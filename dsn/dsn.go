package dsn

import (
	"fmt"
	"net/url"
)

func NormalizeDriver(d string) string {
	switch d {
	case "postgresql", "pg":
		return DriverPostgres
	}
	return d
}

func Format(driver, username, password, host string, port int, dbName string, args ...OptionSetter) (string, error) {
	o := evaluateOptions(args)

	q := make(url.Values)
	q.Set("parseTime", fmt.Sprintf("%t", o.ParseTime))
	var scheme string

	switch driver {
	case DriverPostgres:
		scheme = "postgres"
		if o.SearchPath != "" {
			q.Set("search_path", o.SearchPath)
		}
		// TODO: Add support ssl mode
		if !o.SslMode {
			q.Set("sslmode", "disable")
		}
	case DriverMysql:
		scheme = "mysql"
		q.Set("sslmode", fmt.Sprintf("%t", o.SslMode))
	default:
		return "", fmt.Errorf("nsql: Unsupported driver %s", driver)
	}

	password = url.QueryEscape(password)

	u, err := url.Parse(fmt.Sprintf("%s://%s:%s@%s:%d/%s", scheme, username, password, host, port, dbName))
	if err != nil {
		return "", err
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}
