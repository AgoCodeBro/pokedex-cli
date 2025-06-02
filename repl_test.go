package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := [] struct {
		input    string
		expected []string
	}{
		{
			input    : " hello  world ",
			expected : []string{"hello", "world"},
		},
		{
			input    : "These are SOME words",
			expected : []string{"these", "are", "some", "words"},
		},
		{
			input    : "WhOps IDoknt Nownho w to spleLL",
			expected : []string{"whops", "idoknt", "nownho", "w", "to", "splell"},
		},
	}


	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Expected length: %v \nActual: %v", len(c.expected), len(actual))
		}

		for i, word := range actual {
			if expectedWord := c.expected[i]; word != expectedWord {
				t.Errorf("Expected word at pos %v: %v\nActual %v", i, expectedWord, word)
			}
		}
	}
}
 
