package app

import (
	"github.com/spf13/cobra"

	"github.com/breml/logstash-config/internal/app/check"
)

func makeCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "check [<flags>] [path ...]",
		Short:         "syntax check for logstash config files",
		RunE:          runCheck,
		SilenceErrors: true,
	}

	return cmd
}

func runCheck(cmd *cobra.Command, args []string) error {
	check := check.New()
	return check.Run(args)
}
