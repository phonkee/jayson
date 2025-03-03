package jayson

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestSettings_Validate(t *testing.T) {
	s := Settings{}
	s.Validate()

	assert.Equal(t, http.StatusInternalServerError, s.DefaultErrorStatus)
	assert.Equal(t, "message", s.DefaultErrorMessageKey)
	assert.Equal(t, "code", s.DefaultErrorStatusCodeKey)
	assert.Equal(t, "status", s.DefaultErrorStatusTextKey)
	assert.Equal(t, http.StatusOK, s.DefaultResponseStatus)
}
