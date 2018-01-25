package dpmafirmware

import "time"

// Header holds metadata for a file entry in a firmware package.
type Header struct {
	Name    string // File name without directories
	Path    string // Full file path
	Size    int64
	ModTime time.Time
	Models  ModelSet
}
