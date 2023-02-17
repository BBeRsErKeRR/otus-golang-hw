package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrorEqualFile           = errors.New("input and dst path is equal")
	ErrorNegativeNumber      = errors.New("negative number")
)

func resolveLink(link string) (string, error) {
	resolvedPath, err := os.Readlink(link)
	if err != nil {
		return "", err
	}
	if !filepath.IsAbs(resolvedPath) {
		resolvedPath = path.Join(path.Dir(link), resolvedPath)
	}
	return resolvedPath, nil
}

func checkFilePathEquals(from, to string) error {
	// check string
	if from == to {
		return ErrorEqualFile
	}
	// check relative
	fromPath, err := filepath.Abs(from)
	if err != nil {
		return err
	}
	toPath, err := filepath.Abs(to)
	if err != nil {
		return err
	}
	// check symbolic
	fromStat, err := os.Lstat(from)
	if err != nil {
		return err
	}
	toStat, err := os.Lstat(to)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if fromStat.Mode()&os.ModeSymlink != 0 {
		fromPath, err = resolveLink(fromPath)
		if err != nil {
			return err
		}
	}
	if toStat.Mode()&os.ModeSymlink != 0 {
		toPath, err = resolveLink(toPath)
		if err != nil {
			return err
		}
	}
	if fromPath == toPath {
		return ErrorEqualFile
	}
	return nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	err := checkFilePathEquals(fromPath, toPath)
	if err != nil {
		return err
	}

	if offset < 0 || limit < 0 {
		return ErrorNegativeNumber
	}
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

	_, err = input.Seek(offset, io.SeekStart)
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
	fmt.Printf("Copy file '%v' to '%v': \n", fromPath, toPath)
	fmt.Printf("  limit: %v\n  offset: %v \n", limit, offset)
	pb.Prefix("  progress:")
	pb.Start()
	defer pb.Finish()
	defer pb.Postfix("  SUCCESS")

	var reader io.Reader
	if limit > 0 {
		reader = io.LimitReader(input, limit)
	} else {
		reader = input
	}

	barReader := pb.NewProxyFreezeReader(reader, time.Millisecond*10)
	_, err = io.Copy(destination, barReader)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	return nil
}
