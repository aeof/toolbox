package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"runtime"
	"toolbox/pdown"
)

var (
	downloaderCmd = &cobra.Command{
		Use:   "pdown",
		Short: "Pdown is a parallel downloader",
		Long:  "Pdown is a parallel downloader that use concurrent connections to download files",
		RunE:  RunParallelDownload,
	}

	// number of the concurrent downloading workers. If not specified, the count is set to the number of CPU cores
	numWorkers int
)

func init() {
	downloaderCmd.Flags().IntVar(&numWorkers, "worker", runtime.NumCPU(), "number of concurrent connections")
	rootCmd.AddCommand(downloaderCmd)
}

func RunParallelDownload(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("missing downloading link")
	}

	urlLink := args[0]
	task, err := pdown.NewDownloadingTask(urlLink, numWorkers)
	if err != nil {

		return err
	}
	return task.Start()
}
