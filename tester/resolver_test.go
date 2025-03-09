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

package tester_test

import (
	"github.com/phonkee/jayson/tester"
	"github.com/phonkee/jayson/tester/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestNewGorillaResolver(t *testing.T) {
	t.Run("test invalid router", func(t *testing.T) {
		m := mocks.NewTestingT(t)
		m.On("Errorf", mock.Anything, mock.MatchedBy(matchByStringContains("resolver: Router is nil"))).Once()
		m.On("FailNow").Once()

		tester.NewGorillaResolver(m, nil)
	})

	t.Run("test valid router", func(t *testing.T) {
		m := mocks.NewTestingT(t)
		router := newHealthRouter(m)
		resolver := tester.NewGorillaResolver(m, router)
		assert.NotNil(t, resolver)
	})
}

func TestGorillaResolver_ReverseURL(t *testing.T) {

	t.Run("test simple url", func(t *testing.T) {
		router := newHealthRouter(t)
		m := mocks.NewTestingT(t)
		resolver := tester.NewGorillaResolver(m, router)
		assert.NotNil(t, resolver)
		url := resolver.ReverseURL(t, "api:v1:health")
		assert.Equal(t, "/api/v1/health", url)
	})

	t.Run("test url with arguments", func(t *testing.T) {
		router := newHealthRouter(t)
		m := mocks.NewTestingT(t)
		resolver := tester.NewGorillaResolver(m, router)
		assert.NotNil(t, resolver)
		url := resolver.ReverseURL(t,
			"api:v1:health:extra",
			tester.ResolverArgs(t, "component", "database"),
		)
		assert.Equal(t, "/api/v1/health/database", url)
	})

	t.Run("test url with query", func(t *testing.T) {
		router := newHealthRouter(t)
		m := mocks.NewTestingT(t)
		resolver := tester.NewGorillaResolver(m, router)
		assert.NotNil(t, resolver)
		url := resolver.ReverseURL(t,
			"api:v1:health:extra",
			tester.ResolverArgs(t, "component", "database"),
			tester.ResolverQuery(t, "page", "1"),
		)
		assert.Equal(t, "/api/v1/health/database?page=1", url)
	})

}
