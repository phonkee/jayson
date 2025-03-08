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
	"github.com/phonkee/jayson/tester/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"strings"
	"testing"
)

func TestDeps_Validate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		for _, item := range []struct {
			deps *Deps
		}{
			{&Deps{Router: nil, Handler: http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}), Address: ""}},
			{&Deps{Router: nil, Handler: nil, Address: "localhost:8080"}},
		} {
			m := mocks.NewMockRequireTestingT(t)
			item.deps.Validate(require.New(m))
		}
	})

	t.Run("invalid", func(t *testing.T) {
		for _, item := range []struct {
			deps *Deps
		}{
			{&Deps{Router: nil, Handler: nil, Address: ""}},
		} {
			m := mocks.NewMockRequireTestingT(t)
			m.On("Errorf", mock.Anything, mock.MatchedBy(func(msg string) bool {
				return strings.Contains(msg, "exampleHandler or Address is required")
			})).Once()
			m.On("FailNow").Once()
			item.deps.Validate(require.New(m))
			m.AssertExpectations(t)
			m.AssertNumberOfCalls(t, "Errorf", 1)
			m.AssertNumberOfCalls(t, "FailNow", 1)
		}
	})
}
