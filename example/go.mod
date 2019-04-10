module gen-readme

go 1.12

replace github.com/jpillora/opts => ../

replace github.com/jpillora/md-tmpl => github.com/millergarym/md-tmpl v1.2.0

//replace github.com/jpillora/md-tmpl => github.com/millergarym/md-tmpl v1.2.0

require (
	github.com/jpillora/md-tmpl v0.0.0-20190330134846-187fa65bc021
	github.com/jpillora/opts v0.0.0-20160806153215-3f7962811f23
)
