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
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
)

// AssertAll asserts that all actions are executed
func AssertAll(a ...Action) Action {
	return &actionFunc{
		run: func(t require.TestingT, ctx context.Context, _ any, raw json.RawMessage, err error) error {
			if err != nil {
				return err
			}

			// run all actions until one fails
			for _, action := range a {
				fn := ctx.Value(ContextKeyUnmarshalActionValue).(unmarshalActionValueFunc)

				// get value from action
				value, err := fn(t, raw, action)

				// check if unmarshal was not applied
				if err != nil {
					if errors.Is(err, ErrActionNotApplied) {
						err = nil
					} else {
						return fmt.Errorf("%w: %w", action.BaseError(), err)
					}
				}

				// run action
				if err := action.Run(t, ctx, value, raw, err); err != nil {
					return fmt.Errorf("%w: %w", action.BaseError(), err)
				}
			}

			return nil
		},
		support:   []Support{SupportHeader, SupportJson, SupportStatus},
		baseError: ErrActionAssertAll,
	}
}

// AssertAny asserts that any action succeeded
func AssertAny(a ...Action) Action {
	return &actionFunc{
		run: func(t require.TestingT, ctx context.Context, _ any, raw json.RawMessage, err error) error {
			if err != nil {
				return err
			}

			// iterate over all actions and run them until one succeeds
			for _, action := range a {
				fn := ctx.Value(ContextKeyUnmarshalActionValue).(unmarshalActionValueFunc)

				// get value from action
				value, err := fn(t, raw, action)

				// check if unmarshal was not applied
				if err != nil {
					if errors.Is(err, ErrActionNotApplied) {
						err = nil
					} else {
						return fmt.Errorf("%w: %w", action.BaseError(), err)
					}
				}

				// run action
				if err := action.Run(t, ctx, value, raw, err); err == nil {
					return nil
				}
			}

			return fmt.Errorf("no action succeeded")
		},
		support:   []Support{SupportHeader, SupportJson, SupportStatus},
		baseError: ErrActionAssertAny,
	}
}
