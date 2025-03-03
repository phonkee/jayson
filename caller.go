package jayson

import (
	"regexp"
	"runtime"
)

var (
	// patternJaysonPackage matches jayson package
	patternJaysonPackage = regexp.MustCompile(`github.com/phonkee/jayson\..+`)
)

// newCallerInfo returns new caller info that is outside jayson package
func newCallerInfo(maxDepth int) callerInfo {
	for i := 1; i < maxDepth; i++ {
		pc, file, no, ok := runtime.Caller(i)
		if !ok {
			break
		}
		funcName := runtime.FuncForPC(pc).Name()
		if !patternJaysonPackage.MatchString(funcName) {
			return callerInfo{
				file: file,
				fn:   funcName,
				line: no,
			}
		}
	}
	return callerInfo{
		file: "<unknown>",
		fn:   "<unknown>",
		line: 0,
	}
}

// callerInfo holds caller info
type callerInfo struct {
	file string
	fn   string
	line int
}
