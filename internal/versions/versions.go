package versions

import "fmt"

const (
	versionFmt = "%s (%s %s)"
)

var (
	// Version is the version of the builder.
	Version = "dev"
	// Commit is the commit hash of the builder.
	Commit = "none"
	// Date is the date of the builder.
	Date = "unknown"
)

// VersionString is the version string of the builder.
var VersionString = fmt.Sprintf(versionFmt, Version, Commit, Date)
