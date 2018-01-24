package dpmafirmware

import "strings"

// ModelMap is a map of models.
type ModelMap map[string]struct{}

// Contains returns true if the map contains one or more of the given models.
func (mm ModelMap) Contains(models ...string) bool {
	// Exact match
	for _, model := range models {
		if _, found := mm[strings.ToLower(model)]; found {
			return true
		}
	}

	// Wildcard match
	for model := range mm {
		if model == Wildcard {
			return true
		}
	}

	return false
}
