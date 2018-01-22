package dpmafirmware

import (
	"encoding/json"
	"net/url"
)

// VersionPlaceholder is the placeholder text in a tarball path that should be
// replaced with the version number of a particular release.
const VersionPlaceholder = "{version}"

// Origin describes an origin from which dpma firmware can be downloaded.
type Origin struct {
	Base    *url.URL `json:"path"`    // Base URL
	Tarball string   `json:"tarball"` // Tarball path relative to base
}

// String returns a string representation of the origin.
func (o *Origin) String() string {
	return o.Base.String() + o.Tarball
}

// UnmarshalJSON unmarshals the origin from JSON-encoded data.
func (o *Origin) UnmarshalJSON(b []byte) error {
	raw := struct {
		Base    string `json:"path"`
		Tarball string `json:"tarball"`
	}{}
	json.Unmarshal(b, &raw)

	var err error

	// Ensure that the base URL can be parsed
	o.Base, err = url.Parse(raw.Base)
	if err != nil {
		return err
	}

	// Also verify that the tarball can be parsed, even though we store it as a
	// string
	_, err = url.Parse(raw.Tarball)
	if err != nil {
		return err
	}
	o.Tarball = raw.Tarball

	// Enforce HTTPS (does this belong here?)
	if o.Base.Scheme == "http" {
		o.Base.Scheme = "https"
	}

	return nil
}

func (o *Origin) marshalRawJSON(raw map[string]interface{}) {
	raw["path"] = o.Base.String()
	raw["tarball"] = o.Tarball
}
