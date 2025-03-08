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
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/phonkee/jayson/tester"
	"github.com/phonkee/jayson/tester/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"strings"
	"testing"
	"time"
)

// exampleResponse is a response struct for testing
type exampleResponse struct {
	Status string `json:"status"`
	Host   string `json:"host"`
}

// exampleHandler is a handler for testing
func exampleHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(exampleResponse{
		Status: "something",
		Host:   "localhost",
	}); err != nil {
		panic(err)
	}
}

// MatchByStringContains matches string by substring
func matchByStringContains(s string) func(in string) bool {
	return func(in string) bool {
		return strings.Contains(in, s)
	}
}

func newHealthRouter(t *testing.T) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/health", exampleHandler).Methods(http.MethodGet).Name("api:v1:health")
	return router
}

// ptrTo helper
func ptrTo[T any](v T) *T {
	return &v
}

func TestClient(t *testing.T) {
	t.Run("test handler", func(t *testing.T) {
		router := newHealthRouter(t)
		tester.WithAPI(t, &tester.Deps{
			Router:  router,
			Handler: router,
		}, func(api tester.APIClient) {
			// context first
			ctx := context.Background()

			var (
				host   string
				status string
			)

			// test pointer in AssertJsonKeyEquals
			ptrStatus := ptrTo("something")
			statusValue := "something"

			// response struct
			rr := exampleResponse{}

			// do response
			api.Request(t, http.MethodGet, api.ReverseURL(t, "api:v1:health")).
				Do(t, ctx).
				AssertStatus(t, http.StatusOK).
				AssertJsonEquals(t, `{"status": "something", "host": "localhost"}`).
				Unmarshal(t,
					tester.APIObject(t,
						"status", &status,
						"host", &host,
					),
				).
				Unmarshal(t, &rr).
				AssertJsonKeyEquals(t, "status", "something").
				AssertJsonKeyEquals(t, "status", statusValue).
				AssertJsonKeyEquals(t, "status", ptrStatus).
				AssertJsonKeyEquals(t, "host", "localhost")

			assert.Equal(t, "something", status)
			assert.Equal(t, "localhost", host)
		})
	})

	t.Run("test address", func(t *testing.T) {
		// create router so we have a handler to run server
		router := newHealthRouter(t)

		t.Run("test error", func(t *testing.T) {
			tester.WithHttpServer(t, router, func(t *testing.T, address string) {
				tester.WithAPI(t, &tester.Deps{
					Router:  router,
					Address: address,
				}, func(api tester.APIClient) {
					// context first
					ctx, cf := context.WithTimeout(context.Background(), time.Second*2)
					defer cf()

					// do response
					api.Request(t, http.MethodGet, "/not/exist").
						Do(t, ctx).
						AssertStatus(t, http.StatusNotFound)
				})
			})
		})

		t.Run("test success", func(t *testing.T) {
			tester.WithHttpServer(t, router, func(t *testing.T, address string) {
				tester.WithAPI(t, &tester.Deps{
					Router:  router,
					Address: address,
				}, func(api tester.APIClient) {
					// context first
					ctx, cf := context.WithTimeout(context.Background(), time.Second*2)
					defer cf()

					var (
						host   string
						status string
					)

					// test pointer in AssertJsonKeyEquals
					ptrStatus := ptrTo("something")
					statusValue := "something"

					// response struct
					rr := exampleResponse{}

					// do response
					api.Request(t, http.MethodGet, api.ReverseURL(t, "api:v1:health")).
						Do(t, ctx).
						AssertStatus(t, http.StatusOK).
						AssertJsonEquals(t, `{"status": "something", "host": "localhost"}`).
						Unmarshal(t,
							tester.APIObject(t,
								"status", &status,
								"host", &host,
							),
						).
						Unmarshal(t, &rr).
						AssertJsonKeyEquals(t, "status", "something").
						AssertJsonKeyEquals(t, "status", statusValue).
						AssertJsonKeyEquals(t, "status", ptrStatus).
						AssertJsonKeyEquals(t, "host", "localhost")

					assert.Equal(t, "something", status)
					assert.Equal(t, "localhost", host)
				})

			})
		})

	})

}

func TestClient_MethodAliases(t *testing.T) {
	testMethod := func(t *testing.T, method string, fn func(client tester.APIClient) func(t require.TestingT, path string) tester.APIRequest) {
		name := "test " + method
		t.Run(name, func(t *testing.T) {
			rt := mocks.NewRoundTripper(t)
			resp := &http.Response{
				StatusCode: http.StatusOK,
			}
			// expect round trip
			rt.On("RoundTrip", mock.MatchedBy(func(r *http.Request) bool {
				return r.Method == method
			})).Return(resp, nil)
			hc := &http.Client{
				Transport: rt,
			}

			tester.WithAPI(t, &tester.Deps{
				Address: "http://localhost",
				Client:  hc,
			}, func(api tester.APIClient) {
				// context first
				ctx := context.Background()

				fn(api)(t, "/api/v1/health").Do(t, ctx)
			})

		})
	}

	testMethod(t, http.MethodDelete, func(client tester.APIClient) func(t require.TestingT, path string) tester.APIRequest {
		return client.Delete
	})

	testMethod(t, http.MethodGet, func(client tester.APIClient) func(t require.TestingT, path string) tester.APIRequest {
		return client.Get
	})

	testMethod(t, http.MethodPost, func(client tester.APIClient) func(t require.TestingT, path string) tester.APIRequest {
		return client.Post
	})

	testMethod(t, http.MethodPut, func(client tester.APIClient) func(t require.TestingT, path string) tester.APIRequest {
		return client.Put
	})

}

func TestClient_ReverseURL(t *testing.T) {
	t.Run("test missing router", func(t *testing.T) {
		tester.WithAPI(t, &tester.Deps{
			Handler: http.HandlerFunc(exampleHandler),
		}, func(api tester.APIClient) {
			m := mocks.NewTestingT(t)
			m.On("Errorf", mock.Anything, mock.MatchedBy(matchByStringContains("Deps: Router is nil")))
			m.On("FailNow").Run(func(args mock.Arguments) {
				t.Skip()
			})
			api.ReverseURL(m, "api:v1:health")
		})

	})

	t.Run("test missing route", func(t *testing.T) {
		router := mux.NewRouter()
		tester.WithAPI(t, &tester.Deps{
			Router:  router,
			Handler: http.HandlerFunc(exampleHandler),
		}, func(api tester.APIClient) {
			m := mocks.NewTestingT(t)
			m.On("Errorf", mock.Anything, mock.MatchedBy(matchByStringContains("route `api:v1:health` not found")))
			m.On("FailNow").Run(func(args mock.Arguments) {
				t.Skip()
			})
			api.ReverseURL(m, "api:v1:health")
		})

	})

}
