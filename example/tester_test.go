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
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/phonkee/jayson/tester"
	"github.com/phonkee/jayson/tester/action"
	"github.com/phonkee/jayson/tester/resolver"
	"github.com/stretchr/testify/assert"
	"net/http"
	"regexp"
	"testing"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"users": [{"id":1,"name":"John Doe", "admin": false}]}`))
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
			user := User{}
			userObj := User{}
			api.Get(t, api.ReverseURL(t, "api:v1:users:list")).
				Do(t, ctx).
				Status(t, action.AssertEquals(http.StatusOK)).
				Json(t, "users.0.name", action.AssertNot(action.AssertEquals("Johnson Doe"))).
				Json(t, "users.0.name", action.AssertEquals("John Doe")).
				Json(t, "users.0", action.AssertEquals(json.RawMessage(`{"id":1,"name":"John Doe"}`))).
				Json(t, "users.0.name", action.AssertIn("John Doe", "Peter Vrba")).
				Json(t, "users.0.name", action.AssertNot(action.AssertIn("Johnson Doe", "Peter Vrba"))).
				Json(t, "users", action.AssertExists()).
				Json(t, "user", action.AssertNot(action.AssertExists())).
				Json(t, "users", action.AssertLen(1)).
				Json(t, "users.0", action.AssertKeys("id", "name", "admin")).
				Json(t, "users.0.id", action.AssertGte(1)).
				Json(t, "users.0.id", action.AssertGt(0)).
				Json(t, "users.0.id", action.AssertLt(2)).
				Json(t, "users.0.id", action.AssertLte(1)).
				Json(t, "users.0", action.Unmarshal(&user)).
				Json(t, "users.0.id", action.AssertAll(
					action.AssertGte(0),
					action.AssertLte(1),
					action.AssertExists(),
				)).
				Json(t, "users.0.id", action.AssertAny(
					action.AssertGte(0),
					action.AssertLte(0),
					action.AssertNot(action.AssertExists()),
				)).
				Json(t, "users.0.id", action.AssertRegexMatch(
					regexp.MustCompile(`\d+`),
				)).
				Json(t, "users.0", action.UnmarshalObjectKeys(action.KV{
					"id":   &userObj.ID,
					"name": &userObj.Name,
				})).
				Json(t, "users.0.admin", action.AssertZero()).
				Json(t, "users.0.name", action.AssertNot(action.AssertZero()))

			// test Unmarshal
			assert.Equal(t, user.ID, 1)
			assert.Equal(t, user.Name, "John Doe")

			// test UnmarshalObjectKeys
			assert.Equal(t, userObj.ID, 1)
			assert.Equal(t, userObj.Name, "John Doe")
		})
	})
}
