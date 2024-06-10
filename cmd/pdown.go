package cmd

import (
	"errors"
	"github.com/aeof/toolbox/pdown"
	"github.com/spf13/cobra"
	"runtime"
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
	// target file name to save to the disk
	fileName string
)

func init() {
	downloaderCmd.Flags().IntVar(&numWorkers, "thread", runtime.NumCPU(), "number of concurrent connections")
	downloaderCmd.Flags().StringVar(&fileName, "name", "", "filename to save to the disk")
	rootCmd.AddCommand(downloaderCmd)
}

func RunParallelDownload(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("missing downloading link")
	}

	urlLink := args[0]
	task, err := pdown.NewDownloadingTask(urlLink, numWorkers, pdown.WithName(fileName))
	if err != nil {
		return err
	}
	return task.Start()
}
