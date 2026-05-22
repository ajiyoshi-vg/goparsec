module github.com/ajiyoshi-vg/goparsec/bench

go 1.26.3

replace (
	github.com/ajiyoshi-vg/goparsec => ../
	github.com/ajiyoshi-vg/parseg => /tmp/parseg
)

require (
	github.com/ajiyoshi-vg/goparsec v0.0.0-00010101000000-000000000000
	github.com/ajiyoshi-vg/parseg v0.0.0-00010101000000-000000000000
)
