package dpmafirmware

import (
	"encoding/json"
	"strings"
	"unicode"
)

// Wildcard matches any model.
const Wildcard = "*"

// ModelSet represents a set of models.
type ModelSet []string

// ParseModelSet parses the given string as a model set.
func ParseModelSet(value string) ModelSet {
	// Split on whitespace and commas
	sep := func(c rune) bool {
		return unicode.IsSpace(c) || c == ','
	}
	fields := strings.FieldsFunc(value, sep)

	// Make everything uppercase
	for i := range fields {
		fields[i] = strings.ToUpper(fields[i])
	}

	return ModelSet(fields)
}

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

// Wildcard returns true if the model set includes a wildcard "*" model.
func (ms ModelSet) Wildcard() bool {
	for _, model := range ms {
		if model == Wildcard {
			return true
		}
	}
	return false
}

// Match returns true if any model within the set matches the given matcher.
func (ms ModelSet) Match(matcher Matcher) bool {
	for _, model := range ms {
		if matcher.Match(model) {
			return true
		}
	}
	return false
}

// Map returns a map of models present in the model set.
func (ms ModelSet) Map() (mm ModelMap) {
	mm = make(ModelMap, len(ms))
	for _, model := range ms {
		mm[strings.ToLower(model)] = struct{}{}
	}
	return
}

func (ms ModelSet) String() string {
	return strings.Join(ms, ",")
}
