package jayson

import (
	"context"
	"net/http"
)

// newExecutor creates a new executor
func newExecutor(exts ...[]Extension) *executor {
	return &executor{
		extensions: exts,
	}
}

type executor struct {
	extensions [][]Extension
}

// ExtendResponseObject extends the object
func (e *executor) ExtendResponseObject(ctx context.Context, obj map[string]any) (result bool) {
	for _, exts := range e.extensions {
		for _, ext := range exts {
			if ext.ExtendResponseObject(ctx, obj) {
				result = true
			}
		}
	}
	return result
}

// ExtendResponseWriter extends the response
func (e *executor) ExtendResponseWriter(ctx context.Context, w http.ResponseWriter) (result bool) {
	for _, exts := range e.extensions {
		for _, ext := range exts {
			if ext.ExtendResponseWriter(ctx, w) {
				result = true
			}
		}
	}
	return result
}
