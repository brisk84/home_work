package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var (
	prevRune rune
	slash    bool
)

var ErrInvalidString = errors.New("invalid string")

func unpackRune(curRune rune) (string, error) {
	var outString strings.Builder
	if !unicode.IsDigit(curRune) {
		if (curRune != '\\') && (prevRune == '\\') {
			return "", ErrInvalidString
		}
		if (curRune != '\\') && (prevRune != 0) {
			outString.WriteRune(prevRune)
		}
		prevRune = curRune
		return outString.String(), nil
	}

	if prevRune == 0 {
		return "", ErrInvalidString
	}
	if (prevRune == '\\') && !slash {
		prevRune = curRune
		return outString.String(), nil
	}
	if slash {
		prevRune = '\\'
	}
	count, _ := strconv.Atoi(string(curRune))
	if count != 0 {
		outString.WriteString(strings.Repeat(string(prevRune), count))
	}
	prevRune = 0
	return outString.String(), nil
}

func Unpack(inputString string) (string, error) {
	var outString strings.Builder
	slash = false
	prevRune = 0
	for _, curRune := range inputString {
		s, err := unpackRune(curRune)
		if err != nil {
			return "", err
		}
		outString.WriteString(s)
	}
	if prevRune != 0 {
		outString.WriteRune(prevRune)
	}
	return outString.String(), nil
}
