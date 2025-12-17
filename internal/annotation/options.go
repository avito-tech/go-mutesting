package annotation

type filters struct {
	filteredRegexps []string
}

type options struct {
	global filters
}

type OptionFunc func(*options)

func WithGlobalRegexpFilter(filteredRegexps []string) OptionFunc {
	return func(o *options) {
		o.global.filteredRegexps = filteredRegexps
	}
}
