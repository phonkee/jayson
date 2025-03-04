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
	"github.com/phonkee/jayson/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

// extErrorDetail is an extFunc that adds an error detail to the response object.
func extErrorDetail(detail string) Extension {
	return ExtObjectKeyValue("errorDetail", detail)
}

func assertExtensionResponse(t *testing.T, expect string, obj any, ext ...Extension) {
	jay := New(DefaultSettings())
	rw := httptest.NewRecorder()
	jay.Response(context.Background(), rw, obj, ext...)
	body := rw.Body.String()
	assert.JSONEqf(t, expect, body, "expected expectJSON: `%s`, got: `%s`", expect, body)
}

func TestExtension_Response(t *testing.T) {

	t.Run("test ExtObjectKeyValue", func(t *testing.T) {
		assertExtensionResponse(t, `{"key":[1,2,3]}`, ExtObjectKeyValue("key", []int{1, 2, 3}))
		assertExtensionResponse(t, `[1,2,3]`, []int{1, 2, 3})
		assertExtensionResponse(t, `{"key":[1,2,3]}`, ExtObjectKeyValue("key", []int{1, 2, 3}), ExtNoop())
		assertExtensionResponse(t, `{"errorDetail":"detail","key":[1,2,3]}`, ExtObjectKeyValue("key", []int{1, 2, 3}), extErrorDetail("detail"))
	})

	t.Run("test extErrorDetail", func(t *testing.T) {
		assertExtensionResponse(t, `{"errorDetail":"detail"}`, extErrorDetail("detail"))
		assertExtensionResponse(t, `{"errorDetail":"detail"}`, extErrorDetail("detail"), ExtNoop())
		assertExtensionResponse(t, `{"errorDetail":"detail","key":[1,2,3]}`, ExtObjectKeyValue("key", []int{1, 2, 3}), extErrorDetail("detail"))
	})

	t.Run("test ExtObjectKeyValuef", func(t *testing.T) {
		assertExtensionResponse(t, `{"key":"hello: value"}`, ExtObjectKeyValuef("key", "hello: %s", "value"))
	})
}

func TestExtFirst(t *testing.T) {
	t.Run("test ExtResponseWriter", func(t *testing.T) {
		t.Run("test empty", func(t *testing.T) {
			ext := ExtFirst()
			rw := httptest.NewRecorder()
			assert.False(t, ext.ExtendResponseWriter(context.Background(), rw))
		})
		t.Run("test existing", func(t *testing.T) {
			ext := ExtFirst(
				ExtStatus(http.StatusInternalServerError),
				ExtStatus(http.StatusTeapot),
			)

			rw := httptest.NewRecorder()
			assert.True(t, ext.ExtendResponseWriter(context.Background(), rw))
			assert.Equal(t, http.StatusInternalServerError, rw.Code)
		})
	})

	t.Run("test ExtResponseObject", func(t *testing.T) {
		t.Run("test empty", func(t *testing.T) {
			ext := ExtFirst()

			obj := make(map[string]any)

			assert.False(t, ext.ExtendResponseObject(context.Background(), obj))
		})
		t.Run("test existing", func(t *testing.T) {
			ext := ExtFirst(
				ExtObjectKeyValue("key1", "value1"),
				ExtObjectKeyValue("key1", "value1"),
			)

			obj := make(map[string]any)

			assert.True(t, ext.ExtendResponseObject(context.Background(), obj))
			assert.Equal(t, map[string]any{"key1": "value1"}, obj)
		})
	})

}

func TestExtChain(t *testing.T) {
	t.Run("test empty", func(t *testing.T) {
		ext := ExtChain()
		rw := httptest.NewRecorder()
		assert.False(t, ext.ExtendResponseWriter(context.Background(), rw))
	})
	t.Run("test existing", func(t *testing.T) {
		ext := ExtChain(
			ExtObjectKeyValue("key1", "value1"),
			ExtObjectKeyValue("key2", "value2"),
			ExtHeaderValue("key1", "value1"),
			ExtHeaderValue("key2", "value2"),
		)

		rw := httptest.NewRecorder()
		obj := make(map[string]any)

		ext.ExtendResponseWriter(context.Background(), rw)
		ext.ExtendResponseObject(context.Background(), obj)

		assert.Equal(t, "value1", rw.Header().Get("key1"))
		assert.Equal(t, "value2", rw.Header().Get("key2"))
		assert.Equal(t, map[string]any{"key1": "value1", "key2": "value2"}, obj)

	})

}

func TestExtHeader(t *testing.T) {
	t.Run("test empty", func(t *testing.T) {
		ext := ExtHeader(nil)
		rw := httptest.NewRecorder()
		assert.False(t, ext.ExtendResponseWriter(context.Background(), rw))
	})
	t.Run("test existing", func(t *testing.T) {
		ext := ExtHeader(http.Header{"key1": []string{"value1"}})

		rw := httptest.NewRecorder()
		obj := make(map[string]any)

		ext.ExtendResponseWriter(context.Background(), rw)
		ext.ExtendResponseObject(context.Background(), obj)

		assert.Equal(t, "value1", rw.Header().Get("key1"))
		assert.Empty(t, obj)
	})
}

func TestExtConditional(t *testing.T) {
	t.Run("test empty", func(t *testing.T) {
		m1 := mocks.NewExtension(t)
		m2 := mocks.NewExtension(t)
		ext := ExtConditional(
			ExtFunc(nil, nil),
			m1,
			m2,
		)
		rw := httptest.NewRecorder()
		obj := make(map[string]any)
		assert.False(t, ext.ExtendResponseWriter(context.Background(), rw))
		assert.False(t, ext.ExtendResponseObject(context.Background(), obj))
		m1.AssertExpectations(t)
		m2.AssertExpectations(t)
	})
	t.Run("test existing", func(t *testing.T) {
		t.Run("test ExtendResponseWriter", func(t *testing.T) {
			m1 := mocks.NewExtension(t)
			m1.On("ExtendResponseWriter", mock.Anything, mock.Anything).Return(true)
			m1.On("ExtendResponseObject", mock.Anything, mock.Anything).Return(false)
			m2 := mocks.NewExtension(t)
			m2.On("ExtendResponseWriter", mock.Anything, mock.Anything).Return(true)
			ext := ExtConditional(
				m1,
				m2,
			)
			rw := httptest.NewRecorder()
			obj := make(map[string]any)
			assert.True(t, ext.ExtendResponseWriter(context.Background(), rw))
			assert.False(t, ext.ExtendResponseObject(context.Background(), obj))
			m1.AssertExpectations(t)
			m1.AssertNumberOfCalls(t, "ExtendResponseWriter", 1)
			m2.AssertExpectations(t)
			m2.AssertNumberOfCalls(t, "ExtendResponseWriter", 1)
		})
		t.Run("test ExtendResponseObject", func(t *testing.T) {
			m1 := mocks.NewExtension(t)
			m1.On("ExtendResponseWriter", mock.Anything, mock.Anything).Return(false)
			m1.On("ExtendResponseObject", mock.Anything, mock.Anything).Return(true)
			m2 := mocks.NewExtension(t)
			m2.On("ExtendResponseObject", mock.Anything, mock.Anything).Return(true)
			ext := ExtConditional(
				m1,
				m2,
			)
			rw := httptest.NewRecorder()
			obj := make(map[string]any)
			assert.False(t, ext.ExtendResponseWriter(context.Background(), rw))
			assert.True(t, ext.ExtendResponseObject(context.Background(), obj))
			m1.AssertExpectations(t)
			m1.AssertNumberOfCalls(t, "ExtendResponseWriter", 1)
			m2.AssertExpectations(t)
			m2.AssertNumberOfCalls(t, "ExtendResponseObject", 1)
		})

	})

}

func TestExtOmitObjectKey(t *testing.T) {
	t.Run("test empty", func(t *testing.T) {
		ext := ExtChain(
			ExtObjectKeyValue("key1", "value1"),
			ExtOmitObjectKey("key1"),
		)
		obj := make(map[string]any)
		assert.True(t, ext.ExtendResponseObject(context.Background(), obj))
	})
	t.Run("test existing", func(t *testing.T) {
		ext := ExtOmitObjectKey("key1")
		obj := make(map[string]any)
		assert.False(t, ext.ExtendResponseObject(context.Background(), obj))
	})
}

func TestOmitSettingsKey(t *testing.T) {
	t.Run("test empty", func(t *testing.T) {
		ext := ExtChain(
			ExtOmitSettingsKey(
				func(s Settings) []string {
					return []string{s.DefaultErrorStatusCodeKey}
				},
			),
		)

		ctx := contextWithSettingsValue(context.Background(), DefaultSettings())
		obj := make(map[string]any)
		ext.ExtendResponseObject(ctx, obj)
		assert.Empty(t, obj)
	})
	t.Run("test existing", func(t *testing.T) {
		ext := ExtChain(
			extSettingsKeyValue(func(s Settings) string {
				return s.DefaultErrorStatusCodeKey
			}, "hello"),
			ExtOmitSettingsKey(
				func(s Settings) []string {
					return []string{s.DefaultErrorStatusCodeKey}
				},
			),
		)

		ctx := contextWithSettingsValue(context.Background(), DefaultSettings())
		obj := make(map[string]any)
		ext.ExtendResponseObject(ctx, obj)
		assert.Empty(t, obj)
	})
}
