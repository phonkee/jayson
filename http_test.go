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
