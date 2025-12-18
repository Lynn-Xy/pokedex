package main

import (
	"strings"
)

func cleanInput(s string) []string {
	var results []string
	lowerString := strings.ToLower(s)
	delimitedString := strings.TrimSpace(lowerString)
	results = strings.Fields(delimitedString)
	return results
}
