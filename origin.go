package dpmafirmware

// Origin describes an origin from which dpma firmware can be downloaded.
type Origin struct {
	URL     string `json:"path"` // Root URL
	Tarball string `json:"tarball"`
}

func (o *Origin) marshalRawJSON(raw map[string]interface{}) {
	raw["path"] = o.URL
	raw["tarball"] = o.Tarball
}
