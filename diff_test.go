package config_test

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func printDiff(want, got string) string {
	dmp := diffmatchpatch.New()

	fileAdmp, fileBdmp, dmpStrings := dmp.DiffLinesToChars(want, got)
	diffs := dmp.DiffMain(fileAdmp, fileBdmp, false)
	diffs = dmp.DiffCharsToLines(diffs, dmpStrings)

	var buff bytes.Buffer
	for _, diff := range diffs {
		text := diff.Text

		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			buff.WriteString("\x1b[32m")
			lines := strings.Split(strings.TrimRight(text, "\n"), "\n")
			for _, line := range lines {
				buff.WriteString("+ ")
				buff.WriteString(highlightWhitespaces(line, "\x1b[42m") + "\n")
			}
			buff.WriteString("\x1b[0m")
			if !strings.HasSuffix(text, "\n") {
				buff.WriteString("\n\\ No newline at end of file\n")
			}
		case diffmatchpatch.DiffDelete:
			buff.WriteString("\x1b[31m")
			lines := strings.Split(strings.TrimRight(text, "\n"), "\n")
			for _, line := range lines {
				buff.WriteString("- ")
				buff.WriteString(highlightWhitespaces(line, "\x1b[41m") + "\n")
			}
			buff.WriteString("\x1b[0m")
			if !strings.HasSuffix(text, "\n") {
				buff.WriteString("\n\\ No newline at end of file\n")
			}
			if !strings.HasSuffix(text, "\n") {
				buff.WriteString("\n\\ No newline at end of file\n")
			}
		case diffmatchpatch.DiffEqual:
			lines := strings.Split(strings.TrimRight(text, "\n"), "\n")
			for _, line := range lines {
				buff.WriteString("  ")
				buff.WriteString(line + "\n")
			}
		}
	}
	return buff.String()
}

var tailingWhitespace = regexp.MustCompile(`([ \t]+)$`)

func highlightWhitespaces(in string, color string) string {
	return tailingWhitespace.ReplaceAllString(in, color+"$1\x1b[49m")
}
