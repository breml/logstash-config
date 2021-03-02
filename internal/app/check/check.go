package check

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	config "github.com/breml/logstash-config"
)

type Check struct{}

func New() Check {
	return Check{}
}

func (f Check) Run(args []string) error {
	var result *multierror.Error

	for _, filename := range args {
		stat, err := os.Stat(filename)
		if err != nil {
			result = multierror.Append(result, err)
		}
		if stat.IsDir() {
			continue
		}

		_, err = config.ParseFile(filename)
		if err != nil {
			if errMsg, hasErr := config.GetFarthestFailure(); hasErr {
				if !strings.Contains(err.Error(), errMsg) {
					err = errors.Errorf("%v\n%s", err, errMsg)
				}
			}
			result = multierror.Append(result, err)
			continue
		}
	}

	if result != nil {
		result.ErrorFormat = multiErrFormat
		return result
	}

	return nil
}

func multiErrFormat(es []error) string {
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
