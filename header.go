package dpmafirmware

import "time"

// Header holds metadata for a file entry in a firmware package.
type Header struct {
	Name    string
	Size    int64
	ModTime time.Time
}
