package flag

import "strings"

func camel2dash(s string) string {
	return strings.ToLower(s)
}

func nletters(r rune, n int) string {
	str := make([]rune, n)
	for i, _ := range str {
		str[i] = r
	}
	return string(str)
}
