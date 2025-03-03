package jayson

import (
	"context"
	"go.uber.org/zap"
	"net/http"
)

var (
	// _g is a global instance of Jayson
	_g = New(DefaultSettings())
)

// Debug enables debug mode via zap logger in _g instance of Jayson
func Debug(logger *zap.Logger) {
	_g.Debug(logger)
}

// Error writes error to _g instance of Jayson
func Error(ctx context.Context, w http.ResponseWriter, err error, ext ...Extension) {
	_g.Error(ctx, w, err, ext...)
}

// RegisterError registers extFunc for given error in _g instance of Jayson
func RegisterError(err error, ext ...Extension) error {
	return _g.RegisterError(err, ext...)
}

// RegisterResponse registers extFunc for given object in _g instance of Jayson
func RegisterResponse(what any, ext ...Extension) error {
	return _g.RegisterResponse(what, ext...)
}

// Response writes response to http.ResponseWriter using _g instance of Jayson
func Response(ctx context.Context, w http.ResponseWriter, what any, ext ...Extension) {
	_g.Response(ctx, w, what, ext...)
}
