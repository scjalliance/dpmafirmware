package dpmafirmware

import (
	"net/url"
	"strings"
)

// Release is a dpma firmware release.
type Release struct {
	Date    string   `json:"date"` // "YYYY-MM-DD"
	MD5Sum  string   `json:"md5sum"`
	Models  ModelSet `json:"models"`
	Version Version  `json:"version"`
}

// Branch returns the branch the release is on. If the branch cannot be
// determined, an empty string will be returned.
func (r *Release) Branch() string {
	return r.Version.Branch()
}

// URL returns the release firmware URL for the provided origin. If the origin
// is malformed in some way nil will be returned.
func (r *Release) URL(o *Origin) *url.URL {
	t, _ := url.Parse(strings.Replace(o.Tarball, VersionPlaceholder, r.Version.String(), -1))
	return o.Base.ResolveReference(t)
}

func (r *Release) marshalRawJSON(raw map[string]interface{}) {
	raw["date"] = r.Date
	raw["md5sum"] = r.MD5Sum
	raw["models"] = r.Models
	raw["version"] = r.Version
}
