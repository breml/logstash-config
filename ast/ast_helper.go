package ast

import (
	"bytes"
	"fmt"
	"strings"
)

func prefix(in string, emptyNewline bool) string {
	if len(in) == 0 {
		if emptyNewline {
			return "\n"
		}
		return ""
	}

	var s bytes.Buffer
	s.WriteString("\n")
	lines := strings.Split(strings.TrimRight(in, "\n"), "\n")
	for _, l := range lines {
		if len(l) == 0 {
			s.WriteString("\n")
			continue
		}
		s.WriteString(fmt.Sprintln("  " + l))
	}
	return s.String()
}
