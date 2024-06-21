package pdown

import (
	"errors"
	"fmt"
	"github.com/aeof/toolbox/utils"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/schollz/progressbar/v3"
)

/*
	1. Use HEAD method to get the file size(HTTP response header Content-Length)
	2. Create a temporary file to store the downloaded content
	3. Calculate the range of each slice for each worker
	4. Send a GET request with Range header to pdown the slice
	5. Write the downloaded content to the temporary file

	Note: To avoid concurrent write to the same file, we open the file in each worker.

	TODO: Add handling logic for the following cases:
	- The server sends a 429 Too Many Requests status code, slow down the speed
	- Custom rate limiter to control the pdown speed
	- Add a timeout for the HTTP request
	- Add a retry mechanism for the failed HTTP request
	- Add a checksum to verify the downloaded content
*/

// HTTP status code that may occur:
// 200 OK: we can only pdown the complete file
// 206 Partial Content: we can pdown the file slice separately
// 429 Too Many Requests: we need to slow down the pdown speed

const (
	ProgressBarHint = "Downloading task: "
)

var (
	// ErrDownloadingNotCompleted indicates the task finished incompletely
	ErrDownloadingNotCompleted = errors.New("downloading not completed")
)

// DownloadingTask represents an HTTP file downloading task
type DownloadingTask struct {
	// the URL of the file to pdown
	Url string

	// the file name to save the downloaded content
	FileName string

	// the temporary file name to store the unfinished downloaded content
	TempFileName string

	// the length of the total file in bytes, -1 means unknown file size
	Length int64

	// concurrent workers to pdown the file
	NumWorker int

	// the downloading status of each worker
	CompleteStatus []bool
}

// NewDownloadingTaskOption Options to create specified task
// As Go doesn't allow overriding, we can use options pattern to achieve similar effects
type NewDownloadingTaskOption func(task *DownloadingTask)

func WithName(name string) NewDownloadingTaskOption {
	return func(task *DownloadingTask) {
		if name != "" {
			task.FileName = name
			utils.LogVerbose("Set file name to:", name)
		}
	}
}

func NewDownloadingTask(fileURL string, numWorker int, options ...NewDownloadingTaskOption) (*DownloadingTask, error) {
	utils.LogVerbose("File URL:", fileURL)
	utils.LogVerbose("Number of workers:", numWorker)

	// fetch the content length of the file
	resp, err := http.Head(fileURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	contentLength, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		numWorker = 1
		contentLength = -1
		utils.LogVerbose("Unknown task file size")
	} else {
		utils.LogVerbose("Task file size:", contentLength)
	}

	// create a tmp file to store the downloaded content
	file, err := os.CreateTemp(".", ".multi-downloader-*")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	utils.LogVerbose("Temp file created:", file.Name())

	// create a downloading task and start downloading
	task := &DownloadingTask{
		Url:            fileURL,
		Length:         contentLength,
		NumWorker:      numWorker,
		CompleteStatus: make([]bool, numWorker),
		TempFileName:   file.Name(),
		FileName:       parseFileName(fileURL),
	}
	for _, option := range options {
		option(task)
	}
	return task, nil
}

// IsComplete checks if all the workers have finished downloading
func (t *DownloadingTask) IsComplete() bool {
	for _, status := range t.CompleteStatus {
		if !status {
			return false
		}
	}
	return true
}

// Start starts the downloading task by assigning the downloading task to all the workers.
func (t *DownloadingTask) Start() error {
	// use a channel to receive the downloading status for each slice
	statusStream := make(chan TaskSliceStatus)
	pb := progressbar.DefaultBytes(t.Length, ProgressBarHint) // use progress bar to show the downloading progress
	for i := 0; i < t.NumWorker; i++ {
		t.startSlice(i, statusStream, pb)
	}

	// collect the downloading status from all the workers
	for i := 0; i < t.NumWorker; i++ {
		status := <-statusStream
		if status.Err == nil {
			t.CompleteStatus[status.WorkerId] = true
		} else {
			fmt.Printf("failed to pdown slice %d: %s\n", status.WorkerId, status.Err.Error())
		}
	}

	if t.IsComplete() {
		return os.Rename(t.TempFileName, t.FileName)
	}
	return ErrDownloadingNotCompleted
}

// Resume resumes downloading unfinished slices of the task
func (t *DownloadingTask) Resume() error {
	statusStream := make(chan TaskSliceStatus)
	pb := progressbar.DefaultBytes(-1, ProgressBarHint) // use progress bar to show the downloading progress

	countMissingSlice := 0
	for i := 0; i < t.NumWorker; i++ {
		if !t.CompleteStatus[i] {
			countMissingSlice++
			t.startSlice(i, statusStream, pb)
		}
	}

	for i := 0; i < countMissingSlice; i++ {
		status := <-statusStream
		if status.Err == nil {
			t.CompleteStatus[status.WorkerId] = true
		} else {
			fmt.Printf("failed to pdown slice %d: %s\n", status.WorkerId, status.Err.Error())
		}
	}

	if t.IsComplete() {
		return nil
	}
	return ErrDownloadingNotCompleted
}

func (t *DownloadingTask) startSlice(sliceNumber int, statusStream chan<- TaskSliceStatus, pb io.Writer) {
	go func(workerId int) {
		statusStream <- TaskSliceStatus{
			WorkerId: workerId,
			Err:      downloadSlice(t, workerId, pb),
		}
	}(sliceNumber)
}

type TaskSliceStatus struct {
	WorkerId int
	Err      error
}

func downloadSlice(task *DownloadingTask, workerId int, pb io.Writer) error {
	// open the file in each worker, so we don't have to synchronize the file access
	// each worker's file descriptor has its own offset
	file, err := os.OpenFile(task.TempFileName, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// calculate the range of the slice
	sliceLen := task.Length / int64(task.NumWorker)
	start := sliceLen * int64(workerId)
	end := start + sliceLen - 1
	if workerId == task.NumWorker-1 {
		end = task.Length - 1
	}
	if _, err = file.Seek(start, io.SeekStart); err != nil {
		return err
	}

	w := io.MultiWriter(file, pb)
	req, err := http.NewRequest("GET", task.Url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36 Edg/124.0.0.0")
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent && len(task.CompleteStatus) > 1 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	_, err = io.Copy(w, resp.Body)
	return err
}

// parseFileName extracts the file name from the file URL
// by removing the query string and the path
func parseFileName(fileURL string) string {
	idx := strings.Index(fileURL, "?")
	if idx != -1 {
		fileURL = fileURL[:idx]
	}

	idx = strings.LastIndex(fileURL, "/")
	if idx != -1 {
		return fileURL[idx+1:]
	}
	return fileURL
}
