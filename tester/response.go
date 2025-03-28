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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/phonkee/jayson/tester/action"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
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

// Print prints whole response
func (r *response) Print(writer ...io.Writer) APIResponse {
	var w io.Writer
	if len(writer) > 0 && writer[0] != nil {
		w = writer[0]
	} else {
		w = os.Stdout
	}

	// printf is a helper function to print formatted output and ignore errors
	printf := func(format string, args ...any) {
		_, _ = fmt.Fprintf(w, format+"\n", args...)
	}
	printf("Print:")
	printf("  Request:")
	printf("    URL: %v", r.request.URL.String())
	printf("    Method: %v", r.request.Method)
	printf("    Header: %v", r.request.Header)
	printf("  Response:")
	printf("    Body: %v", string(r.body))
	printf("    Status: %v", r.rw.Code)
	printf("    StatusMessage: %v", http.StatusText(r.rw.Code))
	printf("    Header: %v", r.rw.Header())
	printf("    Size: %v", r.rw.Body.Len())

	return r
}

// Header runs action against given header key
func (r *response) Header(t require.TestingT, key string, a action.Action) APIResponse {
	assert.Truef(t, a.Supports(action.SupportHeader), "action %v is not supported in Header call", a)

	var (
		err error
		raw json.RawMessage = nil
	)

	// check if header is present and get get action and marshal to json (we use json code and it's easier)
	if _, ok := r.rw.Result().Header[key]; ok {
		value := r.rw.Result().Header.Get(key)
		var buffer bytes.Buffer
		if err = json.NewEncoder(&buffer).Encode(value); err == nil {
			raw = buffer.Bytes()
		}
	} else {
		err = action.ErrNotPresent
	}

	// context instance
	ctx := context.WithValue(context.Background(), action.ContextKeyUnmarshalActionValue, r.unmarshalActionValue)

	var value any

	// if we don't have any error, we can try to unmarshal the action
	if err == nil {
		value, err = r.unmarshalActionValue(t, raw, a)
	}

	// check if unmarshal was not applied
	if errors.Is(err, action.ErrActionNotApplied) {
		err = nil
	}

	// now run action
	if err := a.Run(t, ctx, value, raw, err); err != nil {
		// TODO: check if we have action error, otherwise we need to add it somehow
		require.Fail(t, fmt.Sprintf("FAILED: `response.Header`, name: `%s`, %v", key, err.Error()))
	}

	return r
}

// Json runs given action against json path
func (r *response) Json(t require.TestingT, path string, a action.Action) APIResponse {
	assert.Truef(t, a.Supports(action.SupportJson), "action %v is not supported in Json call", a)

	raw := json.RawMessage(r.body)
	path = strings.TrimSpace(path)

	var (
		err error
		ok  bool
	)

main:
	for _, part := range strings.Split(path, ".") {
		if part == "" {
			continue main
		}
		// try to parse the part as a number so we know we need to unmarshal array
		if number, numberError := strconv.ParseUint(part, 10, 64); numberError == nil {
			var arr []json.RawMessage

			// try to unmarshal array
			if err = json.Unmarshal(raw, &arr); err != nil {
				err = fmt.Errorf("cannot unmarshal object into slice")
				break main
			}

			// check if the number is in bounds
			if int(number) >= len(arr) {
				err = fmt.Errorf("%w: out of bounds", action.ErrNotPresent)
				break main
			}

			raw = arr[number]
			continue main
		}
		// parse object
		obj := make(map[string]json.RawMessage)

		// try to unmarshal object, if it fails there is a problem
		if err = json.Unmarshal(raw, &obj); err != nil {
			err = fmt.Errorf("cannot unmarshal action")
			break main
		}

		// check if object contains the part
		if raw, ok = obj[part]; !ok {
			err = action.ErrNotPresent
			break main
		}
	}

	// context instance
	ctx := context.WithValue(context.Background(), action.ContextKeyUnmarshalActionValue, r.unmarshalActionValue)

	var value any

	// if we don't have any error, we can try to unmarshal the action
	if err == nil {
		value, err = r.unmarshalActionValue(t, raw, a)
	}

	// check if unmarshal was not applied
	if errors.Is(err, action.ErrActionNotApplied) {
		err = nil
	}

	// now run action
	if errRun := a.Run(t, ctx, value, raw, err); errRun != nil {
		// TODO: check if we have action error, otherwise we need to add it somehow
		require.Fail(t, fmt.Sprintf("FAILED: `response.Json`, path: `%s`, %v", path, errRun.Error()))
	}

	return r
}

// Status runs given action against status
func (r *response) Status(t require.TestingT, a action.Action) APIResponse {
	assert.Truef(t, a.Supports(action.SupportStatus), "action %v is not supported in Status call", a)

	// prepare raw action (json)
	raw := json.RawMessage(strconv.FormatInt(int64(r.rw.Code), 10))

	// unmarshal action
	value, err := r.unmarshalActionValue(t, raw, a)

	// check if unmarshal was not applied
	if errors.Is(err, action.ErrActionNotApplied) {
		err = nil
	}

	// context instance
	ctx := context.WithValue(context.Background(), action.ContextKeyUnmarshalActionValue, r.unmarshalActionValue)

	// now call a.Run
	if err := a.Run(t, ctx, value, raw, err); err != nil {
		// TODO: check if we have action error, otherwise we need to add it somehow
		require.NoError(t, err)
	}

	return r
}

// unmarshal given message action new action of given type
func (r *response) unmarshalActionValue(t require.TestingT, raw json.RawMessage, a action.Action) (any, error) {
	// check if we have action provided
	if v, ok := a.Value(t); ok {
		var value any
		var val reflect.Value
		typ := reflect.TypeOf(v)
		if typ.Kind() == reflect.Ptr {
			val = reflect.New(typ.Elem())
		} else {
			val = reflect.New(typ)
		}
		// prepare action
		value = val.Interface()

		// unmarshal action
		if err := json.NewDecoder(bytes.NewBuffer(raw)).Decode(value); err != nil {
			return nil, fmt.Errorf("%w: %s", action.ErrUnmarshal, err)
		}

		if typ.Kind() != reflect.Ptr {
			value = reflect.ValueOf(value).Elem().Interface()
		}

		return value, nil
	}

	return nil, action.ErrActionNotApplied
}
