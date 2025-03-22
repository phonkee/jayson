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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/constraints"
	"reflect"
	"regexp"
	"sort"
)

// AssertBetween asserts that given integer value is between the values (inclusive)
func AssertBetween[T constraints.Integer](min, max T) Action {
	return &actionFunc{
		value: func(t require.TestingT) (any, bool) {
			return min, true
		},
		run: func(t require.TestingT, ctx context.Context, v any, raw json.RawMessage, err error) error {
			if err != nil {
				return err
			}
			if !assert.GreaterOrEqual(&requireTestingT{}, v, min) {
				return fmt.Errorf("%w: expected %v to be between [%v, %v]", ErrActionAssertBetween, v, min, max)
			}
			if !assert.LessOrEqual(&requireTestingT{}, v, max) {
				return fmt.Errorf("%w: expected %v to be between [%v, %v]", ErrActionAssertBetween, v, min, max)
			}
			return nil
		},
		support:   []Support{SupportJson, SupportStatus},
		baseError: ErrActionAssertBetween,
	}
}

// AssertContains asserts value contains the given value
func AssertContains(contains string) Action {
	return &actionFunc{
		run: func(t require.TestingT, ctx context.Context, v any, raw json.RawMessage, err error) error {
			if err != nil {
				return err
			}
			if !assert.Contains(&requireTestingT{}, string(raw), contains) {
				return fmt.Errorf("%w: expected %s to contain %s", ErrActionAssertContains, string(raw), contains)
			}
			return nil
		},
		support:   []Support{SupportHeader, SupportJson},
		baseError: ErrActionAssertContains,
	}
}

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

			// special case for json.RawMessage
			// we call json eq to compare the values correctly
			if reflect.TypeOf(v) == reflect.TypeOf(json.RawMessage{}) {
				if !assert.JSONEq(&requireTestingT{}, string(v.(json.RawMessage)), string(raw)) {
					return fmt.Errorf("%w, expected: %#v, got: %#v", ErrActionAssertEquals, value, v)
				}
			} else {
				if !assert.ObjectsAreEqual(v, value) {
					return fmt.Errorf("%w, expected: %#v, got: %#v", ErrActionAssertEquals, value, v)
				}
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
				if errors.Is(err, ErrNotPresent) {
					return fmt.Errorf("%w: expected value to exist, but it does not", ErrActionAssertExists)
				}
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
				return fmt.Errorf("%w: expected %v to be greater than %v", ErrActionAssertGt, v, value)
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

// AssertKeysIn asserts that given map has the given keys
func AssertKeysIn(keys ...string) Action {
	return &actionFunc{
		value: func(t require.TestingT) (any, bool) {
			return map[string]any{}, true
		},
		run: func(t require.TestingT, ctx context.Context, v any, raw json.RawMessage, err error) error {
			if err != nil {
				return err
			}
			var valueMap map[string]any
			if v == nil {
				valueMap = make(map[string]any)
			}
			if cast, ok := v.(map[string]any); ok {
				valueMap = cast
			}

			// collect all the keys from the map
			found := make([]string, 0, len(keys))
			for key := range valueMap {
				found = append(found, key)
			}

			var isThere bool

		outer:
			for _, f := range found {
				for _, k := range keys {
					if assert.ObjectsAreEqual(f, k) {
						isThere = true
						break outer
					}
				}
			}

			if !isThere {
				return fmt.Errorf("%w: expected keys %#v be in %#v", ErrActionAssertKeys, found, keys)
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

// AssertNil asserts that given value is nil
func AssertNil() Action {
	return &actionFunc{
		run: func(t require.TestingT, ctx context.Context, v any, raw json.RawMessage, err error) error {
			if err != nil {
				return err
			}
			if raw == nil {
				return fmt.Errorf("%w: value missing", ErrActionAssertNil)
			}
			if !assert.Equal(&requireTestingT{}, nullRawMessage, json.RawMessage(bytes.ToLower(raw))) {
				return fmt.Errorf("%w: expected value to be nil, got %#v", ErrActionAssertNil, string(raw))
			}
			return nil
		},
		support:   []Support{SupportHeader, SupportJson},
		baseError: ErrActionAssertNil,
	}
}

// AssertRegexMatch asserts that given value matches the regex in the response
// if count provided it will check for the number of matches,
// otherwise it will check if the value matches the regex
func AssertRegexMatch(pattern *regexp.Regexp) Action {
	return &actionFunc{
		run: func(t require.TestingT, ctx context.Context, v any, raw json.RawMessage, err error) error {
			if err != nil {
				return err
			}
			if raw == nil {
				return fmt.Errorf("expected value to exist, but it does not")
			}

			// match regular expression
			if !pattern.Match(raw) {
				return fmt.Errorf("expected %s to match regex: %#v", raw, pattern.String())
			}
			return nil
		},
		support:   []Support{SupportHeader, SupportJson, SupportStatus},
		baseError: ErrActionAssertRegexMatch,
	}
}

// AssertRegexSearch asserts that given regular expression is found given number of times
func AssertRegexSearch(pattern *regexp.Regexp, count int) Action {
	return &actionFunc{
		run: func(t require.TestingT, ctx context.Context, v any, raw json.RawMessage, err error) error {
			if err != nil {
				return err
			}
			if raw == nil {
				return fmt.Errorf("expected value to exist, but it does not")
			}

			// find all matches
			matches := pattern.FindAllSubmatch(raw, -1)
			if len(matches) != count {
				return fmt.Errorf("expected %d matches, got %d for regex: %#v", count, len(matches), pattern.String())
			}

			return nil
		},
		support:   []Support{SupportHeader, SupportJson, SupportStatus},
		baseError: ErrActionAssertRegexSearch,
	}
}

// AssertZero asserts that given value is zero
func AssertZero() Action {
	return &actionFunc{
		run: func(t require.TestingT, ctx context.Context, value any, raw json.RawMessage, err error) error {
			if err != nil {
				return err
			}

			for _, zeroValue := range zeroValues {
				if assert.JSONEq(&requireTestingT{}, string(zeroValue), string(raw)) {
					return nil
				}
			}

			return fmt.Errorf("%w: expected value to be zero, got %#v", ErrAction, raw)

		},

		support: []Support{SupportHeader, SupportJson},
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
			return fmt.Errorf("%s", "cannot get length of value that is not slice or map")
		}
		*l = length(len(m))
	} else {
		*l = length(len(arr))
	}

	return nil
}

// AssertNot asserts that given action does not succeed, otherwise it fails
// This assertion does not need any special handling, it just relies on correct errors (ErrAction, ErrNotPresent)
func AssertNot(a Action) Action {
	return &actionFunc{
		value: func(t require.TestingT) (any, bool) {
			return a.Value(t)
		},
		run: func(t require.TestingT, ctx context.Context, v any, raw json.RawMessage, err error) error {
			// run the provided action
			runErr := a.Run(t, ctx, v, raw, err)

			// check for error
			if runErr != nil {
				if errors.Is(runErr, ErrAction) {
					return nil
				}
				return fmt.Errorf("expected action to fail, but it did not: %w", runErr)
			}
			return fmt.Errorf("expected action to fail, but it did not")
		},
		supportsFunc: a.Supports,
		baseError:    ErrActionAssertNot,
	}
}
