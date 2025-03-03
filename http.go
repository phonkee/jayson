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
	"bytes"
	"net/http"
)

// applyHeader applies headers from src to dst
func applyHeader(dst http.Header, src http.Header) {
	for k, v := range src {
		for _, vv := range v {
			dst.Add(k, vv)
		}
	}
}

// newResponseWriter creates a new response writer
// this is used for compatibility reasons and for various quirks of the http.ResponseWriter
func newResponseWriter(statusCode int) *responseWriter {
	return &responseWriter{
		header:     make(http.Header),
		buffer:     bytes.Buffer{},
		statusCode: statusCode,
	}
}

type responseWriter struct {
	header     http.Header
	buffer     bytes.Buffer
	statusCode int
}

// Header returns the header map
func (r *responseWriter) Header() http.Header {
	return r.header
}

// Write writes the data to the connection as part of an HTTP reply.
func (r *responseWriter) Write(bytes []byte) (int, error) {
	return r.buffer.Write(bytes)
}

// WriteHeader writes status code to writer
func (r *responseWriter) WriteHeader(statusCode int) {
	r.statusCode = statusCode
}

// WriteTo writes the response to the given http.ResponseWriter
func (r *responseWriter) WriteTo(w http.ResponseWriter) {
	applyHeader(w.Header(), r.header)
	w.WriteHeader(r.statusCode)
	_, _ = r.buffer.WriteTo(w)
}
