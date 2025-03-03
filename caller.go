/*
 * MIT License
 *
 * Copyright (c) 2025 Peter Vrba
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

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
