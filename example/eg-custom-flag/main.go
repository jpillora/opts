package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jpillora/opts"
)

//MagicInt is a valid custom type since it implements the flag.Value interface
type MagicInt int

func (b MagicInt) String() string {
	return "{" + strconv.Itoa(int(b)) + "}"
}

func (b *MagicInt) Set(s string) error {
	n, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*b = MagicInt(n + 42)
	return nil
}

type Config struct {
	Mmm   []MagicInt
	Bar   time.Duration
	Zee   bool
	Files []opts.File
	Dir   opts.Dir
}

func main() {
	c := Config{}
	opts.Parse(&c)
	fmt.Printf("%+v\n", c)
}
