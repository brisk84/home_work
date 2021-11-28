package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func ShowProgressbar(progress int) {
	var progressbar = "----------------------------------------------------------------------------------------------------"
	progressbar = strings.Replace(progressbar, "-", "+", progress)
	fmt.Printf("\r%v [%d%%]", progressbar, progress)
}

func Copy(fromPath, toPath string, offset, limit int64) error {

	const bufSize = 64
	var (
		curBufSize  int64
		bytesCopyed int64
		incProgress float64
		progress    int
	)

	fileInfo, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()
	if fileSize == 0 {
		return ErrUnsupportedFile
	}
	if (offset > fileSize) && (fileSize != 0) {
		return ErrOffsetExceedsFileSize
	}
	if limit == 0 {
		limit = fileSize
	}

	fileFrom, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fileFrom.Close()

	curBufSize = limit
	if limit >= bufSize {
		curBufSize = bufSize
	}
	buf := make([]byte, curBufSize)

	fileTo, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer fileTo.Close()

	if offset > 0 {
		_, err := fileFrom.Seek(offset, 0)
		if err != nil {
			return err
		}
	}

	realSizeToCopy := fileSize - offset
	biggestSize := limit
	if realSizeToCopy < limit {
		biggestSize = realSizeToCopy
	}
	incSize := float64(bufSize) * 100 / float64(biggestSize)
	for {
		if bytesCopyed >= limit {
			break
		}
		if bytesCopyed >= (limit - curBufSize) {
			buf = make([]byte, limit-bytesCopyed)
		}
		bytesRead, err := fileFrom.Read(buf)
		bytesCopyed += int64(bytesRead)
		if _, err := fileTo.Write(buf[:bytesRead]); err != nil {
			return err
		}
		if err == io.EOF {
			break
		}

		incProgress += incSize
		if incProgress >= 1 {
			progress++
			incProgress--
		}
		ShowProgressbar(progress)
		time.Sleep(bufSize / 2 * time.Millisecond)
	}
	ShowProgressbar(100)

	return nil
}
