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
