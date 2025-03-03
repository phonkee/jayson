package tester

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// response is the implementation of APIResponse
type response struct {
	rw      *httptest.ResponseRecorder
	request *http.Request
	body    []byte
}

// AssertJsonEquals asserts that response body is equal to given json string or object/map
func (r *response) AssertJsonEquals(t *testing.T, expected any) APIResponse {
	require.NotNilf(t, r.body, "response body is nil")

	var expectedStr string

	switch expected.(type) {
	case string:
		expectedStr = expected.(string)
	case []byte:
		expectedStr = string(expected.([]byte))
	default:
		b, err := json.Marshal(expected)
		require.NoErrorf(t, err, "failed to marshal expected value: %v", expected)
		expectedStr = string(b)
	}

	require.JSONEq(t, expectedStr, string(r.body))
	return r
}

// AssertJsonKeyEquals asserts that response body key is equal to given value
// This method uses bit of magic to unmarshal json object into given value.
// It inspects the type of the given value and unmarshalls the json object into same type.
func (r *response) AssertJsonKeyEquals(t *testing.T, key string, what any) APIResponse {
	require.NotNilf(t, r.body, "response body is nil")

	var val reflect.Value
	typ := reflect.TypeOf(what)
	if typ.Kind() == reflect.Ptr {
		val = reflect.New(typ.Elem())
	} else {
		val = reflect.New(typ).Elem()
	}

	// prepare object
	obj := make(map[string]json.RawMessage)

	assert.NoError(t, json.NewDecoder(bytes.NewReader(r.body)).Decode(&obj))

	v, ok := obj[key]
	require.Truef(t, ok, "key %s not found in response", key)

	target := val.Interface()
	assert.NoError(t, json.NewDecoder(bytes.NewBuffer(v)).Decode(&target))

	assert.Equalf(t, what, target, "expected: %v, got: %v", what, target)

	return r
}

// AssertStatus asserts that response status is equal to given status
func (r *response) AssertStatus(t *testing.T, status int) APIResponse {
	require.Equal(t, status, r.rw.Code)
	return r
}

// Unmarshal unmarshalls whole response body into given value
func (r *response) Unmarshal(t *testing.T, v any) APIResponse {
	require.NotNilf(t, r.body, "response body is nil")
	require.NoError(t, json.NewDecoder(bytes.NewReader(r.body)).Decode(v))
	return r
}
