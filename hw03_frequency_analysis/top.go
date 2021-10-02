package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var re *regexp.Regexp

func Top10(inputString string) []string {
	if len(inputString) == 0 {
		return nil
	}

	wordsSlice := strings.Split(re.ReplaceAllString(inputString, " "), " ")
	wordsCount := make(map[string]int)
	for _, word := range wordsSlice {
		lowerWord := strings.ToLower(word)
		wordsCount[lowerWord]++
	}

	countWords := make(map[int][]string)
	for k, v := range wordsCount {
		if (k != "") && (k != "-") {
			countWords[v] = append(countWords[v], k)
		}
	}

	keys := make([]int, 0, len(countWords))
	for k := range countWords {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(keys)))

	for _, v := range keys {
		sort.Strings(countWords[v])
	}

	tenWords := []string{}
	for _, v := range keys {
		for _, word := range countWords[v] {
			tenWords = append(tenWords, word)
			if len(tenWords) >= 10 {
				return tenWords
			}
		}
	}
	return tenWords
}

func init() {
	re = regexp.MustCompile(`[\s' ,.!]+`)
}
