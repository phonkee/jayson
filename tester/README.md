# tester

Simple testing package for restful APIs.

# Usage

Tester provides a simple way to test restful APIs. It is based on the `testing` package and provides a simple way to
test APIs.
It provides two main functions:

- `WithHttpServer` - starts a http server and runs given closure with the server
- `WithAPI` - runs given closure with the API client

# Example

Let's see an example of how to use the tester package.

```go
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
				Json(t, "users.0", action.AssertKeys("id", "name")).
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
				}))

			// test Unmarshal
			assert.Equal(t, user.ID, 1)
			assert.Equal(t, user.Name, "John Doe")

			// test UnmarshalObjectKeys
			assert.Equal(t, userObj.ID, 1)
			assert.Equal(t, userObj.Name, "John Doe")
		})
	})
}

```

# TODO

- [ ] AssertGt, AssertGte, AssertLt, AssertLte should accept also float64 and float32

# Author

Peter Vrba <phonkee@phonkee.eu>