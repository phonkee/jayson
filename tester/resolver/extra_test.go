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

package resolver

import (
	"github.com/phonkee/jayson/tester/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strings"
	"testing"
)

// MatchByStringContains matches string by substring
func matchByStringContains(s string) func(in string) bool {
	return func(in string) bool {
		return strings.Contains(in, s)
	}
}

func TestResolverExtraImpl_Args(t *testing.T) {
	e := resolverExtra(nil, nil)
	assert.Nil(t, e.Args())
}

func TestResolverExtraImpl_Query(t *testing.T) {
	e := resolverExtra(nil, nil)
	assert.Nil(t, e.Query())
}

func TestArguments(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		a := Arguments(t, "arg", "value")
		assert.Equal(t, []string{"arg", "value"}, a.Args())
	})
	t.Run("invalid", func(t *testing.T) {
		m := mocks.NewTestingT(t)
		m.On("Errorf", mock.Anything, mock.MatchedBy(matchByStringContains("arguments must be in key/value pairs"))).Once()
		m.On("FailNow").Once()
		Arguments(m, "hello")
	})
}

func TestQuery(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		q := Query(t, "key", "value")
		assert.Equal(t, "value", q.Query().Get("key"))
	})
	t.Run("invalid", func(t *testing.T) {
		m := mocks.NewTestingT(t)
		m.On("Errorf", mock.Anything, mock.MatchedBy(matchByStringContains("arguments must be in key/value pairs"))).Once()
		m.On("FailNow").Once()
		Query(m, "key")

	})
}
