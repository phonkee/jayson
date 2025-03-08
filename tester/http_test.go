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
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestWithHttpServer(t *testing.T) {
	t.Run("test invalid args", func(t *testing.T) {
		for _, item := range []struct {
			tt      *testing.T
			handler http.Handler
			fn      func(t *testing.T, address string)
		}{
			{nil, nil, func(t *testing.T, address string) {}},
			{t, nil, func(t *testing.T, address string) {}},
			{t, http.HandlerFunc(func(writer http.ResponseWriter, r *http.Request) {}), nil},
		} {
			assert.Panics(
				t,
				func() {
					WithHttpServer(
						item.tt,
						item.handler,
						item.fn,
					)
				},
			)
		}
	})

	t.Run("test valid args", func(t *testing.T) {
		value := 0
		WithHttpServer(
			t,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				value++
				w.WriteHeader(http.StatusOK)
			}),
			func(t *testing.T, address string) {
				resp, err := http.DefaultClient.Get(fmt.Sprintf("http://%v", address))
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				assert.Equal(t, 1, value)
			},
		)
	})
}
