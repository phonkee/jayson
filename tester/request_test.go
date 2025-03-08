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
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

// noopHandler does nothing for testing purposes
type noopHandler struct{}

func (n noopHandler) ServeHTTP(writer http.ResponseWriter, r *http.Request) {}

func TestRequest_Header(t *testing.T) {
	r := newRequest(http.MethodGet, "http://localhost:8080", &Deps{Handler: noopHandler{}})
	r.Header(t, "key", "value")
	r.Do(t, context.Background())
	assert.Equal(t, []string{ContentTypeJSON}, r.header.Values(ContentTypeHeader))
	assert.Equal(t, ContentTypeJSON, r.header.Get(ContentTypeHeader))
	assert.Equal(t, "value", r.header.Get("key"))
}

func TestRequest_Query(t *testing.T) {

}
