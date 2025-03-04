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

package jayson

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegistryNew(t *testing.T) {
	r := newRegistry[error]()
	assert.NotNil(t, r)
}

func TestRegistryExists(t *testing.T) {
	r := newRegistry[error]()
	assert.NotNil(t, r)

	// check if error Exists
	assert.False(t, r.Exists(assert.AnError))
	assert.NoError(t, r.Register(assert.AnError, nil))
	assert.True(t, r.Exists(assert.AnError))
}

func TestRegistryGet(t *testing.T) {
	r := newRegistry[error]()
	assert.NotNil(t, r)

	// Get should return false
	_, ok := r.Get(assert.AnError)
	assert.False(t, ok)

	// now register type
	ext := []Extension{ExtFunc(nil, nil)}
	assert.NoError(t, r.Register(assert.AnError, ext))

	// Get should return true
	e, ok := r.Get(assert.AnError)
	assert.True(t, ok)
	assert.Equal(t, ext, e)
}
