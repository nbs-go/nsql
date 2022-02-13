package parse

import (
	"strconv"
	"strings"
	"time"
)

const (
	DefaultTimeLayout = "2006-01-02 15:04:05"
)

// IntArgs parse a string into an integer array arguments that will be use in query with bind vars
func IntArgs(v string) []interface{} {
	// Separate by comma
	tmp := strings.Split(v, ",")

	// Parse int to map
	iMap := make(map[int]int)
	sort := 0
	for _, iStr := range tmp {
		// Parse integer
		i, err := strconv.Atoi(iStr)
		if err != nil {
			continue
		}

		// If already exists, then skip
		if _, ok := iMap[i]; ok {
			continue
		}

		iMap[i] = sort
		sort += 1
	}

	// Convert result to map
	result := make([]interface{}, len(iMap))
	for val, i := range iMap {
		result[i] = val
	}

	return result
}

// TimeArg parse a string into a time.Time that will be use in query with bind vars
func TimeArg(v string, tLayout string) []interface{} {
	if v == strings.ToLower("now") {
		return []interface{}{time.Now()}
	}

	// If value can be parsed to number, then convert as epoch
	sec, err := strconv.ParseInt(v, 10, 64)
	if err == nil {
		return []interface{}{time.Unix(sec, 0)}
	}

	// If format is empty string, then set to default layout
	if tLayout == "" {
		tLayout = DefaultTimeLayout
	}

	// Parse time as string
	t, err := time.Parse(tLayout, v)
	if err != nil {
		return nil
	}
	return []interface{}{t}
}
