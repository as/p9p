package p9p

import (
	"log"
)

var debug = 0
var logf = log.Printf

func init() {
	if debug == 0 {
		logf = func(s string, i ...interface{}) { return }
	}
}
