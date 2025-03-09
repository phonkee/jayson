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
	"encoding/json"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/constraints"
)

// Assertion is a type that represents a single assertion
type Assertion interface {
	Assert(t require.TestingT, value json.RawMessage, err error) (any, error)
	UseEmpty() bool
}

// AssertEqual asserts that given value is equal to the value in the response
func AssertEqual(value any) Assertion {
	return newAssertion(
		func(t require.TestingT, what json.RawMessage, err error) (any, error) {
			return nil, nil
		},
		nil,
	)
}

// newAssertion creates a new assertion
func newAssertion(
	assert func(t require.TestingT, value json.RawMessage, err error) (any, error),
	useEmpty func() bool,
) Assertion {
	return &assertion{
		assert:   assert,
		useEmpty: useEmpty,
	}
}

type assertion struct {
	assert   func(t require.TestingT, value json.RawMessage, err error) (any, error)
	useEmpty func() bool
}

func (a *assertion) Assert(t require.TestingT, value json.RawMessage, err error) (any, error) {
	if a.assert != nil {
		return a.assert(t, value, err)
	}
	return nil, nil
}

func (a *assertion) UseEmpty() bool {
	if a.useEmpty != nil {
		return a.useEmpty()
	}
	return false
}

type Number interface {
	constraints.Integer | constraints.Float
}

func AssertLen[T Number](length T) Assertion {
	return nil
}

//	func AssertExists(bool) Assertion {
//		return nil
//	}
//
//	func AssertGt(value int) Assertion {
//		return nil
//	}
//
//	func AssertGte(value int) Assertion {
//		return nil
//	}
//
//func AssertLt(value int) Assertion {
//	return nil
//}
//
//func AssertLte(value int) Assertion {
//	return nil
//}
//
//func AssertNotEqual(value any) Assertion {
//	return nil
//}
