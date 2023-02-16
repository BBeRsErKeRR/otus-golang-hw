package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	input, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer input.Close()

	stat, err := input.Stat()
	if err != nil {
		return err
	}
	size := stat.Size()

	if size == 0 {
		return ErrUnsupportedFile
	}
	if offset > size {
		return ErrOffsetExceedsFileSize
	}

	_, err = input.Seek(offset, 0)
	if err != nil {
		return err
	}

	totalSize := size - offset

	destination, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer destination.Close()

	pb := CreateNew64(totalSize)
	pb.Prefix(fmt.Sprintf("Copy file '%v' to '%v': ", fromPath, toPath))
	pb.Postfix(fmt.Sprintf(" | limit: %v, offset: %v ", limit, offset))
	pb.Start()
	defer pb.Finish()

	if limit == 0 {
		n, err := io.Copy(destination, input)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		pb.Add64(n)
	} else {
		for {
			n, err := io.CopyN(destination, input, limit)
			if err != nil && !errors.Is(err, io.EOF) {
				return err
			}
			pb.Add64(n)
			if n == 0 {
				break
			}
		}
	}
	return nil
}
