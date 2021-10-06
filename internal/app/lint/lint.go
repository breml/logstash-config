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

type Lint struct {
	autoFixID bool
}

func New(autoFixID bool) Lint {
	return Lint{
		autoFixID: autoFixID,
	}
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

		c, err := config.ParseFile(filename, config.ExceptionalCommentsWarning(true))
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
		for _, warning := range conf.Warnings {
			result = multierror.Append(result, errors.New(warning))
		}

		v := validator{
			autoFixID: l.autoFixID,
			allIDs:    make(map[string]bool),
		}

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
		if len(v.duplicateIDs) > 0 {
			errMsg := strings.Builder{}
			errMsg.WriteString(fmt.Sprintf("%s: duplicate IDs found in:\n", filename))
			for _, block := range v.duplicateIDs {
				errMsg.WriteString(block + "\n")
			}
			result = multierror.Append(result, errors.New(errMsg.String()))
		}

		if l.autoFixID && v.changed {
			func() {
				f, err := os.Create(filename)
				if err != nil {
					result = multierror.Append(result, errors.Wrap(err, "failed to open file for writting with automatically fixed ID"))
					return
				}
				defer f.Close()

				_, err = f.WriteString(conf.String())
				if err != nil {
					result = multierror.Append(result, errors.Wrap(err, "failed to write file with automatically fixed ID"))
					return
				}
			}()
		}
	}

	if result != nil {
		result.ErrorFormat = format.MultiErr
		return result
	}

	return nil
}

type validator struct {
	count        int
	noIDs        []string
	autoFixID    bool
	changed      bool
	duplicateIDs []string
	allIDs       map[string]bool
}

func (v *validator) walk(c *astutil.Cursor) {
	v.count++

	id, err := c.Plugin().ID()
	if err != nil {
		if v.autoFixID {
			v.changed = true

			plugin := c.Plugin()
			plugin.Attributes = append(plugin.Attributes, ast.NewStringAttribute("id", fmt.Sprintf("%s-%d", c.Plugin().Name(), v.count), ast.DoubleQuoted))

			c.Replace(plugin)
		} else {
			v.noIDs = append(v.noIDs, fmt.Sprintf("%s: %s", c.Plugin().Pos().String(), c.Plugin().Name()))
		}
	} else if v.allIDs[id] {
		v.duplicateIDs = append(v.duplicateIDs, fmt.Sprintf("%s: %s", c.Plugin().Pos().String(), c.Plugin().Name()))
	} else {
		v.allIDs[id] = true
	}
}
