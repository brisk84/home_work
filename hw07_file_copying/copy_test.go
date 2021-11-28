package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnsupportedFile(t *testing.T) {
	fromFile := "/dev/urandom"
	toFile := "testdata/output.txt"

	err := Copy(fromFile, toFile, 10, 2000)

	require.Error(t, err)
}

func TestOffset(t *testing.T) {
	fromFile := "testdata/input.txt"
	toFile := "testdata/output.txt"

	err := Copy(fromFile, toFile, 7000, 100)
	os.Remove(toFile)

	require.Error(t, err)
}

func TestLimit(t *testing.T) {
	fromFile := "testdata/input.txt"
	toFile := "testdata/output.txt"

	err := Copy(fromFile, toFile, 0, 7000)
	os.Remove(toFile)

	require.NoError(t, err)
}
