# tester

Simple testing package for restful APIs.
Uses mux router for routing and direct http.Handler for handling requests.

# Usage

Tester provides a simple way to test restful APIs. It is based on the `testing` package and provides a simple way to test APIs.
You need to call `WithAPI` function with dependencies and then you provide closure where API will be available.


```go

var (
    router = mux.NewRouter()
)

func init() {
    // create a health check endpoint
    router.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        if err := json.NewEncoder(w).Encode(HealthResponse{
			Status: "ok",
        }); err != nil {
            panic(err)
        }
    }).Methods(http.MethodGet).Name("api:v1:health")
}

type HealthResponse struct {
    Status string `json:"status"`
}

func TestHealthHandler(t *testing.T) {// 
    tester.WithAPI(t, &Deps{
        Router: router,
        Handler: router,
    }, func(api *API) {
        var status string
        
        // unmarshal key from json object to value
        api.Get(t, api.ReverseURL(t, "api:v1:health")).
            Do(context.Background()).
            AssertStatus(t, http.StatusOK).
            Unmarshal(t, 
                APIObject(t, "status", &status),
            )
        assert.Equal(t, "ok", status)

        // direct unmarshal to struct
        response := HealthResponse{}
        api.Get(t, api.ReverseURL(t, "api:v1:health")).
            Do(context.Background()).
            AssertStatus(t, http.StatusOK).
            Unmarshal(t, &response)

        // assert json equals
        api.Get(t, api.ReverseURL(t, "api:v1:health")).
            Do(context.Background()).
            AssertStatus(t, http.StatusOK).
            AssertJsonEquals(t, HealthResponse{
                Status: "ok",	
            })
		
        // assert object key
        api.Get(t, api.ReverseURL(t, "api:v1:health")).
            Do(context.Background()).
            AssertStatus(t, http.StatusOK).
            AssertJsonKeyEquals(t, "status", "ok")
		
    })
}
```

# TODO:

List of future todos

# Author

Peter Vrba <phonkee@phonkee.eu>