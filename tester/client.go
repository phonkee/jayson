package tester

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

// WithAPI runs the given function with a new APIClient instance
func WithAPI(t *testing.T, deps *Deps, fn func(APIClient)) {
	// validate deps first
	deps.Validate(t)

	// call closure with an api client
	fn(newClient(t, deps))
}

// newClient creates a new *client instance
func newClient(t *testing.T, deps *Deps) *client {
	return &client{
		deps: deps,
	}
}

// client is the implementation of rest APIClient testing helper
type client struct {
	deps *Deps
}

// Delete does a DELETE response to the APIClient
func (c *client) Delete(t *testing.T, path string) APIRequest {
	return c.Request(t, http.MethodDelete, path)
}

// Get does a GET response to the APIClient
func (c *client) Get(t *testing.T, path string) APIRequest {
	return c.Request(t, http.MethodGet, path)
}

// Post does a POST response to the APIClient
func (c *client) Post(t *testing.T, path string) APIRequest {
	return c.Request(t, http.MethodPost, path)
}

// Put does a PUT response to the APIClient
func (c *client) Put(t *testing.T, path string) APIRequest {
	return c.Request(t, http.MethodPut, path)
}

// Request does a response to the APIClient
func (c *client) Request(t *testing.T, method string, path string) APIRequest {
	return newRequest(method, path, c.deps)
}

// ReverseURL creates a path by given url name and url arguments
// name is the name of the route
// vars are the variables to replace in the route
// query is the query string to add to the URL
func (c *client) ReverseURL(t *testing.T, name string, vars ...string) string {
	require.NotNil(t, c.deps.Router, "Deps: Router is nil")

	// get route from router by name
	route := c.deps.Router.Get(name)
	require.NotNil(t, route, "route `%s` not found", name)

	// reverse the URL
	reversedURL, err := route.URL(vars...)
	require.NoErrorf(t, err, "failed to reverse URL for route `%s`", name)

	return reversedURL.String()
}
