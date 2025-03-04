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
	"reflect"
)

// extWithResponseTypes returns an extension that sets the response type.
func extWithResponseTypes(extension Extension, typ ...reflect.Type) Extension {
	return &extResponseType{
		ext:   extension,
		types: typ,
	}
}

// extResponseType is an extension that provides the response type.
type extResponseType struct {
	ext   Extension
	types []reflect.Type
}

// responseTypeIdentify returns the response type.
func (e *extResponseType) responseTypes() []reflect.Type { return e.types }

// ExtendResponseWriter extends the response writer.
func (e *extResponseType) ExtendResponseWriter(ctx context.Context, writer http.ResponseWriter) bool {
	return e.ext.ExtendResponseWriter(ctx, writer)
}

// ExtendResponseObject extends the response object.
func (e *extResponseType) ExtendResponseObject(ctx context.Context, m map[string]any) bool {
	return e.ext.ExtendResponseObject(ctx, m)
}
