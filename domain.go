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
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"reflect"
)

// Extension is the interface for Jayson ext.
//
// It is used to alter the response or the object before it is sent to the client.
//
//go:generate mockery --name Extension --filename extension.go
type Extension interface {
	// ExtendResponseWriter extends the response
	ExtendResponseWriter(context.Context, http.ResponseWriter) bool
	// ExtendResponseObject extends the object
	ExtendResponseObject(context.Context, map[string]any) bool
}

// ExtensionApplyResponseWriterFunc is a function that extends the response writer.
type ExtensionApplyResponseWriterFunc func(context.Context, http.ResponseWriter) bool

// ExtensionApplyResponseObjectFunc is a function that extends the response object.
type ExtensionApplyResponseObjectFunc func(context.Context, map[string]any) bool

// Jayson is the interface that provides methods for writing responses to the client.
//
// By default, there is a global instance to be used,
// but for some special purposes you can create your own (multiple servers, different settings, etc.)
type Jayson interface {
	// Debug enables debug mode via zap logger.
	Debug(*zap.Logger)
	// Error writes error to the client.
	Error(context.Context, http.ResponseWriter, error, ...Extension)
	// RegisterError registers extFunc for given error.
	RegisterError(error, ...Extension) error
	// RegisterResponse registers extFunc for given response object.
	RegisterResponse(any, ...Extension) error
	// Response writes given object/error to the client.
	Response(context.Context, http.ResponseWriter, any, ...Extension)
}

var (
	Warning = errors.New("jayson: warning")
	// ErrImproperlyConfigured is error returned when Jayson is improperly configured.
	ErrImproperlyConfigured = errors.New("jayson: improperly configured")
	WarnAlreadyRegistered   = fmt.Errorf("%w: already registered", Warning)
)

const (
	// DebugMaxCallerDepth is the maximum depth of debug caller
	DebugMaxCallerDepth = 255
)

// responseTypes is an internal interface that identifies the response types.
// It is used in ExtObjectUnwrap to identify the type of the object.
type responseTypes interface {
	// responseTypes returns the type of the type passed
	responseTypes() []reflect.Type
}

// Extended is interface for errors that can be extended.
type Extended interface {
	// Extensions returns list of extensions for the error.
	Extensions() []Extension
}
