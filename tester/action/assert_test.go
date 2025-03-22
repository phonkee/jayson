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
