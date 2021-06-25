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

	cmd.Flags().BoolP("write-to-source", "w", false, "write result to (source) file instead  of stdout")

	return cmd
}

func runFormat(cmd *cobra.Command, args []string) error {
	writeToSource, _ := cmd.Flags().GetBool("write-to-source")

	format := format.New(cmd.OutOrStdout(), writeToSource)
	return format.Run(args)
}
