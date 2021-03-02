package app

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	exitCodeNormal = 0
	exitCodeError  = 1
)

func Execute(version string, stdout, stderr io.Writer) int {
	configDir, err := os.UserConfigDir()
	if err == nil {
		configDir = path.Join(configDir, "mustache")
	}

	// Initialize config
	viper.SetConfigName("mustache")
	viper.AddConfigPath(".")
	if configDir != "" {
		viper.AddConfigPath(configDir)
	}

	// Setup default values

	// Read config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintf(stderr, "error processing config file: %v", err)
			return exitCodeError
		}
	}

	rootCmd := makeRootCmd(version)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.SilenceUsage = true

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(stderr, "error: %v", err)
		return exitCodeError
	}

	return exitCodeNormal
}

func makeRootCmd(version string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "mustache",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
		SilenceErrors: true,
		Version:       version,
	}

	rootCmd.InitDefaultVersionFlag()

	rootCmd.AddCommand(makeCheckCmd())

	return rootCmd
}
