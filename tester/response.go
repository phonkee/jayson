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
	"sort"
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

// operation is the type of operation for json path
type operation int

const (
	operationEquals operation = iota
	operationGt
	operationGte
	operationLt
	operationLte
)

// AssertJsonPath asserts that response body json path is equal to given value
func (r *response) AssertJsonPath(t require.TestingT, path string, what any) APIResponse {
	// set operation to equals
	op := operationEquals
	_ = op

	// check if path contains operation
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
		switch part {
		case "__len__":
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
		case "__keys__":
			var obj map[string]json.RawMessage
			require.NoErrorf(t, json.NewDecoder(bytes.NewReader(raw)).Decode(&obj), "failed to unmarshal `%v` into `%T`", path, what)

			keys := make([]string, 0, len(obj))
			for k := range obj {
				keys = append(keys, k)
			}

			original, ok := what.([]string)
			require.Truef(t, ok, "expected `[]string`, got: %T: %v", what)

			// now we need to sort the keys so we can compare them
			sort.Strings(original)
			sort.Strings(keys)

			require.Equalf(t, what, keys, "keys are not equal")
			return r
		case "__exists__":
			return r
		case "__gt__":
			op = operationGt
			break main
		case "__gte__":
			op = operationGte
			break main
		case "__lt__":
			op = operationLt
			break main
		case "__lte__":
			op = operationLte
			break main
		default:
			// pass
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

	switch op {
	case operationGt:
		assert.Greaterf(t, target, what, "value `%v` is not greater than `%v`", target, what)
	case operationGte:
		assert.GreaterOrEqualf(t, target, what, "value `%v` is not greater than or equal `%v`", target, what)
	case operationLt:
		assert.Lessf(t, target, what, "value `%v` is not less than `%v`", target, what)
	case operationLte:
		assert.LessOrEqualf(t, target, what, "value `%v` is not less than or equal `%v`", target, what)
	default:
		assert.Equalf(t, what, target, "expected: `%v`, got: `%v`", what, target)
	}

	return r
}

var (
	// jsonRawMessageType is the type of json.RawMessage
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
