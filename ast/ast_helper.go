package ast

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

// This regular expression splits the config into indentable chunks that are:
// * empty line
// * line of comment
// * line without any quotes
// * line(s) with quoted strings, which are kept together to form a single "line of configuration"
var linesRe = regexp.MustCompile(`(\n|\s*//[^\n]*\n|[^'"\n]*\n|([^"'\n]*("(\\"|[^"])*"|'(\\'|[^'])*')[^"'\n]*)*\n)`)

func prefix(in string, emptyNewline bool) string {
	if len(in) == 0 {
		if emptyNewline {
			return "\n"
		}
		return ""
	}

	var s bytes.Buffer
	s.WriteString("\n")
	lines := linesRe.FindAllString(strings.TrimRight(in, "\n")+"\n", -1)
	for _, l := range lines {
		if len(strings.TrimLeft(l, " \n")) == 0 {
			s.WriteString("\n")
			continue
		}
		s.WriteString(fmt.Sprint("  " + l))
	}
	return s.String()
}
