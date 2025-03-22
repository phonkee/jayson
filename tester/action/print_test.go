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
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrint(t *testing.T) {
	for _, item := range []struct {
		name   string
		value  any
		raw    json.RawMessage
		expect []string
	}{
		{
			name: "test empty value",
			raw:  json.RawMessage(`{}`),
			expect: []string{
				"Raw: {}",
			},
		},
		{
			name:  "test value",
			value: "hello",
			raw:   json.RawMessage(`{}`),
			expect: []string{
				"Value[string]: hello\n  Raw: {}",
			},
		},
		{
			name:  "test empty raw",
			value: "hello",
			expect: []string{
				"Value[string]: hello",
			},
		},
		{
			name: "test empty raw and empty value",
			expect: []string{
				"<nil>",
			},
		},
	} {
		t.Run(item.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			ap := Print(w)
			assert.NoError(t, ap.Run(t, nil, item.value, item.raw, nil))
			got := w.String()
			for _, exp := range item.expect {
				assert.Contains(t, got, exp)
			}
		})
	}
}
