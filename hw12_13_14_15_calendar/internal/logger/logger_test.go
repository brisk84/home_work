package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	l := New("/tmp/log.txt", "error")
	l.Info("Test333")
	l.Error("Error111")
	require.True(t, false)
}

func TestLogger2(t *testing.T) {
	l := New("stdout", "info")
	l.Info("Test333")
	l.Error("Error111")
	require.True(t, false)
}
