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
	"golang.org/x/exp/constraints"
)

type Support int

const (
	SupportHeader Support = iota
	SupportJson
	SupportStatus
)

// Number constraint for number types
type Number interface {
	constraints.Integer | constraints.Float
}

// Action is an action that can be performed on the response. It has currently 2 meanings: Assertion and Unmarshal.
// TODO: add path or id to run so we can print better error messages?
type Action interface {
	Run(t require.TestingT, ctx context.Context, value any, raw json.RawMessage, err error) error
	Value(t require.TestingT) (any, bool)
	Supports(support Support) bool
	BaseError() error
}

// actionNegate is a special action that has custom handling for negation.
type actionNegate interface {
	Run(t require.TestingT, ctx context.Context, value any, raw json.RawMessage, err error) error
}

type contextKey int

const (
	// ContextKeyUnmarshalActionValue is a context key for the unmarshal action value.
	ContextKeyUnmarshalActionValue contextKey = iota
)

// unmarshalActionValueFunc is a function that unmarshal the action value from the response.
type unmarshalActionValueFunc = func(t require.TestingT, raw json.RawMessage, a Action) (any, error)
