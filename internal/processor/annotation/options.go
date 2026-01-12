package annotation

type filters struct {
	filteredRegexps []string
}

type options struct {
	global filters
}

// OptionFunc function that allows you to change options
type OptionFunc func(*options)

// WithGlobalRegexpFilter returns OptionFunc which enables global regexp exclusion for mutators
func WithGlobalRegexpFilter(filteredRegexps ...string) OptionFunc {
	return func(o *options) {
		o.global.filteredRegexps = filteredRegexps
	}
}
