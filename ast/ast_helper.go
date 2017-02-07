package ast

import (
	"bytes"
	"fmt"
	"strings"
)

func prefix(in string) string {
	var s bytes.Buffer
	lines := strings.Split(strings.TrimRight(in, "\n"), "\n")
	for _, l := range lines {
		s.WriteString(fmt.Sprintln("  " + l))
	}
	return s.String()
}
