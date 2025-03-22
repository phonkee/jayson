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
	"errors"
	"fmt"
)

var (
	// ErrUnmarshal	when unmarshalling the response
	ErrUnmarshal = errors.New("unmarshal error")

	// ErrActionNotApplied when action is not applied
	ErrActionNotApplied = errors.New("action not applied")

	// ErrAction when action is not applied
	ErrAction = errors.New("action")
)

// newErrAction creates a new error with the given name and error
func newErrAction(name string) error {
	return fmt.Errorf("%w: `%s`", ErrAction, name)
}

var (
	ErrActionAssertEquals        = newErrAction("AssertEquals")
	ErrActionAssertExists        = newErrAction("AssertExists")
	ErrActionAssertIn            = newErrAction("AssertIn")
	ErrActionAssertKeys          = newErrAction("AssertKeys")
	ErrActionAssertLen           = newErrAction("AssertLen")
	ErrActionAssertNotEquals     = newErrAction("AssertNotEquals")
	ErrActionAssertNotExists     = newErrAction("AssertNotExists")
	ErrActionAssertNotIn         = newErrAction("AssertNotIn")
	ErrActionAssertGt            = newErrAction("AssertGt")
	ErrActionAssertGte           = newErrAction("AssertGte")
	ErrActionAssertLt            = newErrAction("AssertLt")
	ErrActionAssertLte           = newErrAction("AssertLte")
	ErrActionAssertAll           = newErrAction("AssertAll")
	ErrActionAssertAny           = newErrAction("AssertAny")
	ErrActionAssertRegexMatch    = newErrAction("ErrActionAssertRegexMatch")
	ErrActionAssertRegexSearch   = newErrAction("ErrActionAssertRegexSearch")
	ErrActionUnmarshal           = newErrAction("Unmarshal")
	ErrActionUnmarshalObjectKeys = newErrAction("UnmarshalObjectKeys")
	ErrActionAssertNot           = newErrAction("AssertNot")
	ErrNotPresent                = errors.New("not present")
)
