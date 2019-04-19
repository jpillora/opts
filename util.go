package opts

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

var flagValueType = reflect.TypeOf((*flag.Value)(nil)).Elem()

var durationType = reflect.TypeOf(time.Second)

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

func str2str(src string, dst *string) {
	if src != "" {
		*dst = src
	}
}

func str2bool(src string, dst *bool) {
	if src != "" {
		*dst = strings.ToLower(src) == "true" || src == "1"
	}
}

func str2int(src string, dst *int) {
	if src != "" {
		n, err := strconv.Atoi(src)
		if err == nil {
			*dst = n
		}
	}
}

func linkFlagset(flags []*item, flagset *flag.FlagSet) error {
	for _, opt := range flags {
		//2. set config via environment
		envVal := ""
		if opt.useEnv {
			envVal = os.Getenv(opt.envName)
		}
		//3. set config via Go's pkg/flags
		addr := opt.val.Addr().Interface()
		switch addr := addr.(type) {
		case flag.Value:
			flagset.Var(addr, opt.name, "")
			if opt.shortName != "" {
				flagset.Var(addr, opt.shortName, "")
			}
		case *[]string:
			sep := ""
			switch opt.typeName {
			case "commalist":
				sep = ","
			case "spacelist":
				sep = " "
			}
			fv := &sepList{sep: sep, strs: addr}
			flagset.Var(fv, opt.name, "")
			if opt.shortName != "" {
				flagset.Var(fv, opt.shortName, "")
			}
		case *bool:
			str2bool(envVal, addr)
			flagset.BoolVar(addr, opt.name, *addr, "")
			if opt.shortName != "" {
				flagset.BoolVar(addr, opt.shortName, *addr, "")
			}
		case *string:
			str2str(envVal, addr)
			flagset.StringVar(addr, opt.name, *addr, "")
			if opt.shortName != "" {
				flagset.StringVar(addr, opt.shortName, *addr, "")
			}
		case *int:
			str2int(envVal, addr)
			flagset.IntVar(addr, opt.name, *addr, "")
			if opt.shortName != "" {
				flagset.IntVar(addr, opt.shortName, *addr, "")
			}
		case *time.Duration:
			flagset.DurationVar(addr, opt.name, *addr, "")
			if opt.shortName != "" {
				flagset.DurationVar(addr, opt.shortName, *addr, "")
			}
		default:
			return fmt.Errorf("[opts] Option '%s' has unsupported type", opt.name)
		}
	}
	return nil
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

//sepList provides a flag value strings interface
//over a string slice
type sepList struct {
	sep  string
	strs *[]string
}

func (l sepList) String() string {
	if l.strs != nil {
		return strings.Join(*l.strs, l.sep)
	}
	return ""
}

func (l sepList) Set(s string) error {
	*l.strs = strings.Split(s, l.sep)
	return nil
}
