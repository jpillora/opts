#!/bin/bash

#this file generates the READMEs in all the examples.
#requires github.com/jpillora/md-tmpl
for i in *; do
	if [ -d "$i" ] && [ -f "$i/README.md" ]; then
		cd "$i"
		echo "$i"
		md-tmpl README.md || exit 1
		cd ..
	fi
done
