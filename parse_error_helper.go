package config

import "fmt"
import "bytes"

type errPos struct {
	msg string
	c   current
	pos int
}

var (
	farthestFailure []errPos
)

// GetFarthestFailure returns the farthest position where the parser had a parse error.
// The farthest position is normally close to the real source for the error.
func GetFarthestFailure() (string, bool) {
	if len(farthestFailure) > 0 {
		var bb bytes.Buffer
		bb.WriteString(fmt.Sprintf("Parsing error at pos %s and [%d] (after: '%s'):\n", farthestFailure[0].c.pos, farthestFailure[0].pos, string(farthestFailure[0].c.text)))
		for _, e := range farthestFailure {
			bb.WriteString(fmt.Sprintf("-> %s\n", e.msg))
		}
		return bb.String(), true
	}
	return "", false
}

func pos(c *current) int {
	return c.pos.offset + len(c.text)
}

// pushError is used to add potential error states to the farthestFailure slice.
// This function should be used, if there are multiple paths to be considered.
// These potential error states are a valuable source for the error message, if
// the parsing fails.
// The assumption is, that the longest successful parse tree is the most acurate.
func pushError(errorMsg string, c *current) (bool, error) {
	pos := pos(c)
	if len(farthestFailure) == 0 || pos > farthestFailure[0].pos {
		farthestFailure = []errPos{{msg: errorMsg, c: *c, pos: pos}}
	} else {
		if pos == farthestFailure[0].pos {
			for _, failure := range farthestFailure {
				if failure.msg == errorMsg {
					return false, nil
				}
			}
			farthestFailure = append(farthestFailure, errPos{msg: errorMsg, c: *c, pos: pos})
		}
	}
	return false, nil
}

// fatalError is used to abort the parsing immediately due to an unrecoverable parse error.
// In most cases this is a missing closing character of a pair, which was opened before.
// Example: a missing closing square bracket or a missing closing double quote.
func fatalError(errorMsg string, c *current) (bool, error) {
	farthestFailure = []errPos{
		{
			msg: errorMsg,
			c:   *c,
			pos: pos(c),
		},
	}
	var bb bytes.Buffer
	bb.WriteString(fmt.Sprintf("Parsing error at pos %s and [%d] (after: '%s'):\n", c.pos, pos(c), string(c.text)))
	bb.WriteString(fmt.Sprintf("-> %s\n", errorMsg))
	panic(bb.String())
}
