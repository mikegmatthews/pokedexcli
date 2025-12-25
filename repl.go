package main

import "strings"

func cleanInput(text string) []string {
	cleaned := make([]string, 0)

	for f := range strings.FieldsSeq(text) {
		cleaned = append(cleaned, strings.ToLower(f))
	}

	return cleaned
}
