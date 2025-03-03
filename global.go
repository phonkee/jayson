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
	"go.uber.org/zap"
	"net/http"
)

var (
	// _g is a global instance of Jayson
	_g = New(DefaultSettings())
)

// Debug enables debug mode via zap logger in _g instance of Jayson
func Debug(logger *zap.Logger) {
	_g.Debug(logger)
}

// Error writes error to _g instance of Jayson
func Error(ctx context.Context, w http.ResponseWriter, err error, ext ...Extension) {
	_g.Error(ctx, w, err, ext...)
}

// RegisterError registers extFunc for given error in _g instance of Jayson
func RegisterError(err error, ext ...Extension) error {
	return _g.RegisterError(err, ext...)
}

// RegisterResponse registers extFunc for given object in _g instance of Jayson
func RegisterResponse(what any, ext ...Extension) error {
	return _g.RegisterResponse(what, ext...)
}

// Response writes response to http.ResponseWriter using _g instance of Jayson
func Response(ctx context.Context, w http.ResponseWriter, what any, ext ...Extension) {
	_g.Response(ctx, w, what, ext...)
}
