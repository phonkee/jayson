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
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
)

// APIObject gives ability to unmarshall json object into multiple fields
func APIObject(t require.TestingT, kv ...any) any {
	if len(kv)%2 != 0 {
		require.Fail(t, "APIObject: odd number of arguments")
		return nil
	}

	m := make(map[string]any)

	// iterate over key-value pairs
	for i := 0; i < len(kv); i += 2 {
		k, ok := kv[i].(string)
		assert.Truef(t, ok, "APIObject: key `%s` is not a string", kv[i])
		val := reflect.ValueOf(kv[i+1])
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		require.Truef(t, val.CanSet(), "APIObject: value `%s` is not settable", k)
		m[k] = kv[i+1]
	}

	return &object{
		kv: m,
		T:  t,
	}
}

// object is a helper struct for unmarshalling json object into multiple fields
type object struct {
	kv map[string]any
	T  require.TestingT
}

func (o *object) UnmarshalJSON(bytes []byte) error {
	into := make(map[string]json.RawMessage)
	err := json.Unmarshal(bytes, &into)
	require.NoError(o.T, err)

	for k, v := range o.kv {
		raw, ok := into[k]
		require.Truef(o.T, ok, "APIObject: key `%s` not found", k)

		err := json.Unmarshal(raw, v)
		require.NoError(o.T, err)
	}

	return nil
}
