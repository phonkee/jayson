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
	"net/http"
)

// DefaultSettings returns default settings for jayson instance
func DefaultSettings() Settings {
	return Settings{
		DefaultErrorStatus:        http.StatusInternalServerError,
		DefaultErrorMessageKey:    "message",
		DefaultErrorStatusCodeKey: "code",
		DefaultErrorStatusTextKey: "status",
		DefaultResponseStatus:     http.StatusOK,
	}
}

// Settings for jayson instance
type Settings struct {
	DefaultErrorStatus        int
	DefaultErrorMessageKey    string
	DefaultErrorStatusCodeKey string
	DefaultErrorStatusTextKey string
	DefaultResponseStatus     int
}

func (s *Settings) Validate() {
	if s.DefaultErrorStatus == 0 {
		s.DefaultErrorStatus = http.StatusInternalServerError
	}
	if s.DefaultErrorMessageKey == "" {
		s.DefaultErrorMessageKey = "message"
	}
	if s.DefaultErrorStatusCodeKey == "" {
		s.DefaultErrorStatusCodeKey = "code"
	}
	if s.DefaultErrorStatusTextKey == "" {
		s.DefaultErrorStatusTextKey = "status"
	}
	if s.DefaultResponseStatus == 0 {
		s.DefaultResponseStatus = http.StatusOK
	}
}
