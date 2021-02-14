package config

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
