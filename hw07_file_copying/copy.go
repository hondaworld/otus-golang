package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrIdenticalFiles        = errors.New("source and destination are same")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == toPath {
		return ErrIdenticalFiles
	}

	sourceFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("failed to source file %s: %w", fromPath, err)
	}
	defer sourceFile.Close()

	sourceFileInfo, err := sourceFile.Stat()
	if err != nil {
		return fmt.Errorf("failed get info in file %s: %w", fromPath, err)
	}

	if sourceFileInfo.Size() == 0 {
		return ErrUnsupportedFile
	}

	if offset > sourceFileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	_, err = sourceFile.Seek(offset, io.SeekStart)
	if err != nil {
		return ErrOffsetExceedsFileSize
	}

	destinationFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("failed to destination file %s: %w", toPath, err)
	}
	defer destinationFile.Close()

	switch {
	case limit == 0:
		limit = sourceFileInfo.Size() - offset
	case limit > sourceFileInfo.Size()-offset:
		limit = sourceFileInfo.Size() - offset
	}

	bar := pb.Full.Start64(limit)
	defer bar.Finish()

	progressReader := bar.NewProxyReader(sourceFile)

	_, err = io.CopyN(destinationFile, progressReader, limit)
	if err != nil && !errors.Is(err, io.EOF) {
		os.Remove(toPath)
		return fmt.Errorf("failed in copy files: %w", err)
	}

	return nil
}
