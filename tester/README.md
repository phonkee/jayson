# tester

Simple testing package for restful APIs.

# Usage

Tester provides a simple way to test restful APIs. It is based on the `testing` package and provides a simple way to test APIs.
You need to call `WithAPI` function with dependencies and then you provide closure where API will be available.
This library supports http.Handler testing as well as http server testing (Address).

```go
package example_test

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/phonkee/jayson/tester"
	"github.com/phonkee/jayson/tester/action"
	"github.com/phonkee/jayson/tester/resolver"
	"net/http"
	"testing"
)

var (
	// we will use gorilla mux router
	router = mux.NewRouter()
)

func init() {
	// create a health check endpoint
	router.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(HealthResponse{
			StatusDatabase: "ok",
		}); err != nil {
			panic(err)
		}
	}).Methods(http.MethodGet).Name("api:v1:health")
}

// HealthResponse is a simple response struct that returns status
type HealthResponse struct {
	StatusDatabase string `json:"status_db"`
}

func TestHealthHandler(t *testing.T) {
	tester.WithAPI(t, &Deps{
		Resolver: resolver.NewGorillaMuxResolver(t, router), // url resolver for gorilla mux
		Handler:  router,                                    // use router as http.Handler
	}, func(api *API) {
		var status string

		// unmarshal key from json object to value
		api.Get(t, api.ReverseURL(t, "api:v1:health")).
			Do(t, context.TODO()).
			Status(t, action.Unmarshal(&status)).
			//Unmarshal(t,
			//	APIObject(t, "status", &status), // APIObject deconstructs json object to value given key value pairs
			//)
			assert.Equal(t, "ok", status)

		// direct unmarshal to struct
		response := HealthResponse{}
		api.Get(t, api.ReverseURL(t, "api:v1:health")).
			Do(t, context.TODO()).
			Status(t, action.AssertEqual(http.StatusOK)).
			Json(t, "", action.Unmarshal(&response))

		// assert json equals
		api.Get(t, api.ReverseURL(t, "api:v1:health")).
			Do(t, context.TODO()).
			Status(t, action.AssertEqual(http.StatusOK)).
			Json(t, "", action.AssertEqual(HealthResponse{
				StatusDatabase: "ok",
			}))

		// assert object key
		api.Get(t, api.ReverseURL(t, "api:v1:health")).
			Do(t, context.TODO()).
			Status(t, action.AssertEqual(http.StatusOK)).
			Json(t, "status_db", action.AssertEqual("ok"))
	})
}

```

# Response

Response provides multiple methods to inspect response.
Each accepts action which can be one from Assert actions or Unmarshal actions.
Let me show example of all of them

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
				Json(t, "users.0.name", action.AssertNotEquals("Johnson Doe")).
				Json(t, "users.0.name", action.AssertEquals("John Doe")).
				Json(t, "users.0", action.AssertEquals(json.RawMessage(`{"id":1,"name":"John Doe"}`))).
				Json(t, "users.0.name", action.AssertIn("John Doe", "Peter Vrba")).
				Json(t, "users.0.name", action.AssertNotIn("Johnson Doe", "Peter Vrba")).
				Json(t, "users", action.AssertExists()).
				Json(t, "user", action.AssertNotExists()).
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
					action.AssertNotExists(),
				)).
				Json(t, "users.0.id", action.AssertRegex(
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
- [ ] AssertZero - assert zero value (0, "", nil, [], {}, false)

# Author

Peter Vrba <phonkee@phonkee.eu>