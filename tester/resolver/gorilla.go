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

package resolver

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

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
func (g *gorillaResolver) ReverseURL(t require.TestingT, name string, extra ...Extra) string {
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

	// add query parameters to the URL
	q := reversedURL.Query()
	for _, e := range extra {
		qv := e.Query()
		if qv == nil {
			continue
		}
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
