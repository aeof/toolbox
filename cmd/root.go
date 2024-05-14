package cmd

import "github.com/spf13/cobra"

var (
	rootCmd = &cobra.Command{
		Use:   "toolbox",
		Short: "Toolbox is a simple tool box for easy jobs",
		Long:  "Toolbox is an easy tool sets for easy jobs like password generator, parallel downloading",
		Args:  cobra.MinimumNArgs(1),
		Run:   func(cmd *cobra.Command, args []string) {},
	}
)

func Execute() error {
	return rootCmd.Execute()
}
