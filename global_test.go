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
