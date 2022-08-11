package qs

import (
	"strconv"
)

// ParseInt parse integer value from query string.
func ParseInt(str string) (int64, bool) {
	// If empty, then ignore value by setting ok = false
	if str == "" {
		return 0, false
	}

	// Parse int
	i, err := strconv.ParseInt(str, 10, 64)
	return i, err == nil
}

// ParseFloat parse float value from query string.
func ParseFloat(str string) (float64, bool) {
	// If empty, then ignore value by setting ok = false
	if str == "" {
		return 0, false
	}

	// Parse int
	f, err := strconv.ParseFloat(str, 64)
	return f, err == nil
}
