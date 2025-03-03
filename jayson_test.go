package jayson_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/phonkee/jayson"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	ErrorDetailKey     = "errorDetail"
	ErrorMessageKey    = "errorMessage"
	ErrorStatusCodeKey = "statusCode"
	ErrorStatusTextKey = "statusText"
)

func testSettings() jayson.Settings {
	return jayson.Settings{
		DefaultErrorStatus:        http.StatusInternalServerError,
		DefaultResponseStatus:     http.StatusOK,
		DefaultErrorMessageKey:    ErrorMessageKey,
		DefaultErrorStatusCodeKey: ErrorStatusCodeKey,
		DefaultErrorStatusTextKey: ErrorStatusTextKey,
	}
}

// extErrorDetail is an extFunc that adds an error detail to the response object.
func extErrorDetail(detail string) jayson.Extension {
	return jayson.ExtObjectKeyValue("errorDetail", detail)
}

func TestNew(t *testing.T) {
	j := jayson.New(testSettings())
	assert.NotNil(t, j)
}

var (
	// Error1 is test error
	Error1 = errors.New("error1")
	Error2 = fmt.Errorf("error2: %w", Error1)
	Error3 = fmt.Errorf("error3: %w", Error2)
)

func assertErrorJSON(t *testing.T, jay jayson.Jayson, err error, expect string, expectStatus int, header http.Header) {
	// response instance
	rw := httptest.NewRecorder()

	// write error to response
	jay.Error(context.Background(), rw, err)

	// assert status code
	assert.Equal(t, expectStatus, rw.Result().StatusCode)

	// print response now
	zap.L().Debug("response", zap.String("body", rw.Body.String()))

	assert.JSONEqf(t, expect, rw.Body.String(), "expected: %s, got: %s", expect, rw.Body.String())

	if header != nil {
		assert.Equal(t, header, rw.Header())
	}
}

func assertResponseJSON(t *testing.T, jay jayson.Jayson, obj any, expect string, expectStatus int, header http.Header) {
	// response instance
	rw := httptest.NewRecorder()

	// write object to response
	jay.Response(context.Background(), rw, obj)

	// print response now
	zap.L().Debug("response", zap.String("body", rw.Body.String()))

	assert.Equalf(t, expectStatus, rw.Code, "expected: %d, got: %d", expectStatus, rw.Code)

	assert.JSONEqf(t, expect, rw.Body.String(), "expected: %s, got: %s", expect, rw.Body.String())

	if header != nil {
		assert.Equal(t, header, rw.Header())
	}
}

func TestJayson_RegisterError(t *testing.T) {

	t.Run("test register nil error panics", func(t *testing.T) {
		jay := jayson.New(testSettings())
		assert.Panics(t, func() {
			_ = jay.RegisterError(nil, extErrorDetail("woah this is bad"))
		})
	})

	t.Run("test register Any error sets defaults for unknown error", func(t *testing.T) {
		jay := jayson.New(testSettings())
		assert.NoError(t,
			jay.RegisterError(jayson.Any, extErrorDetail("woah this is bad")),
		)

		assertErrorJSON(t, jay, Error3, `{"`+ErrorStatusCodeKey+`":500,"`+ErrorDetailKey+`":"woah this is bad","`+ErrorMessageKey+`":"error3: error2: error1","`+ErrorStatusTextKey+`":"Internal Server Error"}`, http.StatusInternalServerError, nil)
	})

	t.Run("should register extFunc for given error", func(t *testing.T) {
		jay := jayson.New(testSettings())

		// register error
		assert.NoError(t,
			jay.RegisterError(Error1, jayson.ExtObjectKeyValue("answer_to_anything", "42"), jayson.ExtStatus(http.StatusTeapot)),
		)
		assert.NoError(t,
			jay.RegisterError(Error2, jayson.ExtStatus(http.StatusNotFound)),
		)
		assert.NoError(t,
			jay.RegisterError(Error3,
				extErrorDetail("some error"),
				jayson.ExtHeaderValue("X-Hello", "World"),
				jayson.ExtNoop(),
				jayson.ExtFunc(nil, nil),
				jayson.ExtHeader(http.Header{"X-Welcome": []string{"Here"}}),
			),
		)

		assertErrorJSON(t, jay, Error3, `{"answer_to_anything":"42","`+ErrorStatusCodeKey+`":404,"`+ErrorDetailKey+`":"some error","`+ErrorMessageKey+`":"error3: error2: error1","`+ErrorStatusTextKey+`":"Not Found"}`, http.StatusNotFound, nil)
	})

	t.Run("test unknown error", func(t *testing.T) {
		jay := jayson.New(testSettings())
		assert.NoError(t,
			jay.RegisterError(jayson.Any, extErrorDetail("woah this is bad")),
		)
		assert.NoError(t,
			jay.RegisterError(Error3, jayson.ExtStatus(http.StatusTeapot)),
		)

		assertErrorJSON(t, jay, Error2, `{"`+ErrorStatusCodeKey+`":500,"`+ErrorDetailKey+`":"woah this is bad","`+ErrorMessageKey+`":"error2: error1","`+ErrorStatusTextKey+`":"Internal Server Error"}`, http.StatusInternalServerError, nil)
	})

	t.Run("test inherit extensions from wrapped error", func(t *testing.T) {
		jay := jayson.New(testSettings())
		assert.NoError(t,
			jay.RegisterError(Error1, jayson.ExtStatus(http.StatusTeapot)),
		)
		assertErrorJSON(t, jay, Error2, `{"`+ErrorStatusCodeKey+`":418,"`+ErrorMessageKey+`":"error2: error1","`+ErrorStatusTextKey+`":"I'm a teapot"}`, http.StatusTeapot, nil)
	})

}

type testResponse struct {
	Answer int `json:"answer"`
}

func TestJayson_Response(t *testing.T) {
	t.Run("test pointer response with response", func(t *testing.T) {
		jay := jayson.New(testSettings())
		assert.NoError(t,
			jay.RegisterResponse(
				&testResponse{},
				jayson.ExtStatus(http.StatusTeapot),
				jayson.ExtHeaderValue("X-Hello", "World"),
			),
		)

		assertResponseJSON(t, jay,
			&testResponse{Answer: 42},
			`{"answer":42}`,
			http.StatusTeapot,
			http.Header{
				"Content-Type": []string{"application/json"},
				"X-Hello":      []string{"World"},
			},
		)
		assertResponseJSON(t, jay, testResponse{Answer: 42}, `{"answer":42}`, http.StatusTeapot, nil)
	})

	t.Run("test json marshal panics", func(t *testing.T) {
		jay := jayson.New(testSettings())
		assert.Panics(t, func() {
			rw := httptest.NewRecorder()
			jay.Response(context.Background(), rw, SomeWrongType(1))
		})

		assert.Panics(t, func() {
			rw := httptest.NewRecorder()
			jay.Response(context.Background(), rw, jayson.ExtObjectKeyValue("key", SomeWrongType(1)))
		})

	})

	t.Run("test registered pointer with value", func(t *testing.T) {
		jay := jayson.New(testSettings())
		assert.NoError(t,
			jay.RegisterResponse(
				testResponse{},
				jayson.ExtStatus(http.StatusTeapot),
				jayson.ExtHeaderValue("X-Hello", "World"),
			),
		)

		assertResponseJSON(t, jay, &testResponse{Answer: 42}, `{"answer":42}`, http.StatusTeapot, nil)

	})

}

type SomeWrongType int

func (s SomeWrongType) MarshalJSON() ([]byte, error) {
	return nil, errors.New("marshal error")
}

func TestJayson_RegisterResponse_AnyResponse(t *testing.T) {
	jay := jayson.New(testSettings())
	assert.NoError(t,
		jay.RegisterResponse(jayson.Any, jayson.ExtStatus(http.StatusTeapot)),
	)
}

func TestJayson_Debug(t *testing.T) {
	t.Run("test enable debug", func(t *testing.T) {
		t.Run("test register response", func(t *testing.T) {
			observedZapCore, observedLogs := observer.New(zap.DebugLevel)
			observedLogger := zap.New(observedZapCore)

			jay := jayson.New(testSettings())
			jay.Debug(observedLogger)

			jayson.Must(
				jay.RegisterResponse(jayson.Any, jayson.ExtStatus(http.StatusTeapot)),
			)

			assert.Len(t, observedLogs.All(), 1)
		})

		t.Run("test register error", func(t *testing.T) {
			observedZapCore, observedLogs := observer.New(zap.DebugLevel)
			observedLogger := zap.New(observedZapCore)

			jay := jayson.New(testSettings())
			jay.Debug(observedLogger)

			jayson.Must(
				jay.RegisterError(jayson.Any, jayson.ExtStatus(http.StatusTeapot)),
			)

			assert.Len(t, observedLogs.All(), 1)

		})

	})

}

func TestJayson_Error(t *testing.T) {
	t.Run("empty error does nothing", func(t *testing.T) {
		jay := jayson.New(testSettings())
		rw := httptest.NewRecorder()
		jay.Error(context.Background(), rw, nil)
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Empty(t, rw.Body.String())
	})

	t.Run("test marshal error panics", func(t *testing.T) {
		jay := jayson.New(testSettings())
		rw := httptest.NewRecorder()
		assert.Panics(t, func() {
			jay.Error(context.Background(), rw, errors.New("error"), jayson.ExtObjectKeyValue("key", SomeWrongType(1)))
		})

		println(rw.Body.String())

	})

}
