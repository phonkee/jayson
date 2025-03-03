package tester

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

// Deps is the dependencies for the APIClient
type Deps struct {
	// Router - currently required
	Router *mux.Router
	// If the Handler is nil, Router will be used as the handler
	Handler http.Handler
}

// Validate deps
func (d *Deps) Validate(t *testing.T) {
	require.NotNil(t, d.Router, "Deps: Router is nil")
	// if the handler is nil, use the router
	if d.Handler == nil {
		d.Handler = d.Router
	}
}
