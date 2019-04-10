## Completion example with custom prediction 

<!--tmpl,chomp,code=go:cat main.go -->
``` go 
package main

import (
	"fmt"

	"github.com/jpillora/opts"
	"github.com/posener/complete"
)

type Config struct {
	// pstring is a custom type with a Predict function
	Alpha pstring
	Bar
	Foo
}

type Foo struct {
	Ping string
	Pong string
}

type Bar struct {
	Zip string
	Zop string
}

type pstring struct {
	val string
}

func (*pstring) Predict(args complete.Args) []string {
	return []string{"a", "b", "c"}
}
func (p *pstring) Set(val string) (err error) {
	p.val = val
	return
}
func (p *pstring) String() string {
	return p.val
}

type foo struct {
}
type man struct{}
type chew struct{}

func main() {
	config := Config{}
	// The Name must match the exec's name for self completion to work
	opts.New(&config).Name("eg-complete").
		Complete().
		AddCommand(opts.New(&foo{}).Name("fooie").AddCommand(
			opts.New(&man{}).Name("man").AddCommand(
				opts.New(&chew{}).Name("chew")))).
		Parse().
		RunFatal()
}

func (obj *Config) Run() {
	fmt.Printf("%+v\n", obj)
}
```
<!--/tmpl-->

### Build example and get the help text
Self installing self completing executable.

`go build -o eg-complete && ./eg-complete --help`
<!--tmpl,chomp,code=plain:go build -o eg-complete && ./eg-complete --help -->
``` plain 

  Usage: eg-complete [options] <command>

  Options:
  --alpha, -a      default {}
  --help, -h       default false
  --install, -i    install shell-completion (default false)
  --uninstall, -u  uninstall shell-completion (default false)

  Commands:
  • fooie
  • bar
  • foo

```
<!--/tmpl-->

### Install the completion in bash, zsh and fish
`./eg-complete -i`
<!--tmpl,chomp,code=plain:./eg-complete -i -->
``` plain 
Installed
```
<!--/tmpl-->

attempt second install

`./eg-complete -i`
<!--tmpl,chomp,code=plain:./eg-complete -i -->
``` plain 
1 error occurred:
	* already installed in /home/garym/.bashrc

```
<!--/tmpl-->


### Get completion
`COMP_LINE="./eg-complete " ./eg-complete`
<!--tmpl,chomp,code=plain:COMP_DEBUG= COMP_LINE="./eg-complete " ./eg-complete -->
``` plain 
bar
foo
fooie
```
<!--/tmpl-->

### Completion of flags
`COMP_LINE="./eg-complete -" ./eg-complete `
<!--tmpl,chomp,code=plain:COMP_DEBUG= COMP_LINE="./eg-complete -" ./eg-complete -->
``` plain 
-a
--help
-h
--install
-i
--uninstall
-u
--alpha
```
<!--/tmpl-->


### Custom completion
`COMP_LINE="./eg-complete -a " ./eg-complete `
<!--tmpl,chomp,code=plain:COMP_DEBUG= COMP_LINE="./eg-complete -a " ./eg-complete -->
``` plain 
a
b
c
```
<!--/tmpl-->

### Custom completion with debug
`COMP_DEBUG=1 COMP_LINE="./eg-complete -a " ./eg-complete `
<!--tmpl,chomp,code=plain:COMP_DEBUG=1 COMP_LINE="./eg-complete -a " ./eg-complete -->
``` plain 
complete 2019/04/10 21:58:15 flag completion alpha &main.pstring{val:""}
complete 2019/04/10 21:58:15 flag completion help (complete.PredictFunc)(0x50ceb0)
complete 2019/04/10 21:58:15 flag completion install (complete.PredictFunc)(0x50ceb0)
complete 2019/04/10 21:58:15 flag completion uninstall (complete.PredictFunc)(0x50ceb0)
complete 2019/04/10 21:58:15 Failed parsing point : strconv.Atoi: parsing "": invalid syntax
complete 2019/04/10 21:58:15 Completing phrase: ./eg-complete -a 
complete 2019/04/10 21:58:15 Completing last field: 
complete 2019/04/10 21:58:15 Predicting according to flag -a
complete 2019/04/10 21:58:15 Options: [a b c]
complete 2019/04/10 21:58:15 Matches: [a b c]
a
b
c
```
<!--/tmpl-->

### Completion sub command
<!--/tmpl-->
`COMP_LINE="./eg-complete fooie " ./eg-complete`
<!--tmpl,chomp,code=plain:COMP_DEBUG= COMP_LINE="./eg-complete fooie " ./eg-complete -->
``` plain 
man
```
<!--/tmpl-->

### Uninstall
`./eg-complete -u`
<!--tmpl,chomp,code=plain:./eg-complete -u -->
``` plain 
Uninstalled
```
<!--/tmpl-->
