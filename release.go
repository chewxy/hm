// +build !debug

package hm

var READMEMSTATS = true

var TABCOUNT uint32 = 0

func tabcount() int { return 0 }

func enterLoggingContext()                      {}
func leaveLoggingContext()                      {}
func logf(format string, others ...interface{}) {}
