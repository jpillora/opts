package opts

import (
	"bytes"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

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
	for i := range str {
		str[i] = r
	}
	return string(str)
}

func constrain(str string, maxWidth int) string {
	lines := strings.Split(str, "\n")
	for i, line := range lines {
		words := strings.Split(line, " ")
		width := 0
		for i, w := range words {
			remain := maxWidth - width
			wordWidth := len(w) + 1 //+space
			width += wordWidth
			overflow := width > maxWidth
			fits := width-maxWidth > remain
			if overflow && fits {
				width = wordWidth
				w = "\n" + w
			}
			words[i] = w
		}
		lines[i] = strings.Join(words, " ")
	}
	return strings.Join(lines, "\n")
}

//borrowed from https://github.com/huandu/xstrings/blob/master/convert.go#L77
func camel2dash(str string) string {
	if len(str) == 0 {
		return ""
	}
	buf := &bytes.Buffer{}
	var prev, r0, r1 rune
	var size int
	r0 = '-'
	for len(str) > 0 {
		prev = r0
		r0, size = utf8.DecodeRuneInString(str)
		str = str[size:]
		switch {
		case r0 == utf8.RuneError:
			buf.WriteByte(byte(str[0]))
		case unicode.IsUpper(r0):
			if prev != '-' {
				buf.WriteRune('-')
			}
			buf.WriteRune(unicode.ToLower(r0))
			if len(str) == 0 {
				break
			}
			r0, size = utf8.DecodeRuneInString(str)
			str = str[size:]
			if !unicode.IsUpper(r0) {
				buf.WriteRune(r0)
				break
			}
			// find next non-upper-case character and insert `_` properly.
			// it's designed to convert `HTTPServer` to `http_server`.
			// if there are more than 2 adjacent upper case characters in a word,
			// treat them as an abbreviation plus a normal word.
			for len(str) > 0 {
				r1 = r0
				r0, size = utf8.DecodeRuneInString(str)
				str = str[size:]
				if r0 == utf8.RuneError {
					buf.WriteRune(unicode.ToLower(r1))
					buf.WriteByte(byte(str[0]))
					break
				}
				if !unicode.IsUpper(r0) {
					if r0 == '-' || r0 == ' ' || r0 == '_' {
						r0 = '-'
						buf.WriteRune(unicode.ToLower(r1))
					} else {
						buf.WriteRune('-')
						buf.WriteRune(unicode.ToLower(r1))
						buf.WriteRune(r0)
					}
					break
				}
				buf.WriteRune(unicode.ToLower(r1))
			}
			if len(str) == 0 || r0 == '-' {
				buf.WriteRune(unicode.ToLower(r0))
				break
			}
		default:
			if r0 == ' ' || r0 == '_' {
				r0 = '-'
			}
			buf.WriteRune(r0)
		}
	}
	return buf.String()
}

type kv struct {
	m map[string]string
}

func (kv *kv) keys() []string {
	ks := []string{}
	for k := range kv.m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}

func (kv *kv) take(k string) (string, bool) {
	v, ok := kv.m[k]
	if ok {
		delete(kv.m, k)
	}
	return v, ok
}

func newKV(s string) *kv {
	m := map[string]string{}
	key := ""
	mode := true
	sb := strings.Builder{}
	for _, r := range s {
		//key done
		if mode && r == '=' {
			key = sb.String()
			sb.Reset()
			mode = false
			continue
		}
		//value done
		if r == ',' {
			val := sb.String()
			sb.Reset()
			m[key] = val
			key = ""
			val = ""
			mode = true
			continue
		}
		//write to builder
		sb.WriteRune(r)
	}
	//write last key=value
	if key != "" {
		val := sb.String()
		m[key] = val
	}
	return &kv{m: m}
}
