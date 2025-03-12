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

package action

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Unmarshal action unmarshal the given json response to the given value
func Unmarshal(into any) Action {
	return &actionFunc{
		run: func(t require.TestingT, value any, raw json.RawMessage, err error) error {
			if err != nil {
				return err
			}

			if err = json.NewDecoder(bytes.NewBuffer(raw)).Decode(into); err != nil {
				return fmt.Errorf("%w: cannot unmarshal response to %T", ErrActionUnmarshal, into)
			}

			return nil
		},
		support:   []Support{SupportJson},
		baseError: ErrActionUnmarshal,
	}
}

type KV map[string]any

// UnmarshalObjectKeys action unmarshal the given json response object keys to the given values
// kv should be pairs of string key and addressable value
func UnmarshalObjectKeys(keys KV) Action {
	return &actionFunc{
		run: func(t require.TestingT, value any, raw json.RawMessage, err error) error {
			if err != nil {
				return err
			}

			// prepare object to unmarshal
			obj := make(map[string]json.RawMessage)

			// unmarshal response to object
			// TODO: finish this
			assert.NoErrorf(t, json.NewDecoder(bytes.NewBuffer(raw)).Decode(&obj), "cannot unmarshal response to %T", obj)

			// iterate over all keys and unmarshal them to the given values
			for key, val := range keys {
				v, ok := obj[key]

				// check if key exists in the response
				assert.Truef(t, ok, "key %s not found in response", key)

				// unmarshal the value to the given addressable value
				assert.NoErrorf(t, json.NewDecoder(bytes.NewBuffer(v)).Decode(val), "cannot unmarshal response to %T", value)
			}

			return nil
		},
		support:   []Support{SupportJson},
		baseError: ErrActionUnmarshalObjectKeys,
	}
}
