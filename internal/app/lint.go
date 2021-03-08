package app

import (
	"github.com/spf13/cobra"

	"github.com/breml/logstash-config/internal/app/lint"
)

func makeLintCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "lint [<flags>] [path ...]",
		Short:         "lint for logstash config files",
		RunE:          runLint,
		SilenceErrors: true,
	}

	return cmd
}

func runLint(_ *cobra.Command, args []string) error {
	lint := lint.New()
	return lint.Run(args)
}
