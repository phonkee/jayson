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
	"github.com/phonkee/jayson/tester/action"
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
	Value string `json:"action"`
}

func ptrTo[T any](v T) *T {
	return &v
}

func TestResponse_Header(t *testing.T) {
	t.Run("test contains", func(t *testing.T) {
		for _, item := range []struct {
			name    string
			headers http.Header
			key     string
			action  action.Action
		}{
			{
				name:    "header exists",
				headers: map[string][]string{"Hello": {"World"}},
				key:     "Hello",
				action:  action.AssertEquals("World"),
			},
			{
				name:    "header does not exist",
				headers: map[string][]string{"Hello": {"This", "World"}},
				key:     "Hello",
				action:  action.AssertEquals("This"),
			},
		} {
			t.Run(item.name, func(t *testing.T) {
				rw := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				for k, v := range item.headers {
					rw.Header()[k] = v
				}
				r := newResponse(rw, req)
				r.Header(t, item.key, item.action)
			})
		}
	})

	t.Run("test does not contain", func(t *testing.T) {
		for _, item := range []struct {
			name    string
			headers http.Header
			key     string
			action  action.Action
		}{
			{
				name:    "header exists",
				headers: map[string][]string{"Hello": {"World"}},
				key:     "Hello",
				action:  action.AssertExists(),
			},
			{
				name:    "header does not exist",
				headers: map[string][]string{"Hello": {"This", "World"}},
				key:     "Nope",
				action:  action.AssertNotExists(),
			},
		} {
			t.Run(item.name, func(t *testing.T) {
				rw := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				for k, v := range item.headers {
					rw.Header()[k] = v
				}
				r := newResponse(rw, req)

				// mock testing
				m := mocks.NewTestingT(t)
				//m.On("Errorf", mock.MatchedBy(matchByStringContains(item.expect)))
				//m.On("FailNow")
				r.Header(m, item.key, item.action)
				//m.AssertExpectations(t)
				//m.AssertNumberOfCalls(t, "Errorf", 1)
				//m.AssertNumberOfCalls(t, "FailNow", 1)
			})
		}
	})

}

func TestResponse_Json(t *testing.T) {

	t.Run("test valid path", func(t *testing.T) {
		for _, item := range []struct {
			name   string
			body   string
			path   string
			action action.Action
		}{
			{
				name:   "simple strings",
				body:   `{"name": "John"}`,
				path:   "name",
				action: action.AssertEquals("John"),
			},
			{
				name:   "whole object as RawMessage",
				body:   `{"name": "John"}`,
				path:   "",
				action: action.AssertEquals(json.RawMessage(`{"name": "John"}`)),
			},
			{
				name:   "simple strings to pointer",
				body:   `{"name": "John"}`,
				path:   "name",
				action: action.AssertEquals(ptrTo("John")),
			},
			{
				name:   "test embedded struct strings",
				body:   `{"other": {"name": "John"}}`,
				path:   "other.name",
				action: action.AssertEquals("John"),
			},
			{
				name:   "test slice of objects to string",
				body:   `{"other": [{"name": "John"}, {"name": "Doe"}]}`,
				path:   "other.1.name",
				action: action.AssertEquals("Doe"),
			},
			{
				name:   "test slice of objects to object",
				body:   `{"other": [{"name": "John"}, {"name": "Doe"}, {"object": {"action": "Doe"}}]}`,
				path:   "other.2.object",
				action: action.AssertEquals(testStruct{Value: "Doe"}),
			},
			{
				name:   "test slice of objects to pointer to object",
				body:   `{"other": [{"name": "John"}, {"name": "Doe"}, {"object": {"action": "Doe"}}]}`,
				path:   "other.2.object",
				action: action.AssertEquals(&testStruct{Value: "Doe"}),
			},
			{
				name:   "test slice of objects to pointer to integer",
				body:   `{"other": [{"name": "John"}, {"name": "Doe"}, {"object": {"action": 42}}]}`,
				path:   "other.2.object.action",
				action: action.AssertEquals(42),
			},
			{
				name:   "test slice of objects to pointer to integer (with dots)",
				body:   `{"other": [{"name": "John"}, {"name": "Doe"}, {"object": {"action": 42}}]}`,
				path:   "other.2......object.action",
				action: action.AssertEquals(ptrTo(42)),
			},
			{
				name:   "test raw json",
				body:   `{"other": [{"name": "John"}, {"name": "Doe"}, {"object": {"action": 42, "other": 12}}]}`,
				path:   "other.2",
				action: action.AssertEquals(json.RawMessage(`{"object": {"action": 42, "other": 12}}`)),
			},
			{
				name: "test raw json",
				body: `{"other": [{"action": "John"}, {"action": "Doe"}, {"action": "Mark", "extra": "extra"}]}`,
				path: "other",
				action: action.AssertEquals([]testStruct{
					{Value: "John"},
					{Value: "Doe"},
					{Value: "Mark"},
				}),
			},
		} {
			t.Run(item.name, func(t *testing.T) {
				r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
				r.body = []byte(item.body)

				r.Json(t, item.path, item.action)
			})
		}
	})

	t.Run("test Assert", func(t *testing.T) {
		t.Run("test AssertLen", func(t *testing.T) {
			t.Run("valid path", func(t *testing.T) {
				for _, item := range []struct {
					name   string
					body   string
					path   string
					action action.Action
				}{
					{
						name:   "test len of array",
						body:   `{"other": [{"name": "John"}, {"name": "Doe"}, {"object": {"action": 42, "other": 12}}]}`,
						path:   "other",
						action: action.AssertLen(3),
					},
					{
						name:   "test len of array",
						body:   `{"other": [{"name": "John"}, {"name": "Doe"}, {"object": {"action": 42, "other": 12}}]}`,
						path:   "other.2.object",
						action: action.AssertLen(2),
					},
				} {
					t.Run(item.name, func(t *testing.T) {
						r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
						r.body = []byte(item.body)

						r.Json(t, item.path, item.action)
					})
				}
			})
			t.Run("invalid path", func(t *testing.T) {
				for _, item := range []struct {
					name   string
					body   string
					path   string
					what   action.Action
					expect string
				}{
					{
						name:   "test len of string",
						body:   `{"other": [{"name": "John"}, {"name": "Doe"}, {"object": {"action": 42, "other": 12}}]}`,
						path:   "other.0.name",
						what:   action.AssertLen(1),
						expect: "FAILED: `response.Json`, path: `other.0.name`, unmarshal error: cannot get length of value that is not slice or map",
					},
					{
						name:   "test len of int",
						body:   `{"other": [{"name": "John"}, {"name": "Doe"}, {"object": {"action": 42, "other": 12}}]}`,
						path:   "other.2.object.other",
						what:   action.AssertLen(1),
						expect: "FAILED: `response.Json`, path: `other.2.object.other`, unmarshal error: cannot get length of value that is not slice or map",
					},
					{
						name: "test len of int",
						body: `{"other": [{"name": "John"}, {"name": "Doe"}, {"object": {"action": 42, "other": 12}}]}`,
						path: "other.2.object.other",
						what: action.AssertLen(1),
					},
				} {
					t.Run(item.name, func(t *testing.T) {
						r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
						r.body = []byte(item.body)

						m := mocks.NewTestingT(t)
						m.On("Errorf", mock.Anything, mock.MatchedBy(matchByStringContains(item.expect))).Once()
						m.On("FailNow")

						r.Json(m, item.path, item.what)

						m.AssertNumberOfCalls(t, "Errorf", 1)
						m.AssertNumberOfCalls(t, "FailNow", 1)
						m.AssertExpectations(t)
					})
				}

			})
		})
		t.Run("test AssertKeys", func(t *testing.T) {
			for _, item := range []struct {
				name   string
				body   string
				path   string
				action action.Action
			}{
				{
					name:   "test keys of object",
					body:   `{"name": "John", "age": 42}`,
					path:   "",
					action: action.AssertKeys("name", "age"),
				},
				{
					name:   "test keys of object",
					body:   `{"obj": {"name": "John", "age": 42}}`,
					path:   "obj",
					action: action.AssertKeys("name", "age"),
				},
			} {
				t.Run(item.name, func(t *testing.T) {
					r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
					r.body = []byte(item.body)

					r.Json(t, item.path, item.action)
				})
			}

		})
		t.Run("test AssertExists", func(t *testing.T) {
			for _, item := range []struct {
				name string
				body string
				path string
			}{
				{
					name: "test keys of object",
					body: `{"name": "John", "age": 42}`,
					path: "name",
				},
				{
					name: "test keys of object",
					body: `{"obj": {"name": "John", "age": 42}}`,
					path: "obj",
				},
			} {
				t.Run(item.name, func(t *testing.T) {
					r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
					r.body = []byte(item.body)

					r.Json(t, item.path, action.AssertExists())
				})
			}

		})
		t.Run("test AssertGt", func(t *testing.T) {
			t.Run("test valid", func(t *testing.T) {
				for _, item := range []struct {
					name   string
					body   string
					path   string
					action action.Action
				}{
					{
						name:   "test integer",
						body:   `{"object": {"name": 42}}`,
						path:   "object.name",
						action: action.AssertGt(41),
					},
					{
						name:   "test integer",
						body:   `{"object": {"name": -40}}`,
						path:   "object.name",
						action: action.AssertGt(-50),
					},
				} {
					t.Run(item.name, func(t *testing.T) {
						r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
						r.body = []byte(item.body)

						r.Json(t, item.path, item.action)
					})
				}
			})
			t.Run("test invalid", func(t *testing.T) {
				for _, item := range []struct {
					name          string
					body          string
					path          string
					action        action.Action
					expectMessage string
				}{
					{
						name:          "test integer",
						body:          `{"object": {"name": 41}}`,
						path:          "object.name",
						action:        action.AssertGt(42),
						expectMessage: "FAILED: `response.Json`, path: `object.name`, action: `AssertGt`: expected 41 to be greater than 42",
					},
					{
						name:          "test integer",
						body:          `{"object": {"name": -60}}`,
						path:          "object.name",
						action:        action.AssertGt(-50),
						expectMessage: "FAILED: `response.Json`, path: `object.name`, action: `AssertGt`: expected -60 to be greater than -50",
					},
				} {
					t.Run(item.name, func(t *testing.T) {
						r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
						r.body = []byte(item.body)
						m := mocks.NewTestingT(t)
						m.On("Errorf", mock.Anything, mock.MatchedBy(matchByStringContains(item.expectMessage))).Once()
						m.On("FailNow")

						r.Json(m, item.path, item.action)

						m.AssertNumberOfCalls(t, "Errorf", 1)
						m.AssertNumberOfCalls(t, "FailNow", 1)
						m.AssertExpectations(t)

					})
				}
			})
		})
		t.Run("test AssertGte", func(t *testing.T) {
			t.Run("test valid", func(t *testing.T) {
				for _, item := range []struct {
					name   string
					body   string
					path   string
					action action.Action
				}{
					{
						name:   "test integer",
						body:   `{"object": {"name": 42}}`,
						path:   "object.name",
						action: action.AssertGte(41),
					},
					{
						name:   "test integer",
						body:   `{"object": {"name": -40}}`,
						path:   "object.name",
						action: action.AssertGte(-50),
					},
					{
						name:   "test integer",
						body:   `{"object": {"name": -40}}`,
						path:   "object.name",
						action: action.AssertGte(-40),
					},
				} {
					t.Run(item.name, func(t *testing.T) {
						r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
						r.body = []byte(item.body)

						r.Json(t, item.path, item.action)
					})
				}
			})
			t.Run("test invalid", func(t *testing.T) {
				for _, item := range []struct {
					name          string
					body          string
					path          string
					action        action.Action
					expectMessage string
				}{
					{
						name:          "test integer",
						body:          `{"object": {"name": 41}}`,
						path:          "object.name",
						action:        action.AssertGte(42),
						expectMessage: "FAILED: `response.Json`, path: `object.name`, action: `AssertGte`: expected 41 to be greater than or equal to 42",
					},
					{
						name:          "test integer",
						body:          `{"object": {"name": -60}}`,
						path:          "object.name",
						action:        action.AssertGte(-50),
						expectMessage: "FAILED: `response.Json`, path: `object.name`, action: `AssertGte`: expected -60 to be greater than or equal to -50",
					},
				} {
					t.Run(item.name, func(t *testing.T) {
						r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
						r.body = []byte(item.body)
						m := mocks.NewTestingT(t)
						m.On("Errorf", mock.Anything, mock.MatchedBy(matchByStringContains(item.expectMessage))).Once()
						m.On("FailNow")

						r.Json(m, item.path, item.action)
						m.AssertNumberOfCalls(t, "Errorf", 1)
						m.AssertNumberOfCalls(t, "FailNow", 1)
						m.AssertExpectations(t)
					})
				}
			})

		})

		t.Run("test AssertLt", func(t *testing.T) {
			t.Run("test valid", func(t *testing.T) {
				for _, item := range []struct {
					name   string
					body   string
					path   string
					action action.Action
				}{
					{
						name:   "test integer",
						body:   `{"object": {"name": 42}}`,
						path:   "object.name",
						action: action.AssertLt(43),
					},
					{
						name:   "test negative integer",
						body:   `{"object": {"name": -40}}`,
						path:   "object.name",
						action: action.AssertLt(-30),
					},
					{
						name:   "test integer",
						body:   `{"object": {"name": -1}}`,
						path:   "object.name",
						action: action.AssertLt(40),
					},
				} {
					t.Run(item.name, func(t *testing.T) {
						r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
						r.body = []byte(item.body)

						r.Json(t, item.path, item.action)
					})
				}
			})
			t.Run("test invalid", func(t *testing.T) {
				for _, item := range []struct {
					name          string
					body          string
					path          string
					action        action.Action
					expectMessage string
				}{
					{
						name:          "test integer",
						body:          `{"object": {"name": 41}}`,
						path:          "object.name",
						action:        action.AssertLt(40),
						expectMessage: "FAILED: `response.Json`, path: `object.name`, action: `AssertLt`: expected 41 to be less than 40",
					},
					{
						name:          "test integer",
						body:          `{"object": {"name": -60}}`,
						path:          "object.name",
						action:        action.AssertLt(-70),
						expectMessage: "FAILED: `response.Json`, path: `object.name`, action: `AssertLt`: expected -60 to be less than -70",
					},
				} {
					t.Run(item.name, func(t *testing.T) {
						r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
						r.body = []byte(item.body)
						m := mocks.NewTestingT(t)
						m.On("Errorf", mock.Anything, mock.MatchedBy(matchByStringContains(item.expectMessage))).Once()
						m.On("FailNow")

						r.Json(m, item.path, item.action)

						m.AssertNumberOfCalls(t, "Errorf", 1)
						m.AssertNumberOfCalls(t, "FailNow", 1)
						m.AssertExpectations(t)
					})
				}
			})

		})

		t.Run("test AssertLte", func(t *testing.T) {
			t.Run("test valid", func(t *testing.T) {
				for _, item := range []struct {
					name   string
					body   string
					path   string
					action action.Action
				}{
					{
						name:   "test integer",
						body:   `{"object": {"name": 40}}`,
						path:   "object.name",
						action: action.AssertLte(41),
					},
					{
						name:   "test integer",
						body:   `{"object": {"name": -40}}`,
						path:   "object.name",
						action: action.AssertLte(-30),
					},
					{
						name:   "test integer",
						body:   `{"object": {"name": -40}}`,
						path:   "object.name",
						action: action.AssertLte(-40),
					},
				} {
					t.Run(item.name, func(t *testing.T) {
						r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
						r.body = []byte(item.body)

						r.Json(t, item.path, item.action)
					})
				}
			})
			t.Run("test invalid", func(t *testing.T) {
				for _, item := range []struct {
					name          string
					body          string
					path          string
					action        action.Action
					expectMessage string
				}{
					{
						name:          "test integer",
						body:          `{"object": {"name": 41}}`,
						path:          "object.name",
						action:        action.AssertLte(40),
						expectMessage: "FAILED: `response.Json`, path: `object.name`, action: `AssertLte`: expected 41 to be less than or equal to 40",
					},
					{
						name:          "test integer",
						body:          `{"object": {"name": -60}}`,
						path:          "object.name",
						action:        action.AssertLte(-70),
						expectMessage: "FAILED: `response.Json`, path: `object.name`, action: `AssertLte`: expected -60 to be less than or equal to -70",
					},
				} {
					t.Run(item.name, func(t *testing.T) {
						r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
						r.body = []byte(item.body)
						m := mocks.NewTestingT(t)
						m.On("Errorf", mock.Anything, mock.MatchedBy(matchByStringContains(item.expectMessage))).Once()
						m.On("FailNow")

						r.Json(m, item.path, item.action)
						m.AssertNumberOfCalls(t, "Errorf", 1)
						m.AssertNumberOfCalls(t, "FailNow", 1)
						m.AssertExpectations(t)
					})
				}
			})

		})
		t.Run("test AssertNotEquals", func(t *testing.T) {
			t.Run("test valid", func(t *testing.T) {
				for _, item := range []struct {
					name   string
					body   string
					path   string
					action action.Action
				}{
					{
						name:   "test integer",
						body:   `{"object": {"name": 40}}`,
						path:   "object.name",
						action: action.AssertNotEquals(41),
					},
					{
						name:   "test integer",
						body:   `{"object": {"name": -40}}`,
						path:   "object.name",
						action: action.AssertNotEquals(-50),
					},
					{
						name:   "test integer",
						body:   `{"object": {"name": -40}}`,
						path:   "object.name",
						action: action.AssertNotEquals(-120),
					},
				} {
					t.Run(item.name, func(t *testing.T) {
						r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
						r.body = []byte(item.body)

						r.Json(t, item.path, item.action)
					})
				}
			})
			t.Run("test invalid", func(t *testing.T) {
				for _, item := range []struct {
					name          string
					body          string
					path          string
					action        action.Action
					expectMessage string
				}{
					{
						name:          "test integer",
						body:          `{"object": {"name": 41}}`,
						path:          "object.name",
						action:        action.AssertNotEquals(41),
						expectMessage: "FAILED: `response.Json`, path: `object.name`, action: `AssertNotEquals`: expected value to not be equal to 41",
					},
					{
						name:          "test integer",
						body:          `{"object": {"name": -60}}`,
						path:          "object.name",
						action:        action.AssertNotEquals(-60),
						expectMessage: "FAILED: `response.Json`, path: `object.name`, action: `AssertNotEquals`: expected value to not be equal to -60",
					},
					{
						name:          "test integer",
						body:          `{"object": {"name": -60}}`,
						path:          "object.name",
						action:        action.AssertNotEquals(json.RawMessage(`{}`)),
						expectMessage: "FAILED: `response.Json`, path: `object.name`, action: `AssertNotEquals`: expected value to not be equal to {}",
					},
				} {
					t.Run(item.name, func(t *testing.T) {
						r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
						r.body = []byte(item.body)
						m := mocks.NewTestingT(t)
						m.On("Errorf", mock.Anything, mock.MatchedBy(matchByStringContains(item.expectMessage))).Once()
						m.On("FailNow")

						r.Json(m, item.path, item.action)
						m.AssertNumberOfCalls(t, "Errorf", 1)
						m.AssertNumberOfCalls(t, "FailNow", 1)
						m.AssertExpectations(t)
					})
				}
			})
		})
	})

	t.Run("test invalid path or wrong type", func(t *testing.T) {
		for _, item := range []struct {
			name          string
			body          string
			path          string
			action        action.Action
			expectedError string
		}{
			{
				name:          "wrong key in object",
				body:          `{"name": "John"}`,
				path:          "what",
				action:        action.AssertEquals("John"),
				expectedError: "FAILED: `response.Json`, path: `what`, not present",
			},
			{
				name:          "wrong path",
				body:          `{"name": "John"}`,
				path:          "0",
				action:        action.AssertEquals("John"),
				expectedError: "FAILED: `response.Json`, path: `0`, cannot unmarshal object into slice",
			},
			{
				name:          "wrong type",
				body:          `{"name": "John"}`,
				path:          "name",
				action:        action.AssertEquals(ptrTo(42)),
				expectedError: "FAILED: `response.Json`, path: `name`, unmarshal error: json: cannot unmarshal string into Go value of type int",
			},
			{
				name:          "wrong path",
				body:          `[{"name": "John"}]`,
				path:          "name",
				action:        action.AssertEquals(ptrTo(42)),
				expectedError: "FAILED: `response.Json`, path: `name`, cannot unmarshal action",
			},
		} {
			t.Run(item.name, func(t *testing.T) {
				r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
				r.body = []byte(item.body)

				m := mocks.NewTestingT(t)
				m.On("Errorf", mock.Anything, mock.MatchedBy(matchByStringContains(item.expectedError))).Once()
				m.On("FailNow")

				r.Json(m, item.path, item.action)
				m.AssertNumberOfCalls(t, "Errorf", 1)
				m.AssertNumberOfCalls(t, "FailNow", 1)
				m.AssertExpectations(t)
			})
		}
	})
}

func TestResponse_AssertStatus(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusCreated, http.StatusNoContent} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			rw := httptest.NewRecorder()
			rw.WriteHeader(status)
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			r := newResponse(rw, req)
			r.Status(t, action.AssertEquals(status))
		})
	}
}

func TestResponse_Unmarshal(t *testing.T) {
	type User struct {
		Name string `json:"name"`
	}

	for _, item := range []struct {
		name     string
		body     string
		expected User
	}{
		{
			name:     "equal strings",
			body:     `{"name": "John"}`,
			expected: User{Name: "John"},
		},
		{
			name:     "equal strings",
			body:     `{"hello": "world", "name": "John"}`,
			expected: User{Name: "John"},
		},
	} {
		t.Run(item.name, func(t *testing.T) {
			rw := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			r := newResponse(rw, req)
			r.body = []byte(item.body)
			var user User
			r.Json(t, "", action.Unmarshal(&user))
			assert.Equal(t, item.expected, user)
		})
	}

}
