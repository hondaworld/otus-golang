package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	sourceFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0)
	if err != nil {
		return ErrUnsupportedFile
	}
	defer sourceFile.Close()

	sourceFileInfo, err := sourceFile.Stat()
	if err != nil {
		return err
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
		return ErrUnsupportedFile
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
		return err
	}

	return nil
}
