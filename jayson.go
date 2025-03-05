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
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"reflect"
)

// New instantiates custom jayson instance. Usually you don't need to use it, since there is a _g instance.
func New(settings Settings) Jayson {
	// validate settings
	settings.Validate()

	return &jayson{
		settings:              settings,
		registryErrors:        newRegistry[error](),
		registryResponseTypes: newRegistry[reflect.Type](),
	}
}

// jayson implements Jayson interface
type jayson struct {
	debug    *zap.Logger
	settings Settings

	// registry for errors
	registryErrors *registry[error]
	// registry for response types
	registryResponseTypes *registry[reflect.Type]
}

// Debug enables debug mode via zap logger
func (j *jayson) Debug(logger *zap.Logger) {
	j.debug = logger
}

// Error writes error response to the client
func (j *jayson) Error(ctx context.Context, rw http.ResponseWriter, err error, override ...Extension) {
	if err == nil {
		return
	}

	// get error extensions
	ext, _ := j.getErrorExtensions(err, override...)

	// prepare context
	ctx = contextWithErrorValue(
		contextWithSettingsValue(ctx, j.settings),
		err,
	)

	// rwInternal
	obj := map[string]any{
		j.settings.DefaultErrorMessageKey: err.Error(),
	}

	// prepare internal response writer
	rwInternal := newResponseWriter(j.settings.DefaultErrorStatus)

	// prepare executor
	exec := newExecutor(ext)

	// now extend response
	exec.ExtendResponseWriter(ctx, rwInternal)

	// now add additional properties to object
	// handle status code and text
	obj[j.settings.DefaultErrorStatusCodeKey] = rwInternal.statusCode

	// handle status text
	if text := http.StatusText(rwInternal.statusCode); text != "" {
		obj[j.settings.DefaultErrorStatusTextKey] = text
	}

	// now extend object
	exec.ExtendResponseObject(ctx, obj)

	// clear buffer here
	rwInternal.buffer = bytes.Buffer{}
	rwInternal.Header()["Content-Type"] = []string{"application/json"}

	// now write JSON value
	if err := json.NewEncoder(rwInternal).Encode(obj); err != nil {
		// if we can't write JSON, we will panic
		// it's fine now
		// FIXME: we should log this
		panic(err)
	}

	// write to a response writer
	rwInternal.WriteTo(rw)
}

// RegisterError registers extFunc for given error
func (j *jayson) RegisterError(err error, ext ...Extension) error {
	j.debugLogMethod("RegisterError", func() []zap.Field {
		return []zap.Field{
			zap.String("type", reflect.TypeOf(err).String()),
			zap.Int("ext", len(ext)),
		}
	})

	// if err is nil, this is error
	if err == nil {
		panic(fmt.Errorf("%w: error is nil", ErrImproperlyConfigured))
	}

	// check for Extended interface
	if extended, ok := err.(Extended); ok {
		ext = append(extended.Extensions(), ext...)
	}

	// if Any, we will Register ext for any error
	if errors.Is(err, Any) {
		j.registryErrors.AddShared(ext...)
		return nil
	}

	return j.registryErrors.Register(err, ext)
}

// RegisterResponse registers extFunc for given object
func (j *jayson) RegisterResponse(what any, extensions ...Extension) error {

	// log caller
	j.debugLogMethod("RegisterResponse", func() []zap.Field {
		return []zap.Field{
			zap.Any("type", reflect.TypeOf(what).String()),
			zap.Int("ext", len(extensions)),
		}
	})

	// if what is Any, we will Register ext for any response object
	if what == Any {
		j.registryResponseTypes.AddShared(extensions...)
		return nil
	}

	// register response type
	return j.registryResponseTypes.Register(reflect.TypeOf(what), extensions)
}

// Response writes response to the client
func (j *jayson) Response(ctx context.Context, rw http.ResponseWriter, what any, override ...Extension) {
	// add object value to the context along with settings
	ctx = contextWithObjectValue(
		contextWithSettingsValue(ctx, j.settings),
		what,
	)

	// rwInternal is a response writer that will be used to collect response
	rwInternal := newResponseWriter(j.settings.DefaultResponseStatus)

	// if what is an override, we will be having object automatically
	if extension, ok := what.(Extension); ok {
		j.responseExtension(ctx, rwInternal, what, extension, override...)
	} else {
		j.responseRaw(ctx, rwInternal, what, override...)
	}

	// set content type
	rwInternal.Header()["Content-Type"] = []string{"application/json"}
	rwInternal.WriteTo(rw)
}

// responseExtension is called when `what` is an extension
func (j *jayson) responseExtension(ctx context.Context, rw *responseWriter, what any, whatExt Extension, override ...Extension) {
	// create object
	obj := make(map[string]any)

	// this fixes the order of extensions
	override = append([]Extension{whatExt}, override...)

	// types of what so we can Get ext
	var types []reflect.Type

	// check if extension is responseTypes so we can Get type
	if rti, ok := whatExt.(responseTypes); ok {
		types = rti.responseTypes()
	}
	if len(types) == 0 {
		types = []reflect.Type{reflect.TypeOf(what)}
	}

	var (
		ext []Extension
		ok  bool
	)

	// now range over types and first that returns true is used
	for _, typ := range types {
		if ext, ok = j.getResponseTypeExtensions(typ, override...); ok {
			break
		}
	}

	// prepare executor
	exec := newExecutor(ext)

	// extend response writer
	exec.ExtendResponseWriter(ctx, rw)

	// extend object
	exec.ExtendResponseObject(ctx, obj)

	// now clear buffer if someone mistakenly wrote to it
	rw.buffer = bytes.Buffer{}

	// json marshal object
	if err := json.NewEncoder(rw).Encode(obj); err != nil {
		panic(err)
	}
}

// responseRaw is called when `what` is not an extension
func (j *jayson) responseRaw(ctx context.Context, rw *responseWriter, what any, override ...Extension) {
	ext, ok := j.getResponseTypeExtensions(reflect.TypeOf(what), override...)
	_ = ok
	// prepare executor with ext (first what we found in registered types, then what is passed in function)
	exec := newExecutor(ext)

	// now extend response, no object here
	exec.ExtendResponseWriter(ctx, rw)

	// now clear buffer if someone mistakenly wrote to it
	rw.buffer = bytes.Buffer{}

	// now json encode object
	if err := json.NewEncoder(rw).Encode(what); err != nil {
		panic(err)
	}
}

// getErrorExtensions returns all extensions for given error
// it also adds extensions for all parent errors and Any
// even when false is returned, extensions are returned
func (j *jayson) getErrorExtensions(err error, override ...Extension) ([]Extension, bool) {
	var (
		result []Extension
		found  bool
	)

	for err != nil {
		if extended, ok := err.(Extended); ok {
			result = append(extended.Extensions(), result...)
		}

		// we prepend errors
		if ext, ok := j.registryErrors.Get(err); ok {
			result = append(ext, result...)
			found = true
		}

		err = errors.Unwrap(err)
	}

	// add shared extensions
	result = j.registryErrors.WithShared(result...)
	// and append override
	result = append(result, override...)

	return result, found
}

// getResponseTypeExtensionsBare returns all extensions for given response type
// no other extensions are added (no default, no overrides
func (j *jayson) getResponseTypeExtensionsBare(what reflect.Type, level int) ([]Extension, bool) {
	var (
		ext []Extension
		ok  bool
	)

	// check exact type
	if ext, ok = j.registryResponseTypes.Get(what); !ok {
		switch what.Kind() {
		case reflect.Ptr:
			// check if we have extension for pointer type
			if ext, ok = j.registryResponseTypes.Get(what); !ok && level == 0 {

				// if we are on zero level, we will try to Get extension for Elem
				return j.getResponseTypeExtensionsBare(what.Elem(), level+1)
			}
		case reflect.Slice, reflect.Array:
			ext, ok = j.getResponseTypeExtensionsBare(what.Elem(), 0)
		default:
			if ext, ok = j.registryResponseTypes.Get(what); !ok && level == 0 {
				return j.getResponseTypeExtensionsBare(reflect.PointerTo(what), level+1)
			}
		}
	}

	return ext, ok
}

// getResponseTypeExtensions returns all extensions for given response type
// it also adds extensions for pointer types, slices
// even when false is returned, extensions are returned
func (j *jayson) getResponseTypeExtensions(what reflect.Type, override ...Extension) ([]Extension, bool) {
	var (
		ext []Extension
		ok  bool
	)

	// check exact type
	ext, ok = j.getResponseTypeExtensionsBare(what, 0)

	ext = j.registryResponseTypes.WithShared(ext...)

	ext = append(ext, override...)

	return ext, ok
}

// debugLogMethod logs caller info
func (j *jayson) debugLogMethod(method string, fn ...func() []zap.Field) {
	if j.debug == nil {
		return
	}

	// check if we are on debug level
	if ch := j.debug.Named("jayson").Check(zap.DebugLevel, "caller info"); ch != nil {
		// try to Get caller info
		caller := newCallerInfo(DebugMaxCallerDepth)

		fields := []zap.Field{
			zap.String("method", method),
			zap.String("function", caller.fn),
			zap.Int("line", caller.line),
			zap.String("file", caller.file),
		}

		for _, fetch := range fn {
			fields = append(fields, fetch()...)
		}

		ch.Write(fields...)
	}
}
