package main

import (
	"fmt"
	"strconv"

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
	Bar  MagicInt
	Bazz int
}

func main() {
	c := Config{}
	opts.Parse(&c)
	fmt.Printf("%s %d\n", c.Bar, c.Bazz)
}
