package dpmafirmware

import (
	"encoding/json"
	"strings"
)

// ModelSet represents a set of models.
type ModelSet []string

// MarshalJSON marshals the model set as a JSON-encoded comma separated string.
func (ms ModelSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(strings.Join(ms, ","))
}

// UnmarshalJSON unmarshals the model set data from a JSON-encoded comma
// separated string.
func (ms *ModelSet) UnmarshalJSON(b []byte) error {
	var val string
	if err := json.Unmarshal(b, &val); err != nil {
		return err
	}
	*ms = ModelSet(strings.Split(val, ","))
	return nil
}
