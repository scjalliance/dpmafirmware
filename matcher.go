package dpmafirmware

// Matcher is a string or pattern matcher.
type Matcher interface {
	Match(string) bool
}
