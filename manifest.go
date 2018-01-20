package dpmafirmware

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Manifest describes a DPMA firmware origin and a set releases.
type Manifest struct {
	Origin
	Releases ReleaseSet // Most recent first
}

// MarshalJSON marshals the manifest as JSON-encoded data.
//
// The manifest will be sorted as part of the marshaling process.
func (m *Manifest) MarshalJSON() ([]byte, error) {
	// Function for mapping branch names to JSON keys
	key := func(branch string) (name string) {
		switch branch {
		case "1":
			return branchPrefix
		default:
			return fmt.Sprintf("%s%s", branchPrefix, branch)
		}
	}

	// The release set must be sorted before marshaling
	sort.Sort(m.Releases)

	// Branches have dynamic keys, so we'll use a map for marshaling
	raw := make(map[string]interface{})

	// Include the origin data
	m.Origin.marshalRawJSON(raw)

	// Include the release data for each branch
	for _, branch := range m.Releases.Branches() {
		if branch.Name == "1" {
			latest := branch.Releases.Latest()
			latest.marshalRawJSON(raw)
		}
		raw[key(branch.Name)] = branch.Releases
	}

	return json.Marshal(raw)
}

// UnmarshalJSON unmarshals the manifset data from a JSON-encoded data.
func (m *Manifest) UnmarshalJSON(b []byte) error {
	// Origin data
	if err := json.Unmarshal(b, &m.Origin); err != nil {
		return err
	}

	// Release data
	var raw map[string]*json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	var buffer ReleaseSet
	for key, value := range raw {
		if strings.HasPrefix(key, branchPrefix) {
			if err := json.Unmarshal(*value, &buffer); err != nil {
				return err
			}
			m.Releases = append(m.Releases, buffer...)
			delete(raw, key)
		}
	}
	sort.Sort(m.Releases)
	return nil
}
