package jayson

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMust(t *testing.T) {
	assert.Panics(t, func() {
		Must(errors.New("test"))
	})
	assert.Panics(t, func() {
		Must(nil, nil, errors.New("test"))
	})
	assert.NotPanics(t, func() {
		Must(nil)
	})
	assert.NotPanics(t, func() {
		Must()
	})
}
