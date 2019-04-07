## defaults example

<!--tmpl,chomp,code=go:cat main.go -->
``` go 
cat: main.go: No such file or directory
```
<!--/tmpl-->

```
$ defaults --foo hello
```

<!--tmpl,chomp,code=plain:go run main.go --foo hello -->
``` plain 
stat main.go: no such file or directory
```
<!--/tmpl-->

```
$ defaults --help
```

<!--tmpl,chomp,code=plain:go run main.go --help -->
``` plain 
stat main.go: no such file or directory
```
<!--/tmpl-->
