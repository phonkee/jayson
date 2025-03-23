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
	"github.com/phonkee/jayson/internal"
	"github.com/stretchr/testify/require"
)

// JsonPath maps json path to given action it is relative to caller
func JsonPath(path string, action Action) Action {
	return &actionFunc{
		run: func(t require.TestingT, ctx context.Context, _ any, raw json.RawMessage, err error) error {
			if action == nil {
				return fmt.Errorf("%w: no action defined", ErrActionPath)
			}

			// get unmarshal action value function
			avf := ctx.Value(ContextKeyUnmarshalActionValue).(unmarshalActionValueFunc)

			// now we parse raw message
			raw, err = internal.ExtractJsonPath(t, path, raw)

			var value any

			// if we don't have any error, we can try to unmarshal the action
			if err == nil {
				value, err = avf(t, raw, action)
			}

			// check if unmarshal was not applied
			if errors.Is(err, ErrActionNotApplied) {
				err = nil
			}

			// now run action
			return action.Run(t, ctx, value, raw, err)
		},
		support:   []Support{SupportJson},
		baseError: ErrActionPath,
	}
}
