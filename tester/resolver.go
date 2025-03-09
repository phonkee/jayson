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

// ResolverArgs adds arguments to the resolver
func ResolverArgs(t require.TestingT, args ...string) ResolverExtra {
	return resolverExtra(func() []string { return args }, nil)
}

// ResolverQuery adds query parameters to the resolver
func ResolverQuery(t require.TestingT, kv ...string) ResolverExtra {
	return resolverExtra(nil, func() url.Values {
		urlValues := url.Values{}
		for i := 0; i < len(kv); i += 2 {
			if i+1 < len(kv) {
				urlValues.Add(kv[i], kv[i+1])
			}
		}
		return urlValues
	})
}

// resolverExtraImpl is an implementation of ResolverExtra
type resolverExtraImpl struct {
	argsFunc  func() []string
	queryFunc func() url.Values
}

// Args returns a list of additional arguments that can be used to resolve the URL
func (r *resolverExtraImpl) Args() []string {
	return r.argsFunc()
}

// Query returns a list of additional query parameters that can be used to resolve the URL
func (r *resolverExtraImpl) Query() url.Values {
	return r.queryFunc()
}

// resolverExtra is an implementation of ResolverExtra
func resolverExtra(argsFunc func() []string, queryFunc func() url.Values) ResolverExtra {
	if argsFunc == nil {
		argsFunc = func() []string { return nil }
	}
	if queryFunc == nil {
		queryFunc = func() url.Values { return nil }
	}
	return &resolverExtraImpl{
		argsFunc:  argsFunc,
		queryFunc: queryFunc,
	}
}

// NewGorillaResolver creates a new gorilla resolver
func NewGorillaResolver(t require.TestingT, router *mux.Router) Resolver {
	require.NotNilf(t, router, "resolver: Router is nil")
	return &gorillaResolver{router: router}
}

// gorillaResolver is a resolver that uses gorilla mux router to resolve URLs
type gorillaResolver struct {
	router *mux.Router
}

// ReverseURL returns a URL by given name and extra arguments
func (g *gorillaResolver) ReverseURL(t require.TestingT, name string, extra ...ResolverExtra) string {
	require.NotNil(t, g.router, "resolver: Router is nil")

	// get route from router by name
	route := g.router.Get(name)
	require.NotNil(t, route, "route `%s` not found", name)

	// add arguments to the URL
	args := make([]string, 0)
	for _, e := range extra {
		args = append(args, e.Args()...)
	}

	// reverse the URL
	reversedURL, err := route.URL(args...)
	require.NoErrorf(t, err, "failed to reverse URL for route `%s`", name)

	// add query
	q := reversedURL.Query()
	for _, e := range extra {
		for k, v := range e.Query() {
			for _, vv := range v {
				q.Add(k, vv)
			}
		}
	}
	if len(q) > 0 {
		reversedURL.RawQuery = q.Encode()
	}

	return reversedURL.String()
}
