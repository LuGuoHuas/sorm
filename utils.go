package sorm

import (
	"fmt"
	"regexp"
	"runtime"
)

var goSrcRegexp = regexp.MustCompile(`LuGuoHuas/sorm(@.*)?/.*.go`)
var goTestRegexp = regexp.MustCompile(`LuGuoHuas/sorm(@.*)?/.*test.go`)

func fileWithLineNum() string {
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && (!goSrcRegexp.MatchString(file) || goTestRegexp.MatchString(file)) {
			return fmt.Sprintf("%v:%v", file, line)
		}
	}
	return ""
}
