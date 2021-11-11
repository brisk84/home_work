package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		// uncomment if task with asterisk completed
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
		{input: `qwe//#$%5`, expected: `qwe//#$%%%%%`},
		{input: `qwe\03bc0d`, expected: `qwe000bd`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestEscapeChars(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "d\n5abc", expected: "d\n\n\n\n\nabc"},
		{input: "E\0333s\t3#$%", expected: "E\033\033\033s\t\t\t#$%"},
		{input: "z\0330\n1abc", expected: "z\nabc"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackCyrilic(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "аб4в3г2", expected: "аббббвввгг"},
		{input: "Абвгд", expected: "Абвгд"},
		{input: "ааа0б", expected: "ааб"},
		{input: `абц//#$%5`, expected: `абц//#$%%%%%`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackUnicode(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "世3界2", expected: "世世世界界"},
		{input: "世世界界", expected: "世世界界"},
		{input: "世世世0界", expected: "世世界"},
		{input: `世世世//#$界5`, expected: `世世世//#$界界界界界`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b", `qw\ne`}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
