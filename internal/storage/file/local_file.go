package file

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
)

// Writer is a client for writing to a file.
type Writer interface {
	Write(ctx context.Context, filePath string) (io.WriteCloser, error)
}

// Reader is a client for reading from a file.
type Reader interface {
	Read(ctx context.Context, filePath string) (io.ReadCloser, error)
}

// Client is a client for reading and writing to a file.
type Client interface {
	Writer
	Reader
}

var _ Client = (*LocalFileClient)(nil)

// LocalFileClient is a client for reading and writing to a local file system.
type LocalFileClient struct {
}

// NewLocalFileClient creates a new LocalFileClient.
func NewLocalFileClient() *LocalFileClient {
	return &LocalFileClient{}
}

func (c LocalFileClient) Write(_ context.Context, filePath string) (io.WriteCloser, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	return file, nil
}

func (c LocalFileClient) Read(_ context.Context, filePath string) (io.ReadCloser, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return newBufferedFileReader(file), nil
}

type bufferedFileReader struct {
	file           *os.File
	bufferedReader io.Reader
}

func newBufferedFileReader(file *os.File) *bufferedFileReader {
	return &bufferedFileReader{
		file:           file,
		bufferedReader: bufio.NewReader(file),
	}
}

func (r bufferedFileReader) Close() error {
	return r.file.Close()
}

func (r bufferedFileReader) Read(p []byte) (n int, err error) {
	return r.bufferedReader.Read(p)
}
