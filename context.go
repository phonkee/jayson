package jayson

import "context"

// contextKey is a type for storing values in context.
type contextKey int

const (
	// contextSettingsKey is the key used to store the settings value in the context.
	contextSettingsKey contextKey = 1 << iota

	// contextErrorKey is the key used to store the error value in the context.
	contextErrorKey

	// contextObjectKey is the key used to store the object value in the context.
	contextObjectKey
)

// ContextErrorValue returns the error value stored in the context.
func ContextErrorValue(ctx context.Context) (error, bool) {
	err, ok := ctx.Value(contextErrorKey).(error)
	return err, ok
}

// ContextObjectValue returns the object value stored in the context.
func ContextObjectValue[T any](ctx context.Context) (T, bool) {
	val, ok := ctx.Value(contextObjectKey).(T)
	return val, ok
}

// ContextSettingsValue returns the settings value stored in the context.
func ContextSettingsValue(ctx context.Context) Settings {
	var result Settings
	v, ok := ctx.Value(contextSettingsKey).(Settings)
	if ok {
		result = v
	}
	return result
}

// contextWithErrorValue adds the error value to the context.
func contextWithErrorValue(ctx context.Context, err error) context.Context {
	return context.WithValue(ctx, contextErrorKey, err)
}

// contextWithObjectValue adds the object value to the context.
func contextWithObjectValue(ctx context.Context, obj any) context.Context {
	return context.WithValue(ctx, contextObjectKey, obj)
}

// contextWithSettingsValue adds the settings value to the context.
func contextWithSettingsValue(ctx context.Context, settings Settings) context.Context {
	return context.WithValue(ctx, contextSettingsKey, settings)
}
