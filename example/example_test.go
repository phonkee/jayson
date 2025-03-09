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

package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/phonkee/jayson/tester"
	"github.com/phonkee/jayson/tester/resolver"
	"net/http"
	"testing"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"users": [{"id":1,"name":"John Doe"}]}`))
}

func TestTester(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/users", ListUsers).Methods(http.MethodGet).Name("api:v1:users:list")

	tester.WithHttpServer(t, context.Background(), router, func(t *testing.T, ctx context.Context, address string) {
		// test api
		tester.WithAPI(t, &tester.Deps{
			Resolver: resolver.NewGorillaResolver(t, router),
			Address:  address,
		}, func(api tester.APIClient) {
			api.Get(t, api.ReverseURL(t, "api:v1:users:list")).
				Do(t, ctx).
				AssertStatus(t, http.StatusOK).
				//AssertJsonPath2(t, "users", tester.AssertExists(true)).
				//AssertJsonPath2(t, "users", tester.AssertLen(1)).
				//AssertJsonPath2(t, "users.0.id", tester.AssertGte(1)).
				//AssertJsonPath2(t, "users.0.id", tester.AssertGt(0)).
				//AssertJsonPath2(t, "users.0.id", tester.AssertLte(10)).
				//AssertJsonPath2(t, "users.0.id", tester.AssertLt(10)).
				//AssertJsonPath2(t, "users.0", json.RawMessage(`{"id":1,"name":"John Doe"}`)).
				//AssertJsonPath2(t, "users.0.name", tester.AssertNotEqual("John Doe")).
				AssertJsonPath2(t, "users.0.name", tester.AssertEqual("John Doe")).
				AssertJsonPath2(t, "users.0.name", "John Doe")
		})
	})
}
