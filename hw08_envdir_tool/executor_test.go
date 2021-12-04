package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	const refOut = `HELLO is ("hello")
BAR is (bar)
FOO is (   foo
with new line)
UNSET is ()
ADDED is (from original env)
EMPTY is ()
arguments are arg1=1 arg2=2
Exit code: 0
`
	env, err := ReadDir("testdata/env/")
	if err != nil {
		fmt.Println(err)
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	var cmd []string
	os.Setenv("ADDED", "from original env")
	cmd = append(cmd, "/bin/bash", "testdata/echo.sh", "arg1=1", "arg2=2")
	exitCode := RunCmd(cmd, env)
	os.Unsetenv("ADDED")
	fmt.Println("Exit code:", exitCode)

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	w.Close()
	os.Stdout = oldStdout
	out := <-outC

	require.Equal(t, refOut, out)
}
