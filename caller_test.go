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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCallerInfo(t *testing.T) {
	t.Run("test getCallerInfo", func(t *testing.T) {
		ci := getCallerInfo(100)
		assert.NotEmpty(t, ci.file)
		assert.NotEmpty(t, ci.fn)
		assert.NotZero(t, ci.line)
	})

	t.Run("test getCallerInfo with skip", func(t *testing.T) {
		ci := getCallerInfo(100, 3)
		assert.Equal(t, defaultUnknownFile, ci.file)
		assert.Equal(t, defaultUnknownFn, ci.fn)
		assert.Zero(t, ci.line)
	})
}
