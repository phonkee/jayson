package tester

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

// APIObject gives ability to unmarshall json object into multiple fields
func APIObject(t *testing.T, kv ...any) any {
	if len(kv)%2 != 0 {
		t.Fatal("APIObject: odd number of arguments")
	}

	m := make(map[string]any)

	// iterate over key-value pairs
	for i := 0; i < len(kv); i += 2 {
		k, ok := kv[i].(string)
		assert.Truef(t, ok, "APIObject: key `%s` is not a string", kv[i])
		val := reflect.ValueOf(kv[i+1])
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		require.Truef(t, val.CanSet(), "APIObject: value `%s` is not settable", k)
		m[k] = kv[i+1]
	}

	return &object{
		kv: m,
		T:  t,
	}
}

// object is a helper struct for unmarshalling json object into multiple fields
type object struct {
	kv map[string]any
	T  *testing.T
}

func (o *object) UnmarshalJSON(bytes []byte) error {
	into := make(map[string]json.RawMessage)
	err := json.Unmarshal(bytes, &into)
	require.NoError(o.T, err)

	for k, v := range o.kv {
		raw, ok := into[k]
		require.Truef(o.T, ok, "APIObject: key `%s` not found", k)

		err := json.Unmarshal(raw, v)
		require.NoError(o.T, err)
	}

	return nil
}
