package main

import (
	// "fmt"
	"testing"
	// "errors"
	// "encoding/json"
	// "io"
)

func TestCleanInput(t *testing.T) {

	cases := []struct {
		input string
		expected []string
	}{
		{
			input: " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input: "   	   PIkaCHU CHARMANDER      bulbasaur      ",
			expected: []string{"pikachu", "charmander", "bulbasaur"},
		},
		{
			input: "      			MEW MEW MeW MewTWO",
			expected: []string{"mew", "mew", "mew", "mewtwo"},
		},
	}
	
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("actual slice length: %v does not match length of expected slice: %v", len(actual), len(c.expected))
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("actual word: %v does not match expected word: %v from input string: %v", word, expectedWord, c.input)
			}
		}
	}
}
