package format

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"

	config "github.com/breml/logstash-config"
	"github.com/breml/logstash-config/ast"
)

type Format struct {
	out           io.Writer
	writeToSource bool
}

func New(out io.Writer, writeToSource bool) Format {
	return Format{
		out:           out,
		writeToSource: writeToSource,
	}
}

func (f Format) Run(args []string) error {
	for _, filename := range args {
		stat, err := os.Stat(filename)
		if err != nil {
			return errors.Errorf("%s: %v", filename, err)
		}
		if stat.IsDir() {
			continue
		}

		c, err := config.ParseFile(filename)
		if err != nil {
			if errMsg, hasErr := config.GetFarthestFailure(); hasErr {
				if !strings.Contains(err.Error(), errMsg) {
					return errors.Errorf("%s: %v\n%s", filename, err, errMsg)
				}
			}
			return errors.Errorf("%s: %v", filename, err)
		}

		if f.writeToSource {
			err := func() error {
				f, err := os.Create(filename)
				if err != nil {
					return errors.Wrap(err, "failed to open file for writting with automatically fixed ID")
				}
				defer f.Close()

				conf := c.(ast.Config)
				_, err = f.WriteString(conf.String())
				if err != nil {
					return errors.Wrap(err, "failed to write file with automatically fixed ID")
				}

				return nil
			}()
			if err != nil {
				return err
			}
			continue
		}

		fmt.Fprint(f.out, c)
	}

	return nil
}
