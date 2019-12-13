package main

import (
	"testing"
)

func TestCleanTrailingCommas(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "nothing to do",
			input:  `{"foo": "bar"}`,
			output: `{"foo": "bar"}`,
		},
		{
			name:   "trailing comma in object",
			input:  `{"foo": "bar",}`,
			output: `{"foo": "bar"}`,
		},
		{
			name:   "trailing comman in array",
			input:  `[1,2,3,]`,
			output: `[1,2,3]`,
		},
		{
			name: "newlines don't matter",
			input: `{
				"foo": 1,
			}`,
			output: `{
				"foo": 1
			}`,
		},
		{
			name:   "final trailing comma",
			input:  `{"foo": "bar"},`,
			output: `{"foo": "bar"}`,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			out, err := cleanTrailingCommas([]byte(c.input))
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if string(out) != c.output {
				t.Errorf("incorrect output: %s", out)
			}
		})
	}
}
