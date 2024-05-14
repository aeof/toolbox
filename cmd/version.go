package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

const (
	version = "0.0.1"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Long:  "Print the version number like x.y.z",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("toolbox version", version)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
