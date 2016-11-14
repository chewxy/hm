// +build !debug

package hm

func tc() int                                   { return 0 }
func enterLoggingContext()                      {}
func leaveLoggingContext()                      {}
func logf(format string, others ...interface{}) {}
