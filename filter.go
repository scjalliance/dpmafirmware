package dpmafirmware

// Filter is an interface that can determine whether a release should or should
// not be included.
type Filter interface {
	Match(*Release) bool
}

type modelFilter ModelMap

// ModelFilter returns a filter for the
func ModelFilter(models ...string) Filter {
	return modelFilter(ModelSet(models).Map())
}

// Match returns true if the release includes a model within the filter.
func (f modelFilter) Match(r *Release) bool {
	return ModelMap(f).Contains(r.Models...)
}

type modelMatchFilter struct {
	matcher Matcher
}

// ModelMatchFilter returns a filter for the
func ModelMatchFilter(matcher Matcher) Filter {
	return modelMatchFilter{matcher}
}

// Match returns true if the release includes a model matched by the filter.
func (f modelMatchFilter) Match(r *Release) bool {
	return r.Models.Match(f.matcher)
}

type inverseFilter struct {
	base Filter
}

// Invert returns an inverse filter for f. Whenever f.Match would return true,
// the inverted filter returns false.
func Invert(f Filter) Filter {
	return inverseFilter{base: f}
}

// Match returns true if the release includes a model matched by the filter.
func (f inverseFilter) Match(r *Release) bool {
	return !f.base.Match(r)
}
