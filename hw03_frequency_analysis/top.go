package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
	"unicode"
)

type kv struct {
	Word  string
	Count int
}

const PunctuationMarks = `.,?!'"\-:`

var re = regexp.MustCompile(`^[` + PunctuationMarks + `]*(.*?)[` + PunctuationMarks + `]+$`)

func Top10(str string) []string {
	arr := strings.Fields(str)
	var words = make(map[string]int, len(arr))

	for _, val := range arr {
		val = fixWord(val)

		if len(val) > 0 {
			words[val]++
		}
	}

	kvWords := make([]kv, 0, len(words))
	for word, count := range words {
		kvWords = append(kvWords, kv{word, count})
	}

	sort.Slice(kvWords, func(i, j int) bool {
		if kvWords[i].Count == kvWords[j].Count {
			return kvWords[i].Word < kvWords[j].Word
		}

		return kvWords[i].Count > kvWords[j].Count
	})

	result := make([]string, 0, 10)
	i := 0
	for _, val := range kvWords {
		if i == 10 {
			break
		}

		result = append(result, val.Word)
		i++
	}

	return result
}

func fixWord(word string) string {
	matches := re.FindAllStringSubmatch(word, -1)

	if len(matches) > 0 {
		word = matches[0][1]
	}

	return lowercaseLetters(word)
}

func lowercaseLetters(word string) string {
	if len(word) == 0 {
		return ""
	}

	var result strings.Builder

	for _, letter := range word {
		if unicode.IsUpper(letter) {
			letter = unicode.ToLower(letter)
		}

		result.WriteRune(letter)
	}

	return result.String()
}
