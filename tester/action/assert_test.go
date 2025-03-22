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
	"github.com/phonkee/jayson/tester/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"regexp"
	"strconv"
	"testing"
)

func TestAssertBetween(t *testing.T) {
	t.Run("test valid values", func(t *testing.T) {
		for _, test := range []struct {
			name   string
			action Action
			value  any
		}{
			{name: "test positive int", action: AssertBetween(1, 10), value: 5},
			{name: "test positive int8", action: AssertBetween(int8(1), int8(10)), value: int8(5)},
			{name: "test positive int16", action: AssertBetween(int16(1), int16(10)), value: int16(5)},
			{name: "test positive int32", action: AssertBetween(int32(1), int32(10)), value: int32(5)},
			{name: "test positive int64", action: AssertBetween(int64(1), int64(10)), value: int64(5)},
			{name: "test negative int", action: AssertBetween(-1000, -400), value: -500},
			{name: "test negative int8", action: AssertBetween(int8(-100), int8(-50)), value: int8(-75)},
			{name: "test negative int16", action: AssertBetween(int16(-100), int16(-50)), value: int16(-75)},
			{name: "test negative int32", action: AssertBetween(int32(-100), int32(-50)), value: int32(-75)},
			{name: "test negative int64", action: AssertBetween(int64(-100), int64(-50)), value: int64(-75)},
			{name: "test uint", action: AssertBetween(0, 100), value: 50},
			{name: "test uint8", action: AssertBetween(uint8(1), uint8(10)), value: uint8(5)},
			{name: "test uint16", action: AssertBetween(uint16(1), uint16(10)), value: uint16(5)},
			{name: "test uint32", action: AssertBetween(uint32(1), uint32(10)), value: uint32(5)},
			{name: "test uint64", action: AssertBetween(uint64(1), uint64(10)), value: uint64(5)},
		} {
			t.Run(test.name, func(t *testing.T) {
				m := mocks.NewTestingT(t)
				require.NoError(t, test.action.Run(m, context.Background(), test.value, nil, nil))
			})
		}
	})
	t.Run("test invalid values", func(t *testing.T) {
		for _, test := range []struct {
			name   string
			action Action
			value  any
		}{
			{name: "test positive int", action: AssertBetween(1, 10), value: 11},
			{name: "test positive int8", action: AssertBetween(int8(1), int8(10)), value: int8(11)},
			{name: "test positive int16", action: AssertBetween(int16(1), int16(10)), value: int16(11)},
			{name: "test positive int32", action: AssertBetween(int32(1), int32(10)), value: int32(11)},
			{name: "test positive int64", action: AssertBetween(int64(1), int64(10)), value: int64(11)},
			{name: "test negative int", action: AssertBetween(-1000, -400), value: -399},
			{name: "test negative int8", action: AssertBetween(int8(-100), int8(-50)), value: int8(-49)},
			{name: "test negative int16", action: AssertBetween(int16(-100), int16(-50)), value: int16(-49)},
			{name: "test negative int32", action: AssertBetween(int32(-100), int32(-50)), value: int32(-49)},
			{name: "test negative int64", action: AssertBetween(int64(-100), int64(-50)), value: int64(-49)},
			{name: "test uint", action: AssertBetween(0, 100), value: 101},
			{name: "test uint8", action: AssertBetween(uint8(1), uint8(10)), value: uint8(11)},
			{name: "test uint16", action: AssertBetween(uint16(1), uint16(10)), value: uint16(11)},
			{name: "test uint32", action: AssertBetween(uint32(1), uint32(10)), value: uint32(11)},
			{name: "test uint64", action: AssertBetween(uint64(1), uint64(10)), value: uint64(11)},
		} {
			t.Run(test.name, func(t *testing.T) {
				m := mocks.NewTestingT(t)
				require.Error(t, test.action.Run(m, context.Background(), test.value, nil, nil))
			})
		}
	})
}

func TestAssertContains(t *testing.T) {
	t.Run("test valid values", func(t *testing.T) {
		for _, test := range []struct {
			name string
			sub  string
			raw  json.RawMessage
		}{
			{name: "test string", sub: "hello", raw: json.RawMessage(`{"hello": "world"}`)},
			{name: "test integer", sub: "123", raw: json.RawMessage(`{"hello": 123}`)},
		} {
			t.Run(test.name, func(t *testing.T) {
				ac := AssertContains(test.sub)
				m := mocks.NewTestingT(t)
				require.NoError(t, ac.Run(m, context.Background(), nil, test.raw, nil))
			})
		}
	})

	t.Run("test invalid value", func(t *testing.T) {
		for _, test := range []struct {
			name string
			sub  string
			raw  json.RawMessage
		}{
			{name: "test string", sub: "something", raw: json.RawMessage(`{"hello": "world"}`)},
			{name: "test integer", sub: "42", raw: json.RawMessage(`{"hello": 123}`)},
		} {
			t.Run(test.name, func(t *testing.T) {
				ac := AssertContains(test.sub)
				m := mocks.NewTestingT(t)
				require.Error(t, ac.Run(m, context.Background(), nil, test.raw, nil))
			})
		}
	})

}

func TestAssertEquals(t *testing.T) {
	t.Run("test valid values", func(t *testing.T) {
		for _, test := range []struct {
			name   string
			expect any
			value  any
		}{
			{name: "test string", expect: "hello", value: "hello"},
			{name: "test int", expect: 42, value: 42},
			{name: "test uint", expect: uint(42), value: uint(42)},
		} {
			t.Run(test.name, func(t *testing.T) {
				ae := AssertEquals(test.expect)
				m := mocks.NewTestingT(t)
				require.NoError(t, ae.Run(m, context.Background(), test.value, nil, nil))
			})
		}

	})

	t.Run("test invalid value", func(t *testing.T) {
		for _, test := range []struct {
			name          string
			expect        any
			value         any
			errIn         error
			errorContains string
		}{
			{name: "test invalid value", expect: "hello", value: "world", errIn: nil, errorContains: "action: `AssertEquals`, expected: \"hello\", got: \"world\""},
			{name: "test passed error", expect: 42, value: 43, errIn: ErrNotPresent, errorContains: "not present"},
		} {
			t.Run(test.name, func(t *testing.T) {
				ae := AssertEquals(test.expect)
				m := mocks.NewTestingT(t)
				err := ae.Run(m, context.Background(), test.value, nil, test.errIn)
				require.Error(t, err)
				assert.Contains(t, err.Error(), test.errorContains)
			})
		}

	})

}

func TestAssertKeysIn(t *testing.T) {
	t.Run("test valid values", func(t *testing.T) {
		// Test cases for AssertKeysIn
		tests := []struct {
			name  string
			value map[string]any
			keys  []string
		}{
			{name: "test subset", value: map[string]any{"a": 1, "b": 2}, keys: []string{"a"}},
			{name: "test other", value: map[string]any{"a": 1, "b": 2}, keys: []string{"b"}},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				aak := AssertKeysIn(test.keys...)
				m := mocks.NewTestingT(t)
				require.NoError(t, aak.Run(m, context.Background(), test.value, nil, nil))
			})
		}

	})
	t.Run("test invalid values", func(t *testing.T) {
		tests := []struct {
			name  string
			value map[string]any
			keys  []string
		}{
			{name: "test non existing", value: map[string]any{"a": 1, "b": 2}, keys: []string{"c"}},
			{name: "test non existing 2", value: map[string]any{"a": 1, "b": 2}, keys: []string{"e"}},
			{name: "test nil", value: nil, keys: []string{"e"}},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				aak := AssertKeysIn(test.keys...)
				m := mocks.NewTestingT(t)
				require.Error(t, aak.Run(m, context.Background(), test.value, nil, nil))
			})
		}
	})
}

func TestAssertZero(t *testing.T) {
	t.Run("test valid values", func(t *testing.T) {
		// Test cases for AssertZero
		tests := []struct {
			input json.RawMessage
		}{
			{json.RawMessage(`0`)},
			{json.RawMessage(`""`)},
			{json.RawMessage(`[]`)},
			{json.RawMessage(`{}`)},
			{json.RawMessage(`false`)},
		}

		for i, test := range tests {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				ass := AssertZero()
				m := mocks.NewTestingT(t)
				require.NoError(t, ass.Run(m, context.Background(), nil, test.input, nil))
			})
		}
	})

	t.Run("test invalid values", func(t *testing.T) {
		// Test cases for AssertZero
		tests := []struct {
			input json.RawMessage
		}{
			{json.RawMessage(`1`)},
			{json.RawMessage(`"hello"`)},
			{json.RawMessage(`[1]`)},
			{json.RawMessage(`{"hello": "world""}`)},
			{json.RawMessage(`true`)},
		}

		for i, test := range tests {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				ass := AssertZero()
				m := mocks.NewTestingT(t)
				require.Error(t, ass.Run(m, context.Background(), nil, test.input, nil))
			})
		}
	})
}

func TestAssertNotZero(t *testing.T) {
	t.Run("test valid values", func(t *testing.T) {
		// Test cases for AssertZero
		tests := []struct {
			input json.RawMessage
		}{
			{json.RawMessage(`0`)},
			{json.RawMessage(`""`)},
			{json.RawMessage(`[]`)},
			{json.RawMessage(`{}`)},
			{json.RawMessage(`false`)},
		}

		for i, test := range tests {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				ass := AssertNot(AssertZero())
				m := mocks.NewTestingT(t)
				require.Error(t, ass.Run(m, context.Background(), nil, test.input, nil))
			})
		}
	})

	t.Run("test invalid values", func(t *testing.T) {
		// Test cases for AssertZero
		tests := []struct {
			input json.RawMessage
		}{
			{json.RawMessage(`1`)},
			{json.RawMessage(`"hello"`)},
			{json.RawMessage(`[1]`)},
			{json.RawMessage(`{"hello": "world""}`)},
			{json.RawMessage(`true`)},
		}

		for i, test := range tests {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				ass := AssertNot(AssertZero())
				m := mocks.NewTestingT(t)
				require.NoError(m, ass.Run(m, context.Background(), nil, test.input, nil))
			})
		}
	})
}

func TestAssertRegexMatch(t *testing.T) {
	t.Run("test valid values", func(t *testing.T) {
		for _, item := range []struct {
			name    string
			input   json.RawMessage
			pattern *regexp.Regexp
		}{
			{
				name:    "match regex",
				input:   json.RawMessage(`200`),
				pattern: regexp.MustCompile(`2\d+`),
			},
		} {
			t.Run(item.name, func(t *testing.T) {
				ar := AssertRegexMatch(item.pattern)
				assert.NoError(t, ar.Run(t, context.Background(), nil, item.input, nil))
			})

		}

	})
	t.Run("test invalid values", func(t *testing.T) {
		for _, item := range []struct {
			name    string
			input   json.RawMessage
			pattern *regexp.Regexp
		}{
			{
				name:    "match regex",
				input:   json.RawMessage(`hello`),
				pattern: regexp.MustCompile(`2\d+`),
			},
		} {
			t.Run(item.name, func(t *testing.T) {
				ar := AssertRegexMatch(item.pattern)
				assert.Error(t, ar.Run(t, context.Background(), nil, item.input, nil))
			})
		}
	})
}

func TestAssertRegexSearch(t *testing.T) {
	t.Run("test valid values", func(t *testing.T) {
		for _, item := range []struct {
			name    string
			input   json.RawMessage
			pattern *regexp.Regexp
			count   int
		}{
			{
				name:    "match regex",
				input:   json.RawMessage(`123`),
				pattern: regexp.MustCompile(`\d`),
				count:   2,
			},
		} {
			t.Run(item.name, func(t *testing.T) {
				ar := AssertRegexSearch(item.pattern, item.count)
				assert.Error(t, ar.Run(t, context.Background(), nil, item.input, nil))
			})
		}
	})
	t.Run("test invalid values", func(t *testing.T) {
		for _, item := range []struct {
			name    string
			input   json.RawMessage
			pattern *regexp.Regexp
			count   int
		}{
			{
				name:    "match regex",
				input:   json.RawMessage(`123`),
				pattern: regexp.MustCompile(`\d`),
				count:   2,
			},
		} {
			t.Run(item.name, func(t *testing.T) {
				ar := AssertRegexSearch(item.pattern, item.count)
				assert.Error(t, ar.Run(t, context.Background(), nil, item.input, nil))
			})
		}
	})
}
