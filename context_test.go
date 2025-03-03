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
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCtxErrorValue(t *testing.T) {
	ctx := context.Background()
	err, ok := ContextErrorValue(ctx)
	assert.False(t, ok)
	assert.Nil(t, err)

	ctx = contextWithErrorValue(ctx, assert.AnError)
	err, ok = ContextErrorValue(ctx)
	assert.True(t, ok)
	assert.Equal(t, assert.AnError, err)
}

func TestCtxSettingsValue(t *testing.T) {
	ctx := context.Background()
	settings := ContextSettingsValue(ctx)
	assert.Zero(t, settings)

	def := DefaultSettings()
	ctx = contextWithSettingsValue(ctx, def)
	settings = ContextSettingsValue(ctx)
	assert.Equal(t, def, settings)
}

func TestCtxObjectValue(t *testing.T) {
	ctx := context.Background()
	obj, ok := ContextObjectValue[int](ctx)
	assert.False(t, ok)
	assert.Zero(t, obj)

	ctx = contextWithObjectValue(ctx, 42)
	obj, ok = ContextObjectValue[int](ctx)
	assert.True(t, ok)
	assert.Equal(t, 42, obj)
}
