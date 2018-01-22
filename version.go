package dpmafirmware

import "strings"

// Version is a DPMA firmware version.
type Version string

// Branch returns the name of the branch the version is on. If the branch
// cannot be determined, an empty string will be returned.
func (v Version) Branch() string {
	components := strings.SplitN(string(v), "_", 2)
	if len(components) > 0 {
		return components[0]
	}
	return ""
}

func (v Version) String() string {
	return string(v)
}

// TODO: Add glob or regex matching function
