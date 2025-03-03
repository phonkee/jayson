package jayson

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExecutor(t *testing.T) {
	var value int
	ext := []Extension{
		ExtFunc(
			func(ctx context.Context, writer http.ResponseWriter) bool {
				value += 42
				return true
			},
			func(ctx context.Context, m map[string]any) bool {
				value += 1
				return true
			},
		),
	}
	exec := newExecutor(ext, ext, ext, ext)
	exec.ExtendResponseWriter(context.Background(), httptest.NewRecorder())
	exec.ExtendResponseObject(context.Background(), map[string]any{})
	assert.Equal(t, 4*43, value)
}
