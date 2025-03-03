package tester

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type HealthResponse struct {
	Status string `json:"status"`
	Host   string `json:"host"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(HealthResponse{
		Status: "something",
		Host:   "localhost",
	}); err != nil {
		panic(err)
	}
}

func newHealthRouter(t *testing.T) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/health", HealthHandler).Methods(http.MethodGet).Name("api:v1:health")
	return router
}

func TestAPI(t *testing.T) {
	router := newHealthRouter(t)
	WithAPI(t, &Deps{
		Router: router,
	}, func(api APIClient) {
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
		rr := HealthResponse{}

		// do response
		api.Request(t, http.MethodGet, api.ReverseURL(t, "api:v1:health")).
			Do(t, ctx).
			AssertStatus(t, http.StatusOK).
			AssertJsonEquals(t, `{"status": "something", "host": "localhost"}`).
			Unmarshal(t,
				APIObject(t,
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
}

// ptrTo helper
func ptrTo[T any](v T) *T {
	return &v
}
