//go:build !solution

package varfmt

import (
	"fmt"
	"strconv"
	"strings"
)

func Sprintf(format string, args ...interface{}) string {
	var ans strings.Builder
	var ind strings.Builder
	baseind := 0
	i := 0
	for ; i < len(format); i++ {
		switch {
		case format[i] == '{':
			flag := false
			j := i + 1
			for j < len(format) && format[j] != '}' {
				ind.WriteRune(rune(format[j]))
				flag = true
				j++
			}
			i = j
			switch {
			case flag:
				setind, _ := strconv.Atoi(ind.String()) // setind - указанный индекс
				ind.Reset()
				ans.WriteString(fmt.Sprint(args[setind]))
			case !flag:
				ans.WriteString(fmt.Sprint(args[baseind])) // baseind - естественный индекс
			}
			baseind++
		default:
			ans.WriteRune(rune(format[i]))
		}
	}
	return ans.String()
}
