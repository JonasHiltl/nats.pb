package utils

import "testing"

func TestToFirstLowerCase(t *testing.T) {
	casingTests := []struct {
		input string
		want  string
	}{
		{
			"CamelCase",
			"camelCase",
		},
		{
			"C",
			"c",
		},
		{
			"",
			"",
		},
	}

	for _, ct := range casingTests {
		got := ToFirstLowerCase(ct.input)
		if ct.want != got {
			t.Errorf("want: %s, got: %s", ct.want, got)
		}
	}
}
