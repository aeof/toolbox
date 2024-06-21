package cmd

import (
	"github.com/aeof/toolbox/utils"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "toolbox",
		Short: "Toolbox is a simple tool box for easy jobs",
		Long:  "Toolbox is an easy tool sets for easy jobs like password generator, parallel downloading",
		Args:  cobra.MinimumNArgs(1),
		Run:   func(cmd *cobra.Command, args []string) {},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			utils.LogVerbose("Verbose mode is enabled")
		},
	}
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&utils.Verbose, "verbose", "v", false, "enable verbose mode")
}

func Execute() error {
	return rootCmd.Execute()
}
