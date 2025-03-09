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

package tester

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
)

// newResponse creates a new *response instance
func newResponse(rw *httptest.ResponseRecorder, request *http.Request) *response {
	return &response{
		rw:      rw,
		request: request,
	}
}

// response is the implementation of APIResponse
type response struct {
	rw      *httptest.ResponseRecorder
	request *http.Request
	body    []byte
}

// AssertHeaderValue asserts that response header value is equal to given value
func (r *response) AssertHeaderValue(t require.TestingT, key, value string) APIResponse {
	for _, v := range r.rw.Result().Header[key] {
		if v == value {
			return r
		}
	}
	require.Failf(t, "fail", "header `%s` not found or does not have value", key)
	return r
}

// AssertJsonEquals asserts that response body is equal to given json string or object/map
func (r *response) AssertJsonEquals(t require.TestingT, expected any) APIResponse {
	require.NotNilf(t, r.body, "response body is nil")

	var expectedStr string

	switch expected.(type) {
	case string:
		expectedStr = expected.(string)
	case []byte:
		expectedStr = string(expected.([]byte))
	default:
		b, err := json.Marshal(expected)
		require.NoErrorf(t, err, "failed to marshal expectedError value: %v", expected)
		expectedStr = string(b)
	}

	require.JSONEq(t, expectedStr, string(r.body))
	return r
}

// AssertJsonKeyEquals asserts that response body key is equal to given value
// This method uses bit of magic to unmarshal json object into given value.
// It inspects the type of the given value and unmarshalls the json object into same type.
func (r *response) AssertJsonKeyEquals(t require.TestingT, key string, what any) APIResponse {
	require.NotNilf(t, r.body, "response body is nil")

	var val reflect.Value
	typ := reflect.TypeOf(what)
	if typ.Kind() == reflect.Ptr {
		val = reflect.New(typ.Elem())
	} else {
		val = reflect.New(typ).Elem()
	}

	// prepare object
	obj := make(map[string]json.RawMessage)

	assert.NoError(t, json.NewDecoder(bytes.NewReader(r.body)).Decode(&obj))

	v, ok := obj[key]
	require.Truef(t, ok, "key %s not found in response", key)

	target := val.Interface()
	assert.NoError(t, json.NewDecoder(bytes.NewBuffer(v)).Decode(&target))

	assert.Equalf(t, what, target, "expectedError: %v, got: %v", what, target)

	return r
}

// AssertJsonPathEquals asserts that response body json path is equal to given value
func (r *response) AssertJsonPathEquals(t require.TestingT, path string, what any) APIResponse {
	path = strings.TrimSpace(path)
	require.NotZerof(t, path, "path is empty")

	// first we get the splitted path
	splitted := strings.Split(path, ".")

	// prepare response body as raw message
	var raw = json.RawMessage(r.body)

	// iterate over all parts of the path
main:
	for _, part := range splitted {
		if part == "" {
			continue
		}
		// special case for __len__ which is used to get the length of the array or object
		if part == "__len__" {
			var array []json.RawMessage

			// try to unmarshal into array first
			if err := json.NewDecoder(bytes.NewReader(raw)).Decode(&array); err == nil {
				raw = json.RawMessage(strconv.Itoa(len(array)))
			} else {
				var obj map[string]json.RawMessage
				if err := json.NewDecoder(bytes.NewReader(raw)).Decode(&obj); err == nil {
					raw = json.RawMessage(strconv.Itoa(len(obj)))
				} else {
					require.Failf(t, "failed to get length", "failed to get length of `%v`", path)
				}
			}
			break main
		}

		// try to parse the part as a number so we know we need to unmarshal array
		if number, err := strconv.ParseUint(part, 10, 64); err == nil {
			var arr []json.RawMessage
			require.NoErrorf(t, json.Unmarshal(raw, &arr), "failed to unmarshal array `%v` into `%T`", path, what)
			require.Lessf(t, int(number), len(arr), "index out of bounds: %d", number)
			raw = arr[number]
			continue main
		}

		// parse object
		obj := make(map[string]json.RawMessage)

		// unmarshal object
		require.NoErrorf(t, json.Unmarshal(raw, &obj), "failed to unmarshal `%v` into `%T`", path, what)

		require.Containsf(t, obj, part, "key `%s` in path `%s` not found", part, path)

		raw = obj[part]
	}

	// now we have the value in raw, we can unmarshal it into given value
	var val reflect.Value
	typ := reflect.TypeOf(what)
	if typ.Kind() == reflect.Ptr {
		val = reflect.New(typ.Elem())
	} else {
		val = reflect.New(typ)
	}

	// prepare value
	target := val.Interface()
	require.NoErrorf(t, json.NewDecoder(bytes.NewBuffer(raw)).Decode(target), "failed to unmarshal `%v` into `%T`", path, what)

	// when not pointer, we need to get the value
	if typ.Kind() != reflect.Ptr {
		target = val.Elem().Interface()
	}

	// in case of json.RawMessage we call JSONEqf so it's compared as string in the correct order
	if typ == jsonRawMessageType {
		require.JSONEqf(t, string(what.(json.RawMessage)), string(target.(json.RawMessage)), "expectedError: %v, got: %v", what, target)
		return r
	}

	assert.Equalf(t, what, target, "expectedError: %v, got: %v", what, target)

	return r
}

var (
	jsonRawMessageType = reflect.TypeOf(json.RawMessage(nil))
)

// AssertStatus asserts that response status is equal to given status
func (r *response) AssertStatus(t require.TestingT, status int) APIResponse {
	require.Equal(t, status, r.rw.Code)
	return r
}

// Unmarshal unmarshalls whole response body into given value
func (r *response) Unmarshal(t require.TestingT, v any) APIResponse {
	require.NotNilf(t, r.body, "response body is nil")
	require.NoError(t, json.NewDecoder(bytes.NewReader(r.body)).Decode(v))
	return r
}
