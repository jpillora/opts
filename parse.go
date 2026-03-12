package opts

import "fmt"

// parseFlags parses command-line flags from args using the provided flag map.
// When stopAtNonFlag is true, parsing stops at the first non-flag argument
// (used for subcommand boundaries). Otherwise, flags and positional arguments
// can be freely interspersed.
func parseFlags(flags map[string]*item, args []string, stopAtNonFlag bool) (remaining []string, err error) {
	i := 0
	for i < len(args) {
		arg := args[i]
		// "--" terminates flag parsing
		if arg == "--" {
			remaining = append(remaining, args[i+1:]...)
			return
		}
		// not a flag
		if len(arg) == 0 || arg[0] != '-' {
			if stopAtNonFlag {
				remaining = append(remaining, args[i:]...)
				return
			}
			remaining = append(remaining, arg)
			i++
			continue
		}
		// single "-" is not a flag
		if arg == "-" {
			if stopAtNonFlag {
				remaining = append(remaining, args[i:]...)
				return
			}
			remaining = append(remaining, arg)
			i++
			continue
		}
		// flag argument: strip leading dashes
		name := arg[1:]
		if name[0] == '-' {
			name = name[1:]
		}
		if name == "" {
			return remaining, fmt.Errorf("bad flag syntax: %s", arg)
		}
		// handle --flag=value
		value := ""
		hasValue := false
		if eqIdx := indexOf(name, '='); eqIdx >= 0 {
			value = name[eqIdx+1:]
			name = name[:eqIdx]
			hasValue = true
		}
		item, ok := flags[name]
		if !ok {
			return remaining, fmt.Errorf("unknown flag: %s", arg)
		}
		// bool flags don't consume next arg
		if item.IsBoolFlag() {
			if hasValue {
				if err := item.Set(value); err != nil {
					return remaining, fmt.Errorf("invalid value %q for flag -%s: %s", value, name, err)
				}
			} else {
				if err := item.Set("true"); err != nil {
					return remaining, fmt.Errorf("invalid value \"true\" for flag -%s: %s", name, err)
				}
			}
			i++
			continue
		}
		// non-bool flag needs a value
		if hasValue {
			if err := item.Set(value); err != nil {
				return remaining, fmt.Errorf("invalid value %q for flag -%s: %s", value, name, err)
			}
			i++
		} else if i+1 < len(args) {
			if err := item.Set(args[i+1]); err != nil {
				return remaining, fmt.Errorf("invalid value %q for flag -%s: %s", args[i+1], name, err)
			}
			i += 2
		} else {
			return remaining, fmt.Errorf("flag needs an argument: -%s", name)
		}
	}
	return
}

func indexOf(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}
