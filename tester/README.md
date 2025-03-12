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

# TODO


Tester library provides a set of assertions that can be used to test APIs.
Some are basic, some are more complex. Let's go through them.
Let's assume that we have instance of APIClient `api` that is used to make requests.

## Status

Asserts that response status code is equal to provided status code.

```go
api.Get(t, "/api/v1/health").
    Do(t, context.TODO()).
    Status(t, http.StatusOK)
```

## AssertHeaderValue

Asserts that response header into is equal to provided into.

```go
api.Get(t, "/api/v1/health").
    Do(t, context.TODO()).
    AssertHeaderValue(t, "Content-Type", "application/json")
```

## AssertJsonEquals

Asserts that response body is equal to provided object/string/bytes.

```go
type HealthResponse struct {
    Status string `json:"status"`
}

api.Get(t, "/api/v1/health").
    Do(t, context.TODO()).
    AssertJsonEquals(t, HealthResponse{
        Status: "ok",
    })
```

## Unmarshal

Unmarshal is not assertion but it is used to unmarshal response body to provided object.

```go
var response HealthResponse
api.Get(t, "/api/v1/health").
    Do(t, context.TODO()).
    Status(t, http.StatusOK).
    Unmarshal(t, &response)
```

## AssertJsonPath

This assertion is the most complex one. It is used to assert json path in response body.
It not just asserts that path exists but also that into is equal to provided into.
On top of that there are ways to assert that into is not only equal but also greater, less, etc.
Path can contain also array indexes.
Let's see some examples.
Let's suppose the api returns following json object for `/api/v1/users` endpoint.

```json
{
    "status": "ok",
    "data": {
        "users": [
            {
                "id": 1,
                "name": "Peter"
            },
            {
                "id": 2,
                "name": "John"
            }
        ]
    }
}
```

Now let's see some examples how we can assert data by json path.

```go
// assert that status is ok
api.Get(t, "/api/v1/users").
    Do(t, context.TODO()).
    AssertJsonPath(t, "status", "ok").

// assert that users array has length of 2
api.Get(t, "/api/v1/users").
    Do(t, context.TODO()).
    AssertJsonPath(t, "data.users.__len__", 2).

// assert that first user has id 1
api.Get(t, "/api/v1/users").
    Do(t, context.TODO()).
    AssertJsonPath(t, "data.users.0.id", 1).

// assert that data has key users
api.Get(t, "/api/v1/users").
    Do(t, context.TODO()).
    AssertJsonPath(t, "data.users.__keys__", []string{"users"}).

// prepare simple struct for partial unmarshalling
type SimpleUser struct {
    ID int `json:"id"`
}
// assert that users data equals to provided slice
api.Get(t, "/api/v1/users").
    Do(t, context.TODO()).
    AssertJsonPath(t, "data.users", []SimpleUser{
        {ID: 1}, {ID: 2}
    })

// assert key users exists
api.Get(t, "/api/v1/users").
    Do(t, context.TODO()).
    AssertJsonPath(t, "data.users.__exists__", nil)

// assert that first user id is greater than 0
api.Get(t, "/api/v1/users").
    Do(t, context.TODO()).
    AssertJsonPath(t, "data.users.0.id.__gte__", 1)

// assert that name of first user is Peter
api.Get(t, "/api/v1/users").
    Do(t, context.TODO()).
    AssertJsonPath(t, "data.users.0.name", "Peter")

// assert that name of first user is not John
api.Get(t, "/api/v1/users").
    Do(t, context.TODO()).
    AssertJsonPath(t, "data.users.0.name.__neq__", "John")
```

# TODO
- [ ] AssertAny - any of given actions succeeds
- [ ] AssertAll - all of given actions succeeds
- [ ] AssertRegex - assert that response body (path) matches regex

# Author

Peter Vrba <phonkee@phonkee.eu>