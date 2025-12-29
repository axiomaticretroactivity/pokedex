package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    "can i get a mfkin uhhh     soda  ",
			expected: []string{"can", "i", "get", "a", "mfkin", "uhhh", "soda"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(c.expected) != len(actual) {
			t.Errorf("returned slice from cleanInput not the same length as expected")
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if len(word) != len(expectedWord) {
				t.Errorf(`actual word "%s" is not expected word "%s"`, word, expectedWord)
			}
		}
	}
}
