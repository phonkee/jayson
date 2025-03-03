package jayson

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

// Extension is the interface for Jayson extensions.
//
// It is used to alter the response or the object before it is sent to the client.
//
//go:generate mockery --name Extension --filename extension.go
type Extension interface {
	// ExtendResponseWriter extends the response
	ExtendResponseWriter(context.Context, http.ResponseWriter) bool
	// ExtendResponseObject extends the object
	ExtendResponseObject(context.Context, map[string]any) bool
}

// ExtensionApplyResponseWriterFunc is a function that extends the response writer.
type ExtensionApplyResponseWriterFunc func(context.Context, http.ResponseWriter) bool

// ExtensionApplyResponseObjectFunc is a function that extends the response object.
type ExtensionApplyResponseObjectFunc func(context.Context, map[string]any) bool

// Jayson is the interface that provides methods for writing responses to the client.
//
// By default, there is a global instance to be used,
// but for some special purposes you can create your own (multiple servers, different settings, etc.)
type Jayson interface {
	// Debug enables debug mode via zap logger.
	Debug(*zap.Logger)
	// Error writes error to the client.
	Error(context.Context, http.ResponseWriter, error, ...Extension)
	// RegisterError registers extFunc for given error.
	RegisterError(error, ...Extension) error
	// RegisterResponse registers extFunc for given response object.
	RegisterResponse(any, ...Extension) error
	// Response writes given object/error to the client.
	Response(context.Context, http.ResponseWriter, any, ...Extension)
}

var (
	// ErrImproperlyConfigured is error returned when Jayson is improperly configured.
	ErrImproperlyConfigured = errors.New("jayson: improperly configured")
	ErrAlreadyRegistered    = fmt.Errorf("%w: already registered", ErrImproperlyConfigured)
)

const (
	// DebugMaxCallerDepth is the maximum depth of debug caller
	DebugMaxCallerDepth = 255
)
