package check

import (
	"os"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	config "github.com/breml/logstash-config"
	"github.com/breml/logstash-config/internal/format"
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
			result = multierror.Append(result, errors.Errorf("%s: %v", filename, err))
		}
		if stat.IsDir() {
			continue
		}

		_, err = config.ParseFile(filename)
		if err != nil {
			if errMsg, hasErr := config.GetFarthestFailure(); hasErr {
				if !strings.Contains(err.Error(), errMsg) {
					err = errors.Errorf("%s: %v\n%s", filename, err, errMsg)
				}
			}
			result = multierror.Append(result, errors.Errorf("%s: %v", filename, err))
			continue
		}
	}

	if result != nil {
		result.ErrorFormat = format.MultiErr
		return result
	}

	return nil
}
