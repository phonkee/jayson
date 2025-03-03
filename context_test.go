package jayson

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCtxErrorValue(t *testing.T) {
	ctx := context.Background()
	err, ok := ContextErrorValue(ctx)
	assert.False(t, ok)
	assert.Nil(t, err)

	ctx = contextWithErrorValue(ctx, assert.AnError)
	err, ok = ContextErrorValue(ctx)
	assert.True(t, ok)
	assert.Equal(t, assert.AnError, err)
}

func TestCtxSettingsValue(t *testing.T) {
	ctx := context.Background()
	settings := ContextSettingsValue(ctx)
	assert.Zero(t, settings)

	def := DefaultSettings()
	ctx = contextWithSettingsValue(ctx, def)
	settings = ContextSettingsValue(ctx)
	assert.Equal(t, def, settings)
}

func TestCtxObjectValue(t *testing.T) {
	ctx := context.Background()
	obj, ok := ContextObjectValue[int](ctx)
	assert.False(t, ok)
	assert.Zero(t, obj)

	ctx = contextWithObjectValue(ctx, 42)
	obj, ok = ContextObjectValue[int](ctx)
	assert.True(t, ok)
	assert.Equal(t, 42, obj)
}
