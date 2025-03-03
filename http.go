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
