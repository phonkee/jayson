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
	"github.com/phonkee/jayson/tester/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// MatchByStringContains matches string by substring
func matchByStringContains(s string) func(in string) bool {
	return func(in string) bool {
		return strings.Contains(in, s)
	}
}

type testStruct struct {
	Value string `json:"value"`
}

func TestAPIObject(t *testing.T) {
	t.Run("test odd number of arguments", func(t *testing.T) {
		m := mocks.NewTestingT(t)
		m.On("Errorf", mock.Anything, mock.MatchedBy(matchByStringContains("APIObject: odd number of arguments"))).Once()
		m.On("FailNow")

		// test for odd number of arguments
		APIObject(m, "key")

		m.AssertExpectations(t)
		m.AssertNumberOfCalls(t, "Errorf", 1)
		m.AssertNumberOfCalls(t, "FailNow", 1)
	})

	t.Run("test settable", func(t *testing.T) {

		t.Run("test non settable", func(t *testing.T) {
			for _, item := range []struct {
				name  string
				value any
			}{
				{
					name:  "test string",
					value: "world",
				},
				{
					name:  "test int",
					value: int(1),
				},
				{
					name:  "test bool",
					value: false,
				},
				{
					name:  "test struct",
					value: testStruct{},
				},
			} {
				t.Run(item.name, func(t *testing.T) {
					m := mocks.NewTestingT(t)
					m.On("Errorf", mock.Anything, mock.MatchedBy(matchByStringContains("APIObject: value `value` is not settable"))).Once()
					m.On("FailNow")
					APIObject(m, "value", item.value)
					m.AssertExpectations(t)
					m.AssertNumberOfCalls(t, "Errorf", 1)
					m.AssertNumberOfCalls(t, "FailNow", 1)
				})
			}
		})

		t.Run("test settable", func(t *testing.T) {
			for _, item := range []struct {
				name  string
				value any
			}{
				{
					name:  "test string",
					value: new(string),
				},
				{
					name:  "test int",
					value: new(int),
				},
				{
					name:  "test bool",
					value: new(bool),
				},
				{
					name:  "test struct",
					value: new(testStruct),
				},
			} {
				t.Run(item.name, func(t *testing.T) {
					m := mocks.NewTestingT(t)
					//m.On("Errorf", mock.Anything, mock.MatchedBy(matchByStringContains("APIObject: value `value` is not settable"))).Once()
					//m.On("FailNow")
					APIObject(m, "value", item.value)
					m.AssertExpectations(t)
					//m.AssertNumberOfCalls(t, "Errorf", 0)
					//m.AssertNumberOfCalls(t, "FailNow", 0)
				})
			}
		})

		t.Run("test unmarshalling", func(t *testing.T) {
			for _, item := range []struct {
				name   string
				value  any
				body   []byte
				expect any
			}{
				{
					name:   "test string",
					value:  new(string),
					body:   []byte(`{"value":"hello"}`),
					expect: ptrTo("hello"),
				},
				{
					name:   "test int",
					value:  new(int),
					body:   []byte(`{"value":1}`),
					expect: ptrTo(1),
				},
				{
					name:   "test bool",
					value:  new(bool),
					body:   []byte(`{"value":true}`),
					expect: ptrTo(true),
				},
				{
					name:   "test struct",
					value:  new(testStruct),
					body:   []byte(`{"value": {"value": "hello"}}`),
					expect: &testStruct{Value: "hello"},
				},
			} {
				t.Run(item.name, func(t *testing.T) {
					m := mocks.NewTestingT(t)
					resp := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
					resp.body = item.body
					//resp.Json(m, APIObject(m, "value", item.value))
					assert.Equal(t, item.expect, item.value)
					m.AssertExpectations(t)
				})
			}
		})

	})

}

func ptrTo[T any](v T) *T {
	return &v
}
