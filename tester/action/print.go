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
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"os"
)

// Print is an action that prints the value and error to the given writer.
// By default, it prints to os.Stdout.
func Print(target ...io.Writer) Action {
	var w io.Writer
	if len(target) > 0 && target[0] != nil {
		w = target[0]
	} else {
		w = os.Stdout
	}
	// printf is a helper function to print formatted output and ignore errors
	printf := func(format string, args ...any) {
		_, _ = fmt.Fprintf(w, format+"\n", args...)
	}
	return &actionFunc{
		run: func(t require.TestingT, ctx context.Context, value any, raw json.RawMessage, err error) error {
			printf("Print:")
			if err != nil {
				printf("  Error: %v", err)
				return err
			}
			if value != nil {
				printf("  Value[%T]: %v", value, value)
			}
			if raw != nil {
				printf("  Raw: %v", string(raw))
			}
			if value == nil && raw == nil {
				printf("  <nil>")
			}

			return nil
		},
		support: []Support{SupportHeader, SupportJson, SupportStatus},
	}
}
