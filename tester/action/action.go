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

package action

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
)

// actionFunc is a helper to create actions from functions
type actionFunc struct {
	run          func(t require.TestingT, ctx context.Context, value any, raw json.RawMessage, err error) error
	value        func(t require.TestingT) (any, bool)
	support      []Support
	supportsFunc func(supports Support) bool
	baseError    error
}

func (a *actionFunc) BaseError() error {
	return a.baseError
}

func (a *actionFunc) Supports(what Support) bool {
	if a.supportsFunc != nil {
		return a.supportsFunc(what)
	}
	for _, s := range a.support {
		if s == what {
			return true
		}
	}
	return false
}

func (a *actionFunc) Value(t require.TestingT) (any, bool) {
	if a.value != nil {
		return a.value(t)
	}
	return nil, false
}

func (a *actionFunc) Run(t require.TestingT, ctx context.Context, value any, raw json.RawMessage, err error) error {
	if a.run != nil {
		return a.run(t, ctx, value, raw, err)
	}
	return err
}
