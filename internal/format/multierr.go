package format

import (
	"bytes"
	"fmt"
	"strings"
)

func MultiErr(es []error) string {
	if len(es) == 1 {
		return fmt.Sprintf("1 error occurred:\n%s\n", prefix(es[0].Error()))
	}

	points := make([]string, len(es))
	for i, err := range es {
		points[i] = prefix(err.Error())
	}

	return fmt.Sprintf(
		"%d errors occurred:\n%s\n",
		len(es), strings.Join(points, "\n"))
}

func prefix(in string) string {
	var s bytes.Buffer
	lines := strings.Split(strings.TrimRight(in, "\n"), "\n")
	for i, l := range lines {
		if i == 0 {
			s.WriteString(fmt.Sprintln("\t* " + l))
			continue
		}
		s.WriteString(fmt.Sprintln("\t  " + l))
	}
	return s.String()
}
