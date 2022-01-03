package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	l := New("/tmp/log.txt", "debug")
	l.Info("Test")
	require.True(t, false)
}
