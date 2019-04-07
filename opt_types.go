package opts

import "fmt"

type RepeatedStringOpt struct {
	vals []string
}

func (os *RepeatedStringOpt) Set(val string) error {
	os.vals = append(os.vals, val)
	return nil
}

func (os RepeatedStringOpt) String() string {
	if len(os.vals) == 0 {
		return ""
	}
	if len(os.vals) == 1 {
		return "'" + os.vals[0] + "'"
	}
	return fmt.Sprintf("%v", os.vals)
}

func (os *RepeatedStringOpt) GetSlice() []string {
	return os.vals
}
