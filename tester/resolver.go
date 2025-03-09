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

package tester

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"net/url"
)

// Resolver is an interface that can be used to resolve URLs
// from the given name and extra arguments
//
//go:generate mockery --name Resolver --filename resolver.go
type Resolver interface {
	// ReverseURL returns a URL by given name and extra arguments
	ReverseURL(t require.TestingT, name string, extra ...ResolverExtra) string
}

// ResolverExtra is an interface that can be used to pass additional information to the url resolver
//
//go:generate mockery --name ResolverExtra --filename resolver_extra.go
type ResolverExtra interface {
	// Args returns a list of additional arguments that can be used to resolve the URL
	Args() []string
	// Query returns a list of additional query parameters that can be used to resolve the URL
	Query() url.Values
}

// NewGorillaResolver creates a new gorilla resolver
func NewGorillaResolver(t require.TestingT, router *mux.Router) Resolver {
	return &gorillaResolver{router: router}
}

// gorillaResolver is a resolver that uses gorilla mux router to resolve URLs
type gorillaResolver struct {
	router *mux.Router
}

// ReverseURL returns a URL by given name and extra arguments
func (g *gorillaResolver) ReverseURL(t require.TestingT, name string, extra ...ResolverExtra) string {
	require.NotNil(t, g.router, "Deps: Router is nil")

	// get route from router by name
	route := g.router.Get(name)
	require.NotNil(t, route, "route `%s` not found", name)

	// reverse the URL
	reversedURL, err := route.URL()
	require.NoErrorf(t, err, "failed to reverse URL for route `%s`", name)

	return reversedURL.String()
}
