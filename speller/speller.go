//go:build !solution

package speller

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	hundred   = 1
	thousands = 4
	millions  = 7
	billions  = 10
	trillions = 13
)

var prefixs = map[int]string{
	0:  "zero",
	1:  "one",
	2:  "two",
	3:  "three",
	4:  "four",
	5:  "five",
	6:  "six",
	7:  "seven",
	8:  "eight",
	9:  "nine",
	10: "ten",
	11: "eleven",
	12: "twelve",
	13: "thirteen",
	14: "fourteen",
	15: "fifteen",
	16: "sixteen",
	17: "seventeen",
	18: "eighteen",
	19: "nineteen",
	20: "twenty",
	30: "thirty",
	40: "forty",
	50: "fifty",
	60: "sixty",
	70: "seventy",
	80: "eighty",
	90: "ninety",
}

func DecodeThreeBlock(ThreeBlock string) string {
	var ans strings.Builder
	number, _ := strconv.Atoi(ThreeBlock)
	if number == 0 && len(ThreeBlock) > 1 {
		return ""
	}
	ans.Grow(3)
	hundred := number / 100 // 121/100 = 1
	last := number % 100    // 121%100 = 21
	tens := last / 10       // 21/10 = 2
	digits := last % 10     // 21%10 = 1
	if hundred >= 1 {
		ans.WriteString(fmt.Sprintf("%s hundred", prefixs[hundred]))
	}
	if number%100 == 0 && len(ThreeBlock) > 1 {
		return ans.String()
	} else if ans.Len() != 0 {
		ans.WriteString(" ")
	}
	if last >= 21 {
		ans.WriteString(prefixs[tens*10])
		if digits != 0 {
			ans.WriteString(fmt.Sprintf("-%s", prefixs[digits]))
		}
	} else {
		ans.WriteString(prefixs[last])
	}
	return ans.String()
}

func GetBlock(number string, begin *int, prefix int) string {
	ThreeBlock := make([]rune, 0, 3)
	i := (*begin)
	for ; len(number)-i >= prefix; i++ {
		ThreeBlock = append(ThreeBlock, rune(number[i]))
	}
	(*begin) = i
	return string(ThreeBlock)
}

func Spell(n int64) string {
	var ans strings.Builder
	number := strconv.FormatInt(n, 10)
	NumBlocks := len(number) / 3
	if len(number)%3 != 0 {
		NumBlocks++
	}
	i := 0
	if number[i] == '-' {
		ans.WriteString("minus")
		i += 1
	}
	for block := 1; block <= NumBlocks; block++ {
		switch {
		case len(number)-i >= trillions:
			if TrillionsBlock := DecodeThreeBlock(GetBlock(number, &i, trillions)); len(
				TrillionsBlock,
			) != 0 {
				ans.WriteString(fmt.Sprintf("%s trillion", TrillionsBlock))
			}

			i++
		case len(number)-i >= billions:
			if BillionsBlock := DecodeThreeBlock(GetBlock(number, &i, billions)); len(
				BillionsBlock,
			) != 0 {
				if ans.Len() != 0 {
					ans.WriteString(" ")
				}
				ans.WriteString(fmt.Sprintf("%s billion", BillionsBlock))
			}
		case len(number)-i >= millions:
			if MillionsBlock := DecodeThreeBlock(GetBlock(number, &i, millions)); len(
				MillionsBlock,
			) != 0 {
				if ans.Len() != 0 {
					ans.WriteString(" ")
				}
				ans.WriteString(fmt.Sprintf("%s million", MillionsBlock))
			}
		case len(number)-i >= thousands:
			if ThousandsBlock := DecodeThreeBlock(GetBlock(number, &i, thousands)); len(
				ThousandsBlock,
			) != 0 {
				if ans.Len() != 0 {
					ans.WriteString(" ")
				}
				ans.WriteString(fmt.Sprintf("%s thousand", ThousandsBlock))
			}
		case len(number)-i >= 1:
			if HudredsBlock := DecodeThreeBlock(GetBlock(number, &i, hundred)); len(
				HudredsBlock,
			) != 0 {
				if ans.Len() != 0 {
					ans.WriteString(" ")
				}
				ans.WriteString(HudredsBlock)
			}
		}
	}
	return ans.String()
}
