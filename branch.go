package dpmafirmware

const branchPrefix = "versions" // For JSON encoding and decoding

// Branch is a set of releases grouped by their branch.
type Branch struct {
	Name     string
	Releases ReleaseSet
}
