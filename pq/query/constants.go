package query

// Flags, start with underscore to prevent table naming collision

const (
	forceWriteFlag = "__force__" // Skip table reference checking
	fromTableFlag  = "__from__"  // Use table that is declared in from
	skipTableFlag  = "__skip__"  // Mark query part will be excluded

	AllColumns = "*"
)
