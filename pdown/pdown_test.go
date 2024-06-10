package pdown

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

func TestNewDownloadingTask(t *testing.T) {
	const contentLength = 1024
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodHead {
			w.Header().Set("Content-Length", strconv.Itoa(contentLength))
		} else if r.Method == http.MethodGet {
			w.Header().Set("Content-Length", strconv.Itoa(contentLength))
			_, err := w.Write(make([]byte, contentLength))
			assert.Nil(t, err, "error while writing response body")
		}
	}))
	defer server.Close()

	t.Run("default task", func(t *testing.T) {
		const numWorker = 4
		task, err := NewDownloadingTask(server.URL, numWorker)
		assert.Nil(t, err, "failed to new downloading task of %d workers", numWorker)
		defer os.Remove(task.TempFileName)

		assert.Equal(t, int64(contentLength), task.Length, "task content length not match")
		assert.Equal(t, numWorker, task.NumWorker, "task worker number not match")
	})

	t.Run("task with specified filename", func(t *testing.T) {
		const targetFileName = "file"
		task, err := NewDownloadingTask(server.URL, 4, WithName(targetFileName))
		assert.Nil(t, err, "failed to new downloading task of specified filename %s", targetFileName)
		defer os.Remove(task.TempFileName)

		assert.Equal(t, targetFileName, task.FileName)
	})
}

func TestRangeFile(t *testing.T) {
	const contentLength = 1024
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rangeHeader := r.Header.Get("Range"); rangeHeader == "" {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Length", strconv.Itoa(contentLength))
			_, err := w.Write(make([]byte, contentLength))
			assert.Nil(t, err, "error while writing response body")
		} else {
			var start, end int
			_, err := fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end)
			assert.Nil(t, err, "failed to parse range")

			w.WriteHeader(http.StatusPartialContent)
			w.Header().Set("Content-Length", strconv.Itoa(end-start+1))
			_, err = w.Write(make([]byte, end-start+1))
			assert.Nil(t, err, "error while writing response body")
		}
	}))
	defer server.Close()

	t.Run("test file completeness", func(t *testing.T) {
		task, err := NewDownloadingTask(server.URL, 4, WithName("filename"))
		assert.Nil(t, err, "failed to new downloading task: %s", err)

		err = task.Start()
		assert.Nil(t, err, "failed to start task: %s", err)
		assert.Equal(t, true, task.IsComplete())

		stat, err := os.Stat(task.FileName)
		assert.Nil(t, err, "file not existing")
		assert.Equal(t, int64(contentLength), stat.Size())
		_ = os.Remove(task.FileName)
	})

	// TODO: Check file correctness
}
