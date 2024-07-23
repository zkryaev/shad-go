//go:build !solution

package reverse

import (
	"strings"
	"unicode/utf8"
)

func Reverse(input string) string {
	var ans strings.Builder
	ans.Grow(len(input))
	for len(input) > 0 {
		r, size := utf8.DecodeLastRuneInString(input)
		ans.WriteRune(r)
		input = input[:len(input)-size]
	}
	return ans.String()
}
