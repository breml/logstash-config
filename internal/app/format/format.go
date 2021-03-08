package format

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"

	config "github.com/breml/logstash-config"
)

type Format struct {
	out io.Writer
}

func New(out io.Writer) Format {
	return Format{
		out: out,
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

		fmt.Fprint(f.out, c)
	}

	return nil
}
