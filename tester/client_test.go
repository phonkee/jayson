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
	"github.com/phonkee/jayson/tester/resolver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

// exampleResponse is a response struct for testing
type exampleResponse struct {
	Status string `json:"status"`
	Host   string `json:"host"`
}

// exampleHandler is a handler for testing
func exampleHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(exampleResponse{
		Status: "something",
		Host:   "localhost",
	}); err != nil {
		panic(err)
	}
}

func newHealthRouter(t require.TestingT) *mux.Router {
	router := mux.NewRouter()
	assert.NotNil(t, router)
	router.HandleFunc("/api/v1/health", exampleHandler).Methods(http.MethodGet).Name("api:v1:health")
	router.HandleFunc("/api/v1/health/{component}", exampleHandler).Methods(http.MethodGet).Name("api:v1:health:extra")
	return router
}

func TestClient(t *testing.T) {
	t.Run("test handler", func(t *testing.T) {
		router := newHealthRouter(t)
		tester.WithAPI(t, &tester.Deps{
			Resolver: resolver.NewGorillaResolver(t, router),
			Handler:  router,
		}, func(api tester.APIClient) {
			// context first
			ctx := context.Background()

			var (
				host   string
				status string
			)

			// response struct
			rr := exampleResponse{}

			// do response
			api.Request(t, http.MethodGet, api.ReverseURL(t, "api:v1:health")).
				Do(t, ctx).
				Status(t, http.StatusOK).
				AssertJsonEquals(t, `{"status": "something", "host": "localhost"}`).
				Unmarshal(t,
					tester.APIObject(t,
						"status", &status,
						"host", &host,
					),
				).
				Unmarshal(t, &rr)

			assert.Equal(t, "something", status)
			assert.Equal(t, "localhost", host)
		})
	})

	t.Run("test reverse url", func(t *testing.T) {
		router := newHealthRouter(t)
		tester.WithAPI(t, &tester.Deps{
			Resolver: resolver.NewGorillaResolver(t, router),
			Handler:  router,
		}, func(api tester.APIClient) {
			assert.Equal(t,
				"/api/v1/health",
				api.ReverseURL(t, "api:v1:health"),
			)
			assert.Equal(t,
				"/api/v1/health/database?page=1",
				api.ReverseURL(t,
					"api:v1:health:extra",
					api.ReverseArgs(t, "component", "database"),
					api.ReverseQuery(t, "page", "1"),
				),
			)
		})
	})

	t.Run("test address", func(t *testing.T) {
		t.Run("test error", func(t *testing.T) {
			// create router so we have a handler to run server
			router := newHealthRouter(t)

			tester.WithHttpServer(t, context.Background(), router, func(t *testing.T, ctx context.Context, address string) {
				tester.WithAPI(t, &tester.Deps{
					Resolver: resolver.NewGorillaResolver(t, router),
					Address:  address,
				}, func(api tester.APIClient) {
					// context first
					ctx, cf := context.WithTimeout(context.Background(), time.Second*2)
					defer cf()

					// do response
					api.Request(t, http.MethodGet, "/not/exist").
						Do(t, ctx).
						Status(t, http.StatusNotFound)
				})
			})
		})

		t.Run("test success", func(t *testing.T) {
			router := newHealthRouter(t)
			tester.WithHttpServer(t, context.Background(), router, func(t *testing.T, ctx context.Context, address string) {
				tester.WithAPI(t, &tester.Deps{
					Resolver: resolver.NewGorillaResolver(t, router),
					Address:  address,
				}, func(api tester.APIClient) {
					// context first
					ctx, cf := context.WithTimeout(context.Background(), time.Second*2)
					defer cf()

					var (
						host   string
						status string
					)

					// response struct
					rr := exampleResponse{}

					// do response
					api.Request(t, http.MethodGet, api.ReverseURL(t, "api:v1:health")).
						Do(t, ctx).
						Status(t, http.StatusOK).
						AssertJsonEquals(t, `{"status": "something", "host": "localhost"}`).
						AssertJsonPath(t, "status", "something").
						AssertJsonPath(t, "__len__", 2).
						AssertJsonPath(t, "__keys__", []string{"status", "host"}).
						Unmarshal(t,
							tester.APIObject(t,
								"status", &status,
								"host", &host,
							),
						).
						Unmarshal(t, &rr)

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
			// action round trip
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
