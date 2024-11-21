package search

import (
	"strings"
	"unicode"

	snowballeng "github.com/kljensen/snowball/english"
)

func analyze(text string) []string {
	tokens := tokenize(text)
	tokens = lowercaseFilter(tokens)
	tokens = stopwordFilter(tokens)
	tokens = stemmerFilter(tokens)
	return tokens
}

func tokenize(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

// filter out all tokens that are not letters or numbers and convert all tokens to lowercase -> ["Hello", "world"] -> ["hello", "world"]
func lowercaseFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = strings.ToLower(token)
	}
	return r
}

// filter out stop words in query "i wanted this baseball cap" -> "wanted baseball cap"
func stopwordFilter(tokens []string) []string {
	r := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if !snowballeng.IsStopWord(token) {
			r = append(r, token)
		}
	}

	return r
}

func stemmerFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = snowballeng.Stem(token, false)
	}

	return r
}
