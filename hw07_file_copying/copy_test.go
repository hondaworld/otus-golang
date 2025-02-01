package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const file1 = "testdata/input.txt"

func TestCopy(t *testing.T) {
	t.Run("offset 0, limit 0", func(t *testing.T) {
		file2, err := os.CreateTemp("", "test-*")
		file3 := "testdata/out_offset0_limit0.txt"

		if err != nil {
			log.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(file2.Name())

		errCopy := Copy(file1, file2.Name(), 0, 0)

		equal, err := filesAreEqual(file2.Name(), file3)

		file2.Close()

		require.Truef(t, equal, "files ars equal")
		require.Nil(t, errCopy)
		require.Nil(t, err)
	})
	t.Run("offset 0, limit 10", func(t *testing.T) {
		file2, err := os.CreateTemp("", "test-*")
		file3 := "testdata/out_offset0_limit10.txt"

		if err != nil {
			log.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(file2.Name())
		defer file2.Close()

		errCopy := Copy(file1, file2.Name(), 0, 10)

		equal, err := filesAreEqual(file2.Name(), file3)

		require.Truef(t, equal, "files ars equal")
		require.Nil(t, errCopy)
		require.Nil(t, err)
	})
	t.Run("offset 6000, limit 1000", func(t *testing.T) {
		file2, err := os.CreateTemp("", "test-*")
		file3 := "testdata/out_offset6000_limit1000.txt"

		if err != nil {
			log.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(file2.Name())
		defer file2.Close()

		errCopy := Copy(file1, file2.Name(), 6000, 1000)

		equal, err := filesAreEqual(file2.Name(), file3)

		require.Truef(t, equal, "files ars equal")
		require.Nil(t, errCopy)
		require.Nil(t, err)
	})
	t.Run("offset 10000, limit 1000", func(t *testing.T) {
		file2, err := os.CreateTemp("", "test-*")

		if err != nil {
			log.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(file2.Name())
		defer file2.Close()

		errCopy := Copy(file1, file2.Name(), 10000, 1000)

		require.Truef(t, errors.Is(errCopy, ErrOffsetExceedsFileSize), "actual err - %v", errCopy)
	})
	t.Run("not file exists", func(t *testing.T) {
		file1Wrong := "testdata1/input.txt"
		file2, err := os.CreateTemp("", "test-*")

		if err != nil {
			log.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(file2.Name())
		defer file2.Close()

		errCopy := Copy(file1Wrong, file2.Name(), 10000, 1000)

		require.Truef(t, errors.Is(errCopy, ErrUnsupportedFile), "actual err - %v", errCopy)
	})
}

func filesAreEqual(file1, file2 string) (bool, error) {
	sourceFile1, err := os.Open(file1)
	if err != nil {
		return false, err
	}
	defer sourceFile1.Close()

	sourceFile2, err := os.Open(file2)
	if err != nil {
		return false, err
	}
	defer sourceFile2.Close()

	const bufSize = 4096
	buf1 := make([]byte, bufSize)
	buf2 := make([]byte, bufSize)

	for {
		n1, err1 := sourceFile1.Read(buf1)
		n2, err2 := sourceFile2.Read(buf2)

		if n1 != n2 || (err1 == io.EOF && err2 != io.EOF) || (err1 != io.EOF && err2 == io.EOF) {
			return false, nil
		}

		if err1 == io.EOF && err2 == io.EOF {
			break
		}

		if !bytes.Equal(buf1[:n1], buf2[:n2]) {
			return false, nil
		}

		if err1 != nil && err1 != io.EOF {
			return false, err1
		}
		if err2 != nil && err2 != io.EOF {
			return false, err2
		}
	}

	return true, nil
}
