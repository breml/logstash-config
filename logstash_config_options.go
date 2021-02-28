package config

const ignoreComment = "__ignoreComments"

// The IgnoreComments option controls if the parse ignores comments (true).
// Otherwise the comments are parsed and returned, if the configuration is
// returned (default).
func IgnoreComments(enabled bool) Option {
	return func(p *parser) Option {
		old, _ := p.cur.globalStore[ignoreComment].(bool)
		p.cur.globalStore[ignoreComment] = enabled
		return IgnoreComments(old)
	}
}

const exceptionalCommentWarning = "__exceptionalCommentWarning"

// The ExceptionalCommentsWarning option controls if the parser does emit warnings
// for comments in exceptional locations (true).
// Otherwise these comments are silently ignored and are lost, during the
// parsing of the configuration (default).
func ExceptionalCommentsWarning(enabled bool) Option {
	return func(p *parser) Option {
		old, _ := p.cur.globalStore[exceptionalCommentWarning].(bool)
		p.cur.globalStore[exceptionalCommentWarning] = enabled
		return ExceptionalCommentsWarning(old)
	}
}
