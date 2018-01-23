package dpmafirmware

import (
	"fmt"
	"net/http"
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

// Get attempts to download the firmware from the given origin.
func (r *Release) Get(o *Origin) (reader *Reader, err error) {
	res, err := http.Get(r.URL(o).String())
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		res.Body.Close()
		return nil, fmt.Errorf("expected http status code 200, got %d", res.StatusCode)
	}

	reader, err = NewReader(res.Body)
	if err != nil {
		res.Body.Close()
		return
	}

	reader.stream = res.Body // Make sure reader closes the body when finished
	return
}

func (r *Release) marshalRawJSON(raw map[string]interface{}) {
	raw["date"] = r.Date
	raw["md5sum"] = r.MD5Sum
	raw["models"] = r.Models
	raw["version"] = r.Version
}
