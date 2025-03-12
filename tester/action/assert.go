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
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/constraints"
	"regexp"
	"sort"
)

// AssertEquals asserts that given value is equal to the value in the response
func AssertEquals(value any) Action {
	return &actionFunc{
		value: func(t require.TestingT) (any, bool) {
			return value, true
		},
		run: func(t require.TestingT, ctx context.Context, v any, raw json.RawMessage, err error) error {
			if err != nil {
				return err
			}
			if !assert.ObjectsAreEqual(v, value) {
				return fmt.Errorf("%w, expected: %#v, got: %#v", ErrActionAssertEquals, value, v)
			}

			return nil
		},
		support:   []Support{SupportHeader, SupportJson, SupportStatus},
		baseError: ErrActionAssertEquals,
	}
}

// AssertExists asserts that the value exists in the response based on given argument
func AssertExists() Action {
	return &actionFunc{
		run: func(t require.TestingT, ctx context.Context, v any, raw json.RawMessage, err error) error {
			if err != nil {
				return err
			}
			if raw == nil {
				return fmt.Errorf("%w: expected value to exist, but it does not", ErrActionAssertExists)
			}
			return nil
		},
		support:   []Support{SupportHeader, SupportJson},
		baseError: ErrActionAssertExists,
	}
}

// AssertGt asserts that given value is greater than the value in the response
func AssertGt[T constraints.Integer](value T) Action {
	return &actionFunc{
		value: func(t require.TestingT) (any, bool) {
			return value, true
		},
		run: func(t require.TestingT, ctx context.Context, v any, raw json.RawMessage, err error) error {
			if err != nil {
				return err
			}
			if !assert.Greater(&requireTestingT{}, v, value) {
				return fmt.Errorf("%w: expected %v to be greater than %v", ErrActionAssertGt, value, v)
			}
			return nil
		},
		support:   []Support{SupportJson, SupportStatus},
		baseError: ErrActionAssertGt,
	}
}

// AssertGte asserts that given value is greater than or equal to the value in the response
func AssertGte[T constraints.Integer](value T) Action {
	return &actionFunc{
		value: func(t require.TestingT) (any, bool) {
			return value, true
		},
		run: func(t require.TestingT, ctx context.Context, v any, raw json.RawMessage, err error) error {
			if err != nil {
				return err
			}

			if !assert.GreaterOrEqual(&requireTestingT{}, v, value) {
				return fmt.Errorf("%w: expected %v to be greater than or equal to %v", ErrActionAssertGte, v, value)
			}
			return nil
		},
		support:   []Support{SupportJson, SupportStatus},
		baseError: ErrActionAssertGte,
	}
}

// AssertIn asserts that value is in the given list of values
func AssertIn[T any](values ...T) Action {
	return &actionFunc{
		value: func(t require.TestingT) (any, bool) {
			return values[0], true
		},
		run: func(t require.TestingT, ctx context.Context, v any, raw json.RawMessage, err error) error {
			if err != nil {
				return err
			}
			ok, found := containsElement(values, v)
			if !ok || !found {
				return fmt.Errorf("%w: expected value %#v to be in %#v", ErrActionAssertIn, v, values)
			}

			return nil
		},
		support:   []Support{SupportHeader, SupportJson, SupportStatus},
		baseError: ErrActionAssertIn,
	}
}

// AssertKeys asserts that given map has the given keys
func AssertKeys(keys ...string) Action {
	return &actionFunc{
		value: func(t require.TestingT) (any, bool) {
			return map[string]any{}, true
		},
		run: func(t require.TestingT, ctx context.Context, v any, raw json.RawMessage, err error) error {
			if err != nil {
				return err
			}
			// collect all the keys from the map
			found := make([]string, 0, len(keys))
			for key := range v.(map[string]any) {
				found = append(found, key)
			}
			// sort both so we can compare them
			sort.Strings(found)
			sort.Strings(keys)

			if !assert.ObjectsAreEqual(keys, found) {
				return fmt.Errorf("%w: expected keys %#v, got %#v", ErrActionAssertKeys, keys, found)
			}
			return nil
		},
		support:   []Support{SupportJson},
		baseError: ErrActionAssertKeys,
	}
}

// AssertLen asserts that given array/object has the given length
func AssertLen[T constraints.Integer](l T) Action {
	return &actionFunc{
		value: func(t require.TestingT) (any, bool) {
			return length(0), true
		},
		run: func(t require.TestingT, ctx context.Context, v any, raw json.RawMessage, err error) error {
			if err != nil {
				return err
			}
			if !assert.ObjectsAreEqual(v, length(l)) {
				return fmt.Errorf("%w: expected length %d, got %d", ErrActionAssertLen, l, v)
			}
			return nil
		},
		support:   []Support{SupportJson},
		baseError: ErrActionAssertLen,
	}
}

// AssertLt asserts that given value is less than the value in the response
func AssertLt[T constraints.Integer](value T) Action {
	return &actionFunc{
		value: func(t require.TestingT) (any, bool) {
			return value, true
		},
		run: func(t require.TestingT, ctx context.Context, v any, raw json.RawMessage, err error) error {
			if err != nil {
				return err
			}
			if !assert.Less(&requireTestingT{}, v, value) {
				return fmt.Errorf("%w: expected %v to be less than %v", ErrActionAssertLt, v, value)
			}
			return nil
		},
		support:   []Support{SupportJson, SupportStatus},
		baseError: ErrActionAssertLt,
	}
}

// AssertLte asserts that given value is less than or equal to the value in the response
func AssertLte[T constraints.Integer](value T) Action {
	return &actionFunc{
		value: func(t require.TestingT) (any, bool) {
			return value, true
		},
		run: func(t require.TestingT, ctx context.Context, v any, raw json.RawMessage, err error) error {
			if err != nil {
				return err
			}
			if !assert.LessOrEqual(&requireTestingT{}, v, value) {
				return fmt.Errorf("%w: expected %v to be less than or equal to %v", ErrActionAssertLte, v, value)

			}
			return nil
		},
		support:   []Support{SupportJson, SupportStatus},
		baseError: ErrActionAssertLte,
	}
}

// AssertNotEquals asserts that given value is not equal to the value in the response
func AssertNotEquals(value any) Action {
	return &actionFunc{
		value: func(t require.TestingT) (any, bool) {
			return value, true
		},
		run: func(t require.TestingT, ctx context.Context, v any, raw json.RawMessage, err error) error {
			if err != nil {
				return err
			}
			if assert.ObjectsAreEqual(value, v) {
				return fmt.Errorf("%w: expected value to not be equal to %#v", ErrActionAssertNotEquals, value)
			}
			return nil
		},
		support:   []Support{SupportHeader, SupportJson, SupportStatus},
		baseError: ErrActionAssertNotEquals,
	}
}

// AssertNotExists asserts that the value does not exist in the response based on given argument
func AssertNotExists() Action {
	return &actionFunc{
		run: func(t require.TestingT, ctx context.Context, v any, raw json.RawMessage, err error) error {
			if err != nil {
				return err
			}
			if raw != nil {
				return fmt.Errorf("%w: expected value to not exist, but it does", ErrActionAssertNotExists)
			}
			return nil
		},
		support:   []Support{SupportHeader, SupportJson},
		baseError: ErrActionAssertNotExists,
	}
}

// AssertNotIn asserts that value is not in the given list of values
func AssertNotIn[T any](values ...T) Action {
	return &actionFunc{
		value: func(t require.TestingT) (any, bool) {
			return values[0], true
		},
		run: func(t require.TestingT, ctx context.Context, v any, raw json.RawMessage, err error) error {
			if err != nil {
				return err
			}
			ok, found := containsElement(values, v)
			if !ok || !found {
				return nil
			}
			return fmt.Errorf("%w: expected value %#v to not be in %#v", ErrActionAssertNotIn, v, values)
		},
		support:   []Support{SupportHeader, SupportJson, SupportStatus},
		baseError: ErrActionAssertNotIn,
	}
}

// AssertRegex asserts that given value matches the regex in the response
// if count provided it will check for the number of matches,
// otherwise it will check if the value matches the regex
func AssertRegex(pattern *regexp.Regexp, count ...int) Action {
	return &actionFunc{
		run: func(t require.TestingT, ctx context.Context, v any, raw json.RawMessage, err error) error {
			if err != nil {
				return err
			}
			if len(count) > 0 {
				matches := pattern.FindAllSubmatch(raw, -1)
				if len(matches) != count[0] {
					return fmt.Errorf("expected %d matches, got %d", count[0], len(matches))
				}
			} else {
				if !pattern.Match(raw) {
					return fmt.Errorf("expected %s to match %s", raw, pattern)
				}
			}
			return nil
		},
		support:   []Support{SupportHeader, SupportJson, SupportStatus},
		baseError: ErrActionAssertRegex,
	}
}

// length type to handle the length of a slice or map
type length int

// UnmarshalJSON unmarshal the JSON data into a length type
func (l *length) UnmarshalJSON(b []byte) error {
	// first unmarshal slice
	var arr []any

	if err := json.Unmarshal(b, &arr); err != nil {
		// second unmarshal map
		var m map[string]any
		if err := json.Unmarshal(b, &m); err != nil {
			return err
		}
		*l = length(len(m))
	}
	*l = length(len(arr))

	return nil
}
