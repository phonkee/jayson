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
	"fmt"
	"github.com/phonkee/jayson/tester/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestResponse_AssertHeaderValue(t *testing.T) {
	t.Run("test contains", func(t *testing.T) {
		for _, item := range []struct {
			name    string
			headers http.Header
			key     string
			value   string
		}{
			{
				name:    "header exists",
				headers: map[string][]string{"Hello": {"World"}},
				key:     "Hello",
				value:   "World",
			},
			{
				name:    "header does not exist",
				headers: map[string][]string{"Hello": {"This", "World"}},
				key:     "Hello",
				value:   "World",
			},
		} {
			t.Run(item.name, func(t *testing.T) {
				rw := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				for k, v := range item.headers {
					rw.Header()[k] = v
				}
				r := newResponse(rw, req)
				r.AssertHeaderValue(t, item.key, item.value)
			})
		}
	})

	t.Run("test does not contain", func(t *testing.T) {
		for _, item := range []struct {
			name    string
			headers http.Header
			key     string
			value   string
		}{
			{
				name:    "header exists",
				headers: map[string][]string{"Hello": {"World"}},
				key:     "Hello",
				value:   "Other",
			},
			{
				name:    "header does not exist",
				headers: map[string][]string{"Hello": {"This", "World"}},
				key:     "Nope",
				value:   "World",
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
				m.On("Errorf", mock.Anything, mock.MatchedBy(func(in string) bool {
					return strings.Contains(in, fmt.Sprintf("header `%v` not found or does not have value", item.key))
				})).Return()
				m.On("FailNow").Once()
				r.AssertHeaderValue(m, item.key, item.value)
				m.AssertExpectations(t)
				m.AssertNumberOfCalls(t, "Errorf", 1)
				m.AssertNumberOfCalls(t, "FailNow", 1)
			})
		}
	})

}

func TestResponse_AssertJsonEquals(t *testing.T) {
	type User struct {
		Name string `json:"name"`
	}

	for _, item := range []struct {
		name     string
		body     string
		expected any
	}{
		{
			name:     "equal strings",
			body:     `{"name": "John"}`,
			expected: `{"name": "John"}`,
		},
		{
			name:     "equal bytes",
			body:     `{"hello": "world", "name": "John"}`,
			expected: []byte(`{"name": "John", "hello": "world"}`),
		},
		{
			name:     "equal struct",
			body:     `{"name": "John"}`,
			expected: User{Name: "John"},
		},
	} {
		t.Run(item.name, func(t *testing.T) {
			rw := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			r := newResponse(rw, req)
			r.body = []byte(item.body)
			r.AssertJsonEquals(t, item.expected)
		})
	}
}

func TestResponse_AssertJsonPath(t *testing.T) {

	t.Run("test valid path", func(t *testing.T) {
		for _, item := range []struct {
			name     string
			body     string
			path     string
			expected any
		}{
			{
				name:     "simple strings",
				body:     `{"name": "John"}`,
				path:     "name",
				expected: "John",
			},
			{
				name:     "whole object as RawMessage",
				body:     `{"name": "John"}`,
				path:     ".",
				expected: json.RawMessage(`{"name": "John"}`),
			},
			{
				name:     "simple strings to pointer",
				body:     `{"name": "John"}`,
				path:     "name",
				expected: ptrTo("John"),
			},
			{
				name:     "test embedded struct strings",
				body:     `{"other": {"name": "John"}}`,
				path:     "other.name",
				expected: "John",
			},
			{
				name:     "test slice of objects to string",
				body:     `{"other": [{"name": "John"}, {"name": "Doe"}]}`,
				path:     "other.1.name",
				expected: "Doe",
			},
			{
				name:     "test slice of objects to object",
				body:     `{"other": [{"name": "John"}, {"name": "Doe"}, {"object": {"value": "Doe"}}]}`,
				path:     "other.2.object",
				expected: testStruct{Value: "Doe"},
			},
			{
				name:     "test slice of objects to pointer to object",
				body:     `{"other": [{"name": "John"}, {"name": "Doe"}, {"object": {"value": "Doe"}}]}`,
				path:     "other.2.object",
				expected: &testStruct{Value: "Doe"},
			},
			{
				name:     "test slice of objects to pointer to integer",
				body:     `{"other": [{"name": "John"}, {"name": "Doe"}, {"object": {"value": 42}}]}`,
				path:     "other.2.object.value.__eq__",
				expected: 42,
			},
			{
				name:     "test slice of objects to pointer to integer (with dots)",
				body:     `{"other": [{"name": "John"}, {"name": "Doe"}, {"object": {"value": 42}}]}`,
				path:     "other.2......object.value",
				expected: ptrTo(42),
			},
			{
				name:     "test raw json",
				body:     `{"other": [{"name": "John"}, {"name": "Doe"}, {"object": {"value": 42, "other": 12}}]}`,
				path:     "other.2",
				expected: json.RawMessage(`{"object": {"other": 12, "value": 42}}`),
			},
			{
				name: "test raw json",
				body: `{"other": [{"value": "John"}, {"value": "Doe"}, {"value": "Mark", "extra": "extra"}]}`,
				path: "other",
				expected: []testStruct{
					{Value: "John"},
					{Value: "Doe"},
					{Value: "Mark"},
				},
			},
		} {
			t.Run(item.name, func(t *testing.T) {
				r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
				r.body = []byte(item.body)

				r.AssertJsonPath(t, item.path, item.expected)
			})
		}
	})

	t.Run("test special operations", func(t *testing.T) {
		t.Run("test special operation: __len__", func(t *testing.T) {
			t.Run("valid path", func(t *testing.T) {
				for _, item := range []struct {
					name   string
					body   string
					path   string
					expect any
				}{
					{
						name:   "test len of array",
						body:   `{"other": [{"name": "John"}, {"name": "Doe"}, {"object": {"value": 42, "other": 12}}]}`,
						path:   "other.__len__",
						expect: 3,
					},
					{
						name:   "test len of array",
						body:   `{"other": [{"name": "John"}, {"name": "Doe"}, {"object": {"value": 42, "other": 12}}]}`,
						path:   "other.2.object.__len__",
						expect: 2,
					},
				} {
					t.Run(item.name, func(t *testing.T) {
						r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
						r.body = []byte(item.body)

						r.AssertJsonPath(t, item.path, item.expect)
					})
				}
			})
			t.Run("invalid path", func(t *testing.T) {
				for _, item := range []struct {
					name   string
					body   string
					path   string
					expect string
				}{
					{
						name:   "test len of string",
						body:   `{"other": [{"name": "John"}, {"name": "Doe"}, {"object": {"value": 42, "other": 12}}]}`,
						path:   "other.0.name.__len__",
						expect: "path: `other.0.name.__len__`, __len__ is only supported for arrays and objects",
					},
					{
						name:   "test len of int",
						body:   `{"other": [{"name": "John"}, {"name": "Doe"}, {"object": {"value": 42, "other": 12}}]}`,
						path:   "other.2.object.other.__len__",
						expect: "path: `other.2.object.other.__len__`, __len__ is only supported for arrays and objects",
					},
				} {
					t.Run(item.name, func(t *testing.T) {
						r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
						r.body = []byte(item.body)

						m := mocks.NewTestingT(t)
						m.On("Errorf", mock.Anything, mock.MatchedBy(matchByStringContains(item.expect))).Once()
						m.On("FailNow").Run(func(args mock.Arguments) {
							t.Skip()
						})

						r.AssertJsonPath(m, item.path, 1)
					})
				}

			})
		})
		t.Run("test special operation: __keys__", func(t *testing.T) {
			for _, item := range []struct {
				name   string
				body   string
				path   string
				expect any
			}{
				{
					name:   "test keys of object",
					body:   `{"name": "John", "age": 42}`,
					path:   "__keys__",
					expect: []string{"name", "age"},
				},
				{
					name:   "test keys of object",
					body:   `{"obj": {"name": "John", "age": 42}}`,
					path:   "obj.__keys__",
					expect: []string{"name", "age"},
				},
			} {
				t.Run(item.name, func(t *testing.T) {
					r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
					r.body = []byte(item.body)

					r.AssertJsonPath(t, item.path, item.expect)
				})
			}

		})
		t.Run("test special operation: __exists__", func(t *testing.T) {
			for _, item := range []struct {
				name string
				body string
				path string
			}{
				{
					name: "test keys of object",
					body: `{"name": "John", "age": 42}`,
					path: "name.__exists__",
				},
				{
					name: "test keys of object",
					body: `{"obj": {"name": "John", "age": 42}}`,
					path: "obj.__exists__",
				},
			} {
				t.Run(item.name, func(t *testing.T) {
					r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
					r.body = []byte(item.body)

					r.AssertJsonPath(t, item.path, nil)
				})
			}

		})
		t.Run("test special operation: __gt__", func(t *testing.T) {
			t.Run("test valid", func(t *testing.T) {
				for _, item := range []struct {
					name   string
					body   string
					path   string
					expect any
				}{
					{
						name:   "test integer",
						body:   `{"object": {"name": 42}}`,
						path:   "object.name.__gt__",
						expect: 41,
					},
					{
						name:   "test integer",
						body:   `{"object": {"name": -40}}`,
						path:   "object.name.__gt__",
						expect: -50,
					},
				} {
					t.Run(item.name, func(t *testing.T) {
						r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
						r.body = []byte(item.body)

						r.AssertJsonPath(t, item.path, item.expect)
					})
				}
			})
			t.Run("test invalid", func(t *testing.T) {
				for _, item := range []struct {
					name          string
					body          string
					path          string
					expect        any
					expectMessage string
				}{
					{
						name:          "test integer",
						body:          `{"object": {"name": 41}}`,
						path:          "object.name.__gt__",
						expect:        42,
						expectMessage: "value `41` is not greater than `42`",
					},
					{
						name:          "test integer",
						body:          `{"object": {"name": -60}}`,
						path:          "object.name.__gt__",
						expect:        -50,
						expectMessage: "value `-60` is not greater than `-50`",
					},
				} {
					t.Run(item.name, func(t *testing.T) {
						r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
						r.body = []byte(item.body)
						m := mocks.NewTestingT(t)
						m.On("Errorf", mock.Anything, mock.MatchedBy(matchByStringContains(item.expectMessage))).Once()

						r.AssertJsonPath(m, item.path, item.expect)
					})
				}
			})
		})
		t.Run("test special operation: __gte__", func(t *testing.T) {
			t.Run("test valid", func(t *testing.T) {
				for _, item := range []struct {
					name   string
					body   string
					path   string
					expect any
				}{
					{
						name:   "test integer",
						body:   `{"object": {"name": 42}}`,
						path:   "object.name.__gte__",
						expect: 41,
					},
					{
						name:   "test integer",
						body:   `{"object": {"name": -40}}`,
						path:   "object.name.__gte__",
						expect: -50,
					},
					{
						name:   "test integer",
						body:   `{"object": {"name": -40}}`,
						path:   "object.name.__gte__",
						expect: -40,
					},
				} {
					t.Run(item.name, func(t *testing.T) {
						r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
						r.body = []byte(item.body)

						r.AssertJsonPath(t, item.path, item.expect)
					})
				}
			})
			t.Run("test invalid", func(t *testing.T) {
				for _, item := range []struct {
					name          string
					body          string
					path          string
					expect        any
					expectMessage string
				}{
					{
						name:          "test integer",
						body:          `{"object": {"name": 41}}`,
						path:          "object.name.__gte__",
						expect:        42,
						expectMessage: "value `41` is not greater than or equal `42`",
					},
					{
						name:          "test integer",
						body:          `{"object": {"name": -60}}`,
						path:          "object.name.__gte__",
						expect:        -50,
						expectMessage: "value `-60` is not greater than or equal `-50`",
					},
				} {
					t.Run(item.name, func(t *testing.T) {
						r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
						r.body = []byte(item.body)
						m := mocks.NewTestingT(t)
						m.On("Errorf", mock.Anything, mock.MatchedBy(matchByStringContains(item.expectMessage))).Once()

						r.AssertJsonPath(m, item.path, item.expect)
					})
				}
			})

		})

		t.Run("test special operation: __lt__", func(t *testing.T) {
			t.Run("test valid", func(t *testing.T) {
				for _, item := range []struct {
					name   string
					body   string
					path   string
					expect any
				}{
					{
						name:   "test integer",
						body:   `{"object": {"name": 42}}`,
						path:   "object.name.__lt__",
						expect: 43,
					},
					{
						name:   "test integer",
						body:   `{"object": {"name": -40}}`,
						path:   "object.name.__lt__",
						expect: -30,
					},
					{
						name:   "test integer",
						body:   `{"object": {"name": 0}}`,
						path:   "object.name.__lt__",
						expect: 40,
					},
				} {
					t.Run(item.name, func(t *testing.T) {
						r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
						r.body = []byte(item.body)

						r.AssertJsonPath(t, item.path, item.expect)
					})
				}
			})
			t.Run("test invalid", func(t *testing.T) {
				for _, item := range []struct {
					name          string
					body          string
					path          string
					expect        any
					expectMessage string
				}{
					{
						name:          "test integer",
						body:          `{"object": {"name": 41}}`,
						path:          "object.name.__lt__",
						expect:        40,
						expectMessage: "value `41` is not less than `40`",
					},
					{
						name:          "test integer",
						body:          `{"object": {"name": -60}}`,
						path:          "object.name.__lt__",
						expect:        -70,
						expectMessage: "value `-60` is not less than `-70`",
					},
				} {
					t.Run(item.name, func(t *testing.T) {
						r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
						r.body = []byte(item.body)
						m := mocks.NewTestingT(t)
						m.On("Errorf", mock.Anything, mock.MatchedBy(matchByStringContains(item.expectMessage))).Once()

						r.AssertJsonPath(m, item.path, item.expect)
					})
				}
			})

		})

		t.Run("test special operation: __lte__", func(t *testing.T) {
			t.Run("test valid", func(t *testing.T) {
				for _, item := range []struct {
					name   string
					body   string
					path   string
					expect any
				}{
					{
						name:   "test integer",
						body:   `{"object": {"name": 40}}`,
						path:   "object.name.__lte__",
						expect: 41,
					},
					{
						name:   "test integer",
						body:   `{"object": {"name": -40}}`,
						path:   "object.name.__lte__",
						expect: -30,
					},
					{
						name:   "test integer",
						body:   `{"object": {"name": -40}}`,
						path:   "object.name.__lte__",
						expect: -40,
					},
				} {
					t.Run(item.name, func(t *testing.T) {
						r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
						r.body = []byte(item.body)

						r.AssertJsonPath(t, item.path, item.expect)
					})
				}
			})
			t.Run("test invalid", func(t *testing.T) {
				for _, item := range []struct {
					name          string
					body          string
					path          string
					expect        any
					expectMessage string
				}{
					{
						name:          "test integer",
						body:          `{"object": {"name": 41}}`,
						path:          "object.name.__lte__",
						expect:        40,
						expectMessage: "value `41` is not less than or equal `40`",
					},
					{
						name:          "test integer",
						body:          `{"object": {"name": -60}}`,
						path:          "object.name.__lte__",
						expect:        -70,
						expectMessage: "value `-60` is not less than or equal `-70`",
					},
				} {
					t.Run(item.name, func(t *testing.T) {
						r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
						r.body = []byte(item.body)
						m := mocks.NewTestingT(t)
						m.On("Errorf", mock.Anything, mock.MatchedBy(matchByStringContains(item.expectMessage))).Once()

						r.AssertJsonPath(m, item.path, item.expect)
					})
				}
			})

		})
		t.Run("test special operation: __neq__", func(t *testing.T) {
			t.Run("test valid", func(t *testing.T) {
				for _, item := range []struct {
					name   string
					body   string
					path   string
					expect any
				}{
					{
						name:   "test integer",
						body:   `{"object": {"name": 40}}`,
						path:   "object.name.__neq__",
						expect: 41,
					},
					{
						name:   "test integer",
						body:   `{"object": {"name": -40}}`,
						path:   "object.name.__neq__",
						expect: -50,
					},
					{
						name:   "test integer",
						body:   `{"object": {"name": -40}}`,
						path:   "object.name.__neq__",
						expect: -120,
					},
				} {
					t.Run(item.name, func(t *testing.T) {
						r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
						r.body = []byte(item.body)

						r.AssertJsonPath(t, item.path, item.expect)
					})
				}
			})
			t.Run("test invalid", func(t *testing.T) {
				for _, item := range []struct {
					name          string
					body          string
					path          string
					expect        any
					expectMessage string
				}{
					{
						name:          "test integer",
						body:          `{"object": {"name": 41}}`,
						path:          "object.name.__neq__",
						expect:        41,
						expectMessage: "value `41`, should not equal to `41`, but it did",
					},
					{
						name:          "test integer",
						body:          `{"object": {"name": -60}}`,
						path:          "object.name.__neq__",
						expect:        -60,
						expectMessage: "value `-60`, should not equal to `-60`, but it did",
					},
				} {
					t.Run(item.name, func(t *testing.T) {
						r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
						r.body = []byte(item.body)
						m := mocks.NewTestingT(t)
						m.On("Errorf", mock.Anything, mock.MatchedBy(matchByStringContains(item.expectMessage))).Once()

						r.AssertJsonPath(m, item.path, item.expect)
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
			into          any
			expectedError string
		}{
			{
				name:          "wrong key in object",
				body:          `{"name": "John"}`,
				path:          "what",
				expectedError: "key `what` in path `what` not found",
			},
			{
				name:          "wrong path",
				body:          `{"name": "John"}`,
				path:          "0",
				expectedError: "failed to unmarshal array `0` into",
			},
			{
				name:          "wrong type",
				body:          `{"name": "John"}`,
				path:          "name",
				into:          ptrTo(42),
				expectedError: "failed to unmarshal `name` into `*int`",
			},
			{
				name:          "wrong path",
				body:          `[{"name": "John"}]`,
				path:          "name",
				into:          ptrTo(42),
				expectedError: "failed to unmarshal `name` into `*int`",
			},
		} {
			t.Run(item.name, func(t *testing.T) {
				r := newResponse(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
				r.body = []byte(item.body)

				m := mocks.NewTestingT(t)
				m.On("Errorf", mock.Anything, mock.MatchedBy(matchByStringContains(item.expectedError))).Once()
				m.On("FailNow").Run(func(args mock.Arguments) {
					// weird but this is the only way
					t.Skip()
				})

				r.AssertJsonPath(m, item.path, item.into)
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
			r.AssertStatus(t, status)
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
			r.Unmarshal(t, &user)
			assert.Equal(t, item.expected, user)
		})
	}

}
