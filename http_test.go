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
	"net/http"
	"testing"
)

func TestResponseWriter(t *testing.T) {
	rw := newResponseWriter(http.StatusOK)
	assert.Equal(t, http.StatusOK, rw.statusCode)

	rw.WriteHeader(http.StatusTeapot)
	assert.Equal(t, http.StatusTeapot, rw.statusCode)

	what := []byte("hello")

	n, err := rw.Write(what)
	assert.Nil(t, err)
	assert.Equal(t, len(what), n)
	assert.Equal(t, "hello", rw.buffer.String())

	rw.Header().Set("Content-Type", "text/plain")
	assert.Equal(t, "text/plain", rw.Header().Get("Content-Type"))

	newRw := newResponseWriter(http.StatusOK)
	rw.WriteTo(newRw)

	assert.Equal(t, http.StatusTeapot, newRw.statusCode)
	assert.Equal(t, "hello", newRw.buffer.String())
	assert.Equal(t, "text/plain", newRw.Header().Get("Content-Type"))

}
