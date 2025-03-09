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
	"github.com/stretchr/testify/require"
	"net/url"
)

// Arguments adds arguments to the resolver, these are named arguments
func Arguments(t require.TestingT, args ...string) Extra {
	return resolverExtra(func() []string { return args }, nil)
}

// Query adds query parameters to the resolver
func Query(t require.TestingT, kv ...string) Extra {
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

// resolverExtraImpl is an implementation of Extra
// it is for internal use only
type resolverExtraImpl struct {
	argsFunc  func() []string
	queryFunc func() url.Values
}

// Args returns a list of additional arguments that can be used to resolve the URL
func (r *resolverExtraImpl) Args() []string {
	if r.argsFunc == nil {
		return nil
	}
	return r.argsFunc()
}

// Query returns a list of additional query parameters that can be used to resolve the URL
func (r *resolverExtraImpl) Query() url.Values {
	if r.queryFunc == nil {
		return nil
	}
	return r.queryFunc()
}

// resolverExtra is an implementation of Extra
func resolverExtra(argsFunc func() []string, queryFunc func() url.Values) Extra {
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
