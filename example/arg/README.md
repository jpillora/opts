## arg example

<tmpl,code=go:cat arg.go>
``` go 

```
</tmpl>
```
$ arg --foo hello --bar world
```
<tmpl,code:go run arg.go --foo hello --bar world>
``` plain 
flag provided but not defined: -foo

  Usage: arg [options] <foo>
  
  foo is a very important argument
  
  Options:
  --bar, -b 
  --help, -h
  
```
</tmpl>
```
$ arg --help
```
<tmpl,code:go run arg.go --help>
``` plain 

  Usage: arg [options] <foo>
  
  foo is a very important argument
  
  Options:
  --bar, -b 
  --help, -h
  
```
</tmpl>