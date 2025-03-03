package jayson

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCallerInfo(t *testing.T) {
	t.Run("test newCallerInfo", func(t *testing.T) {
		ci := newCallerInfo(100)
		assert.NotEmpty(t, ci.file)
		assert.NotEmpty(t, ci.fn)
		assert.NotZero(t, ci.line)
	})
}
