package lint

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	config "github.com/breml/logstash-config"
	"github.com/breml/logstash-config/ast"
	"github.com/breml/logstash-config/ast/astutil"
	"github.com/breml/logstash-config/internal/format"
)

type Lint struct{}

func New() Lint {
	return Lint{}
}

func (l Lint) Run(args []string) error {
	var result *multierror.Error

	for _, filename := range args {
		stat, err := os.Stat(filename)
		if err != nil {
			result = multierror.Append(result, errors.Errorf("%s: %v", filename, err))
		}
		if stat.IsDir() {
			continue
		}

		c, err := config.ParseFile(filename)
		if err != nil {
			if errMsg, hasErr := config.GetFarthestFailure(); hasErr {
				if !strings.Contains(err.Error(), errMsg) {
					err = errors.Errorf("%s: %v\n%s", filename, err, errMsg)
				}
			}
			result = multierror.Append(result, errors.Errorf("%s: %v", filename, err))
			continue
		}
		conf := c.(ast.Config)

		v := validator{}

		for i := range conf.Input {
			astutil.ApplyPlugins(conf.Input[i].BranchOrPlugins, v.walk)
		}
		for i := range conf.Filter {
			astutil.ApplyPlugins(conf.Filter[i].BranchOrPlugins, v.walk)
		}
		for i := range conf.Output {
			astutil.ApplyPlugins(conf.Output[i].BranchOrPlugins, v.walk)
		}

		if len(v.noIDs) > 0 {
			errMsg := strings.Builder{}
			errMsg.WriteString(fmt.Sprintf("%s: no IDs found for:\n", filename))
			for _, block := range v.noIDs {
				errMsg.WriteString(block + "\n")
			}
			result = multierror.Append(result, errors.New(errMsg.String()))
		}
	}

	if result != nil {
		result.ErrorFormat = format.MultiErr
		return result
	}

	return nil
}

type validator struct {
	noIDs []string
}

func (v *validator) walk(c *astutil.Cursor) {
	_, err := c.Plugin().ID()
	if err != nil {
		v.noIDs = append(v.noIDs, c.Plugin().String())
	}
}
