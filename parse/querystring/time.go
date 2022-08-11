package qs

import (
	"strconv"
	"strings"
	"time"
)

// ParseTime parse time value from query string. Value that can be parse are:
//  - "now"
//  - Unix Epoch
//  - Time formatted string. Format can be optionally set in args, otherwise it will parse to time.RFC3339 layout
func ParseTime(str string, args ...string) (time.Time, bool) {
	// If empty, then ignore value by setting ok = false
	if str == "" {
		return time.Time{}, false
	}

	if strings.ToLower(str) == "now" {
		return time.Now(), true
	}

	// Check if string can be parsed to epoch
	if i, err := strconv.ParseInt(str, 10, 64); err == nil {
		return time.Unix(i, 0), true
	}

	// Else, try to time.Parse
	// -- Get format from args if set
	var timeLayout string
	if len(args) > 0 {
		timeLayout = args[0]
	} else {
		timeLayout = time.RFC3339
	}

	t, err := time.Parse(timeLayout, str)
	if err != nil {
		return time.Time{}, false
	}

	return t, true
}
