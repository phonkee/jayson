package jayson

import (
	"context"
	"fmt"
	"net/http"
)

// ExtChain returns an extFunc that chains multiple extensions together.
func ExtChain(extensions ...Extension) Extension {
	return ExtFunc(
		func(ctx context.Context, w http.ResponseWriter) (result bool) {
			for _, e := range extensions {
				if e.ExtendResponseWriter(ctx, w) {
					result = true
				}
			}
			return result
		},
		func(ctx context.Context, m map[string]any) (result bool) {
			for _, e := range extensions {
				if e.ExtendResponseObject(ctx, m) {
					result = true
				}
			}
			return result
		},
	)
}

// ExtConditional calls first extFunc, and if it returns true, it calls all the extensions.
// This is useful for conditional extensions based on context (such as debug mode or any context values).
func ExtConditional(condition Extension, ext ...Extension) Extension {
	return ExtFunc(
		func(ctx context.Context, w http.ResponseWriter) bool {
			var result bool
			if condition.ExtendResponseWriter(ctx, w) {
				for _, e := range ext {
					if e.ExtendResponseWriter(ctx, w) {
						result = true
					}
				}
			}
			return result
		},
		func(ctx context.Context, m map[string]any) bool {
			var result bool
			if condition.ExtendResponseObject(ctx, m) {
				for _, e := range ext {
					if e.ExtendResponseObject(ctx, m) {
						result = true
					}
				}
			}
			return result
		},
	)
}

// ExtErrorDetail is an extFunc that adds an error detail to the response object.
func ExtErrorDetail(detail string) Extension {
	return extSettingsKeyValue(func(s Settings) string {
		return s.DefaultErrorDetailKey
	}, detail)
}

// ExtFirst returns an extFunc that returns the first extFunc that extends the response.
func ExtFirst(ext ...Extension) Extension {
	return ExtFunc(
		func(ctx context.Context, w http.ResponseWriter) bool {
			for _, e := range ext {
				if e.ExtendResponseWriter(ctx, w) {
					return true
				}
			}
			return false
		},
		func(ctx context.Context, m map[string]any) bool {
			for _, e := range ext {
				if e.ExtendResponseObject(ctx, m) {
					return true
				}
			}
			return false
		},
	)
}

// ExtFunc is an extFunc that calls the given function to extend the response or the response object.
func ExtFunc(
	extResponseWriter ExtensionApplyResponseWriterFunc,
	extResponseObject ExtensionApplyResponseObjectFunc,
) Extension {
	return &extFunc{
		extendResponseWriter: func(ctx context.Context, writer http.ResponseWriter) bool {
			if extResponseWriter != nil {
				return extResponseWriter(ctx, writer)
			}
			return false
		},
		extendResponseObject: func(ctx context.Context, m map[string]any) bool {
			if extResponseObject != nil {
				return extResponseObject(ctx, m)
			}
			return false
		},
	}
}

// ExtHeader extends the response with the given HTTP headers.
func ExtHeader(h http.Header) Extension {
	return ExtFunc(
		func(ctx context.Context, w http.ResponseWriter) bool {
			if h == nil {
				return false
			}
			for k, v := range h {
				for _, vv := range v {
					w.Header().Add(k, vv)
				}
			}
			return true
		},
		nil,
	)
}

// ExtHeaderValue is an extFunc that adds a single HTTP header to the response.
func ExtHeaderValue(k string, v string) Extension {
	return ExtFunc(
		func(ctx context.Context, w http.ResponseWriter) bool {
			w.Header().Add(k, v)
			return true
		},
		nil,
	)
}

// ExtNoop does not extend the response or the response object.
func ExtNoop() Extension {
	return ExtFunc(nil, nil)
}

// ExtObjectKeyValue is an extFunc that adds a single key-value pair to the response object.
func ExtObjectKeyValue(key string, value any) Extension {
	return ExtFunc(
		nil,
		func(ctx context.Context, m map[string]any) bool {
			m[key] = value
			return true
		},
	)
}

// ExtObjectKeyValuef is an extFunc that adds a single key-value pair to the response object based on a format string.
func ExtObjectKeyValuef(key string, format string, args ...any) Extension {
	return ExtFunc(
		nil,
		func(ctx context.Context, m map[string]any) bool {
			m[key] = fmt.Sprintf(format, args...)
			return true
		},
	)
}

// ExtOmitObjectKey is an extFunc that removes the given keys from the response object.
// This extFunc needs to be called at the end,
// so it removes the keys after all other extensions have added their keys.
func ExtOmitObjectKey(keys ...string) Extension {
	return ExtFunc(
		nil,
		func(ctx context.Context, m map[string]any) (result bool) {
			for _, key := range keys {
				if _, ok := m[key]; ok {
					delete(m, key)
					result = true
				}
			}
			return result
		},
	)
}

// ExtStatus is an extension that sets the HTTP status code of the response.
func ExtStatus(status int) Extension {
	return ExtFunc(
		func(ctx context.Context, w http.ResponseWriter) bool {
			w.WriteHeader(status)
			return true
		},
		nil,
	)
}

// extSettingsKeyValue is an extFunc that adds a single key-value pair to the response object based on the settings.
func extSettingsKeyValue(fn func(s Settings) string, value any) Extension {
	return ExtFunc(
		nil,
		func(ctx context.Context, m map[string]any) bool {
			s := ContextSettingsValue(ctx)
			if key := fn(s); key != "" {
				m[key] = value
			}
			return true
		},
	)
}

// extOmitSettingsKey is an extFunc that removes the given keys from the response object based on the settings.
func extOmitSettingsKey(fn func(settings Settings) []string) Extension {
	return ExtFunc(
		nil,
		func(ctx context.Context, m map[string]any) (result bool) {
			s := ContextSettingsValue(ctx)
			for _, key := range fn(s) {
				if _, ok := m[key]; ok {
					delete(m, key)
					result = true
				}
			}
			return result
		},
	)
}

// extFunc is a generic extFunc that can extend the response or the response object.
// it is used to implement the Extension interface.
type extFunc struct {
	extendResponseWriter func(context.Context, http.ResponseWriter) bool
	extendResponseObject func(context.Context, map[string]any) bool
}

// ExtendResponseWriter calls the extFunc's extendResponseWriter function.
func (e *extFunc) ExtendResponseWriter(ctx context.Context, w http.ResponseWriter) bool {
	return e.extendResponseWriter(ctx, w)
}

// ExtendResponseObject calls the extFunc's extendResponseObject function.
func (e *extFunc) ExtendResponseObject(ctx context.Context, m map[string]any) bool {
	return e.extendResponseObject(ctx, m)
}
