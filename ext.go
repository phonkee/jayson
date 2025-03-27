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
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

// ExtChain returns an extFunc that chains multiple ext together.
func ExtChain(extensions ...Extension) Extension {
	return ExtFunc(
		func(ctx context.Context, w http.ResponseWriter) (result bool) {
			for _, e := range extensions {
				if e.ExtendResponseWriter(ctx, w) {
					result = true
				}
			}
			return result
		},
		func(ctx context.Context, m map[string]any) (result bool) {
			for _, e := range extensions {
				if e.ExtendResponseObject(ctx, m) {
					result = true
				}
			}
			return result
		},
	)
}

// ExtConditional calls first extFunc, and if it returns true, it calls all the ext.
// This is useful for conditional ext based on context (such as debug mode or any context values).
func ExtConditional(condition Extension, ext ...Extension) Extension {
	return ExtFunc(
		func(ctx context.Context, w http.ResponseWriter) bool {
			var result bool
			if condition.ExtendResponseWriter(ctx, w) {
				for _, e := range ext {
					if e.ExtendResponseWriter(ctx, w) {
						result = true
					}
				}
			}
			return result
		},
		func(ctx context.Context, m map[string]any) bool {
			var result bool
			if condition.ExtendResponseObject(ctx, m) {
				for _, e := range ext {
					if e.ExtendResponseObject(ctx, m) {
						result = true
					}
				}
			}
			return result
		},
	)
}

// ExtFirst returns an extFunc that returns the first extFunc that extends the response.
func ExtFirst(ext ...Extension) Extension {
	return ExtFunc(
		func(ctx context.Context, w http.ResponseWriter) bool {
			for _, e := range ext {
				if e.ExtendResponseWriter(ctx, w) {
					return true
				}
			}
			return false
		},
		func(ctx context.Context, m map[string]any) bool {
			for _, e := range ext {
				if e.ExtendResponseObject(ctx, m) {
					return true
				}
			}
			return false
		},
	)
}

// ExtFunc is an extFunc that calls the given function to extend the response or the response object.
func ExtFunc(
	extResponseWriter ExtensionApplyResponseWriterFunc,
	extResponseObject ExtensionApplyResponseObjectFunc,
) Extension {
	return &extFunc{
		extendResponseWriter: func(ctx context.Context, writer http.ResponseWriter) bool {
			if extResponseWriter != nil {
				return extResponseWriter(ctx, writer)
			}
			return false
		},
		extendResponseObject: func(ctx context.Context, m map[string]any) bool {
			if extResponseObject != nil {
				return extResponseObject(ctx, m)
			}
			return false
		},
	}
}

// ExtHeader extends the response with the given HTTP headers.
func ExtHeader(h http.Header) Extension {
	return ExtFunc(
		func(ctx context.Context, w http.ResponseWriter) bool {
			if h == nil {
				return false
			}
			for k, v := range h {
				for _, vv := range v {
					w.Header().Add(k, vv)
				}
			}
			return true
		},
		nil,
	)
}

// ExtHeaderValue is an extFunc that adds a single HTTP header to the response.
func ExtHeaderValue(k string, v string) Extension {
	return ExtFunc(
		func(ctx context.Context, w http.ResponseWriter) bool {
			w.Header().Add(k, v)
			return true
		},
		nil,
	)
}

// ExtNoop does not extend the response or the response object.
func ExtNoop() Extension {
	return ExtFunc(nil, nil)
}

// ExtObjectKeyValue is an extFunc that adds a single key-value pair to the response object.
func ExtObjectKeyValue(key string, value any) Extension {
	return extWithResponseTypes(
		ExtFunc(
			nil,
			func(ctx context.Context, m map[string]any) bool {
				m[key] = value
				return true
			},
		),
		reflect.TypeOf(value),
	)
}

// ExtObjectKeyValuef is an extFunc that adds a single key-value pair to the response object based on a format string.
func ExtObjectKeyValuef(key string, format string, args ...any) Extension {
	var types []reflect.Type

	for _, arg := range args {
		types = append(types, reflect.TypeOf(arg))
	}

	return extWithResponseTypes(
		ExtFunc(
			nil,
			func(ctx context.Context, m map[string]any) bool {
				m[key] = fmt.Sprintf(format, args...)
				return true
			},
		),
		types...,
	)
}

var (
	// jsonRawMessageType is the type of json.RawMessage
	jsonRawMessageType = reflect.TypeOf(json.RawMessage{})
)

// ExtObjectUnwrap is an Extension that converts the given object to the response object.
// It is useful if you want to add key/values to the response object (by altering it via Extensions).
// If the object is a struct, it will be converted to a map[string]any.
func ExtObjectUnwrap(obj any) Extension {
	// unmarshal object to map[string]any
	objMap := make(map[string]any)

	// try to inspect struct/map type and add it to the response object
	val := reflect.ValueOf(obj)

	// if object is a pointer, dereference it
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	switch val.Kind() {
	case reflect.Struct:
		structToMap(val, objMap)
	case reflect.Map:
		for key, value := range val.MapKeys() {
			objMap[fmt.Sprintf("%v", key)] = val.MapIndex(value).Interface()
		}
	default:
		// check some types here
		switch val.Type() {
		case jsonRawMessageType: // json.RawMessage can be deconstructed, so we parse it to a map[string]json.RawMessage
			// custom handling
			var into map[string]json.RawMessage
			if err := json.Unmarshal(val.Bytes(), &into); err == nil {
				for k, v := range into {
					objMap[k] = v
				}
			} else {
				objMap = nil
			}
		default:
			objMap = nil
		}
	}

	return extWithResponseTypes(
		ExtFunc(
			nil,
			func(ctx context.Context, m map[string]any) bool {
				if objMap == nil {
					s := ContextSettingsValue(ctx)
					m[s.DefaultUnwrapObjectKey] = obj
				} else {
					for k, v := range objMap {
						m[k] = v
					}
				}
				return true
			},
		),
		reflect.TypeOf(obj),
	)
}

// isEmptyValue checks if a value is empty (zero value).
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	default:
		// no-op
	}
	return false
}

// parseJSONTag parses json tag and returns field name, omitempty and skip values.
func parseJSONTag(tag string, name string) (fieldName string, omitempty bool, skip bool) {
	// Check for "omitempty" in the tag
	tagParts := strings.Split(tag, ",")
	fieldName = tagParts[0]
	if fieldName == "" {
		fieldName = name
	}
	omitempty = len(tagParts) > 1 && tagParts[1] == "omitempty"
	skip = len(tagParts) > 0 && tagParts[0] == "-"
	return fieldName, omitempty, skip
}

// structToMap converts a struct to a map[string]any.
func structToMap(val reflect.Value, into map[string]any) {
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		// skip unexported fields
		if field.PkgPath != "" {
			if !field.Anonymous {
				continue
			}
			structToMap(val.Field(i), into)
			continue
		}

		// parse json tag
		fieldName, omitEmpty, skip := parseJSONTag(field.Tag.Get("json"), field.Name)
		if skip {
			continue
		}
		if omitEmpty && isEmptyValue(val.Field(i)) {
			continue
		}

		into[fieldName] = val.Field(i).Interface()
	}
}

// ExtOmitObjectKey is an extFunc that removes the given keys from the response object.
// This extFunc needs to be called at the end,
// so it removes the keys after all other ext have added their keys.
func ExtOmitObjectKey(keys ...string) Extension {
	return ExtFunc(
		nil,
		func(ctx context.Context, m map[string]any) (result bool) {
			for _, key := range keys {
				if _, ok := m[key]; ok {
					delete(m, key)
					result = true
				}
			}
			return result
		},
	)
}

// ExtStatus is an extension that sets the HTTP status code of the response.
func ExtStatus(status int) Extension {
	return ExtFunc(
		func(ctx context.Context, w http.ResponseWriter) bool {
			w.WriteHeader(status)
			return true
		},
		nil,
	)
}

// extSettingsKeyValue is an extFunc that adds a single key-value pair to the response object based on the settings.
func extSettingsKeyValue(fn func(s Settings) string, value any) Extension {
	return ExtFunc(
		nil,
		func(ctx context.Context, m map[string]any) bool {
			s := ContextSettingsValue(ctx)
			if key := fn(s); key != "" {
				m[key] = value
			}
			return true
		},
	)
}

// ExtOmitSettingsKey is an extFunc that removes the given keys from the response object based on the settings.
func ExtOmitSettingsKey(fn func(settings Settings) []string) Extension {
	return ExtFunc(
		nil,
		func(ctx context.Context, m map[string]any) (result bool) {
			s := ContextSettingsValue(ctx)
			for _, key := range fn(s) {
				if _, ok := m[key]; ok {
					delete(m, key)
					result = true
				}
			}
			return result
		},
	)
}

// extFunc is a generic extFunc that can extend the response or the response object.
// it is used to implement the Extension interface.
type extFunc struct {
	extendResponseWriter func(context.Context, http.ResponseWriter) bool
	extendResponseObject func(context.Context, map[string]any) bool
}

// ExtendResponseWriter calls the extFunc's extendResponseWriter function.
func (e *extFunc) ExtendResponseWriter(ctx context.Context, w http.ResponseWriter) bool {
	return e.extendResponseWriter(ctx, w)
}

// ExtendResponseObject calls the extFunc's extendResponseObject function.
func (e *extFunc) ExtendResponseObject(ctx context.Context, m map[string]any) bool {
	return e.extendResponseObject(ctx, m)
}
