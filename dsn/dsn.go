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
	// TODO: Add support ssl mode
	q.Set("sslmode", fmt.Sprintf("%t", o.SslMode))
	q.Set("parseTime", fmt.Sprintf("%t", o.ParseTime))
	var scheme string

	switch driver {
	case DriverPostgres:
		scheme = "postgres"
		if o.SearchPath != "" {
			q.Set("search_path", o.SearchPath)
		}
	case DriverMysql:
		scheme = "mysql"
	default:
		return "", fmt.Errorf("nsql: Unsupported driver %s", driver)
	}

	u, err := url.Parse(fmt.Sprintf("%s://%s:%s@%s:%d/%s", scheme, username, password, host, port, dbName))
	if err != nil {
		return "", err
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}
