package opts

import (
	"bytes"
	"strings"
)

func camel2dash(s string) string {
	return strings.ToLower(s)
}

func camel2const(s string) string {
	b := bytes.Buffer{}
	var c rune
	start := 0
	end := 0
	for end, c = range s {
		if c >= 'A' && c <= 'Z' {
			//uppercase all prior letters and add an underscore
			if start < end {
				b.WriteString(strings.ToTitle(s[start:end] + "_"))
				start = end
			}
		}
	}
	//write remaining string
	b.WriteString(strings.ToTitle(s[start : end+1]))
	return b.String()
}

func nletters(r rune, n int) string {
	str := make([]rune, n)
	for i, _ := range str {
		str[i] = r
	}
	return string(str)
}
