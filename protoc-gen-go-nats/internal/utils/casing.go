package utils

import "strings"

func ToFirstLowerCase(str string) string {
	if len(str) == 0 {
		return str
	}

	var b strings.Builder

	b.WriteString(strings.ToLower(string(str[0])))
	b.WriteString(str[1:])

	return b.String()

}
