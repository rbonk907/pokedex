package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "    hello world    ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charizard BULBasaur PIKACHU  ",
			expected: []string{"charizard", "bulbasaur", "pikachu"},
		},
		{
			input:    "",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("cleanInput returned slice of length: %v, expecting length: %v", len(actual), len(c.expected))
		}

		for i := range actual {
			word := actual[i]
			expected := c.expected[i]

			if word != expected {
				t.Errorf("expected word: %s, actual word received: %s", expected, word)
			}
		}
	}
}
