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

import "context"

// contextKey is a type for storing values in context.
type contextKey int

const (
	// contextSettingsKey is the key used to store the settings value in the context.
	contextSettingsKey contextKey = 1 << iota

	// contextErrorKey is the key used to store the error value in the context.
	contextErrorKey

	// contextObjectKey is the key used to store the object value in the context.
	contextObjectKey
)

// ContextErrorValue returns the error value stored in the context.
func ContextErrorValue(ctx context.Context) (error, bool) {
	err, ok := ctx.Value(contextErrorKey).(error)
	return err, ok
}

// ContextObjectValue returns the object value stored in the context.
func ContextObjectValue[T any](ctx context.Context) (T, bool) {
	val, ok := ctx.Value(contextObjectKey).(T)
	return val, ok
}

// ContextSettingsValue returns the settings value stored in the context.
func ContextSettingsValue(ctx context.Context) Settings {
	var result Settings
	v, ok := ctx.Value(contextSettingsKey).(Settings)
	if ok {
		result = v
	}
	return result
}

// contextWithErrorValue adds the error value to the context.
func contextWithErrorValue(ctx context.Context, err error) context.Context {
	return context.WithValue(ctx, contextErrorKey, err)
}

// contextWithObjectValue adds the object value to the context.
func contextWithObjectValue(ctx context.Context, obj any) context.Context {
	return context.WithValue(ctx, contextObjectKey, obj)
}

// contextWithSettingsValue adds the settings value to the context.
func contextWithSettingsValue(ctx context.Context, settings Settings) context.Context {
	return context.WithValue(ctx, contextSettingsKey, settings)
}
