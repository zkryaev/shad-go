//go:build !solution

package spacecollapse

import (
	"strings"
	"unicode"
)

func CollapseSpaces(input string) string {
	var collapsedString strings.Builder
	flag := false
	for _, r := range input {
		if unicode.IsSpace(r) {
			if !flag {
				collapsedString.WriteRune(' ')
				flag = true
			}
		} else {
			collapsedString.WriteRune(r)
			flag = false
		}
	}
	return collapsedString.String()
}
