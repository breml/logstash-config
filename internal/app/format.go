package app

import (
	"github.com/spf13/cobra"

	"github.com/breml/logstash-config/internal/app/format"
)

func makeFormatCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "format [<flags>] [path ...]",
		Short:         "pretty format logstash config files",
		RunE:          runFormat,
		SilenceErrors: true,
	}

	return cmd
}

func runFormat(cmd *cobra.Command, args []string) error {
	format := format.New(cmd.OutOrStdout())
	return format.Run(args)
}
