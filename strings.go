package db

import (
	"reflect"
	"strings"
)

func SplitStringFunc(s string, f func(rune) bool) []string {
	// A span is used to record a slice of s of the form s[start:end].
	// The start index is inclusive and the end index is exclusive.
	type span struct {
		start int
		end   int
	}
	spans := make([]span, 0, 32)

	start := 0
	for end, r := range s {
		if f(r) {
			if end > start {
				spans = append(spans, span{start, end})
				start = end
			}
		}
	}

	// Last field might end at EOF.
	if start < len(s) {
		spans = append(spans, span{start, len(s)})
	}

	// Create strings from recorded field indices.
	a := make([]string, len(spans))
	for i, span := range spans {
		a[i] = s[span.start:span.end]
	}

	return a
}

func SnakeCase(s string) string {
	words := SplitStringFunc(s, func(r rune) bool {
		return r >= 'A' && r <= 'Z'
	})

	for i, word := range words {
		words[i] = strings.ToLower(word)
	}

	return strings.Join(words, "_")
}

func ModelName(model interface{}) string {
	typ := reflect.TypeOf(model)
	words := strings.Split(typ.String(), ".")
	return words[len(words)-1]
}
