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
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	Error1 = errors.New("error1")
)

func TestRegisterError(t *testing.T) {
	assert.NoError(t, RegisterError(Error1, ExtStatus(http.StatusTeapot)))
	err := RegisterError(Error1, ExtStatus(http.StatusTeapot))
	assert.ErrorIs(t, err, ErrAlreadyRegistered)
	assert.ErrorIs(t, err, ErrAlreadyRegistered)
	assert.ErrorContains(t, err, "already registered")
}

type Some struct{}

func TestRegisterResponse(t *testing.T) {
	assert.NoError(t, RegisterResponse(Some{}, ExtStatus(http.StatusTeapot)))
	err := RegisterResponse(Some{}, ExtStatus(http.StatusTeapot))
	assert.ErrorIs(t, err, ErrAlreadyRegistered)
	assert.ErrorIs(t, err, ErrAlreadyRegistered)
	assert.ErrorContains(t, err, "already registered")
}

func TestDebug(t *testing.T) {
	Debug(nil)
	assert.Nil(t, _g.(*jayson).debug)
	Debug(zap.NewNop())
	assert.NotNil(t, _g.(*jayson).debug)
}

func TestResponse(t *testing.T) {
	rw := httptest.NewRecorder()
	Response(context.Background(), rw, 1)
	assert.Equal(t, DefaultSettings().DefaultResponseStatus, rw.Code)
}

func TestError(t *testing.T) {
	rw := httptest.NewRecorder()
	Error(context.Background(), rw, errors.New("error"))
	assert.Equal(t, DefaultSettings().DefaultErrorStatus, rw.Code)
}
