package dpmafirmware

import (
	"context"
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
	req, err := http.NewRequest("GET", r.URL(o).String(), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to prepare HTTP request: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	req = req.WithContext(ctx)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		cancel()
		return nil, err
	}
	if res.StatusCode != 200 {
		cancel()
		res.Body.Close()
		return nil, fmt.Errorf("expected http status code 200, got %d", res.StatusCode)
	}

	reader, err = NewReader(res.Body)
	if err != nil {
		cancel()
		res.Body.Close()
		return
	}

	reader.stream = res.Body // Make sure reader closes the body when finished
	reader.cancel = cancel   // Allow the reader to cancel the request
	return
}

func (r *Release) marshalRawJSON(raw map[string]interface{}) {
	raw["date"] = r.Date
	raw["md5sum"] = r.MD5Sum
	raw["models"] = r.Models
	raw["version"] = r.Version
}
