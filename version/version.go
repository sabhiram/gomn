package version

////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"
)

////////////////////////////////////////////////////////////////////////////////

const (
	VersionMajor = 0
	VersionMid   = 0
	VersionMinor = 1
)

////////////////////////////////////////////////////////////////////////////////

var (
	VersionString = fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMid, VersionMinor)
)

////////////////////////////////////////////////////////////////////////////////