package jayson

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"reflect"
	"sync"
)

// New instantiates custom jayson instance. Usually you don't need to use it, since there is a _g instance.
func New(settings Settings) Jayson {
	// validate settings
	settings.Validate()

	return &jayson{
		settings: settings,
		errors:   make(map[error]*registeredError),
		types:    make(map[reflect.Type]*registeredType),
	}
}

// jayson implements Jayson interface
type jayson struct {
	debug              *zap.Logger
	settings           Settings
	anyErrorExtensions []Extension
	anyTypeExtensions  []Extension
	errors             map[error]*registeredError
	types              map[reflect.Type]*registeredType
	mutex              sync.RWMutex
}

// Debug enables debug mode via zap logger
func (j *jayson) Debug(logger *zap.Logger) {
	j.debug = logger
}

// Error writes error response to the client
func (j *jayson) Error(ctx context.Context, rw http.ResponseWriter, err error, extension ...Extension) {
	if err == nil {
		return
	}

	// mark original error
	originalError := err

	// get registered error
	var reg *registeredError

	// find registered error
	var ok bool

	// collect all registered errors
	allRegisteredErrors := make([]*registeredError, 0)

	// lock errors for read access
	j.mutex.RLock()
outer:
	for {
		// check if we have registered error
		if reg, ok = j.errors[err]; ok {
			allRegisteredErrors = append(allRegisteredErrors, reg)
		}

		// get parent error so we can check it
		err = errors.Unwrap(err)

		if err == nil {
			break outer
		}
	}
	j.mutex.RUnlock()

	// no error found, apply unknown error extensions
	if len(allRegisteredErrors) > 0 {
		// reverse errors
		allRegisteredErrors = reverseSlice(allRegisteredErrors)
	}

	// add error value to the context along with settings
	ctx = contextWithErrorValue(
		contextWithSettingsValue(ctx, j.settings),
		err,
	)

	// rwInternal
	obj := map[string]any{
		j.settings.DefaultErrorMessageKey: originalError.Error(),
	}
	rwInternal := newResponseWriter(j.settings.DefaultErrorStatus)

	// prepare extensions
	allExtensions := [][]Extension{
		j.anyErrorExtensions,
	}
	// add all found extensions (from registered errors)
	for _, reg := range allRegisteredErrors {
		allExtensions = append(allExtensions, reg.extensions)
	}
	// add extensions passed in function
	allExtensions = append(allExtensions, extension)

	// prepare executor
	exec := newExecutor(allExtensions...)

	// now extend response
	exec.ExtendResponseWriter(ctx, rwInternal)

	// now extend object
	exec.ExtendResponseObject(ctx, obj)

	// handle status code and text
	obj[j.settings.DefaultErrorStatusCodeKey] = rwInternal.statusCode
	if text := http.StatusText(rwInternal.statusCode); text != "" {
		obj[j.settings.DefaultErrorStatusTextKey] = text
	}

	// clear buffer here
	rwInternal.buffer = bytes.Buffer{}
	rwInternal.Header()["Content-Type"] = []string{"application/json"}

	// now write JSON value
	if err := json.NewEncoder(rwInternal).Encode(obj); err != nil {
		// if we can't write JSON, we will panic
		// it's fine now
		// FIXME: we should log this
		panic(err)
	}

	// write to a response writer
	rwInternal.WriteTo(rw)
}

// RegisterError registers extFunc for given error
func (j *jayson) RegisterError(err error, ext ...Extension) error {
	j.debugLogMethod("RegisterError", func() []zap.Field {
		return []zap.Field{
			zap.String("type", reflect.TypeOf(err).String()),
			zap.Int("extensions", len(ext)),
		}
	})

	j.mutex.Lock()
	defer j.mutex.Unlock()

	// if err is nil, this is error
	if err == nil {
		panic(fmt.Errorf("%w: error is nil", ErrImproperlyConfigured))
	}

	// if Any, we will register extensions for any error
	if errors.Is(err, Any) {
		j.anyErrorExtensions = append(j.anyErrorExtensions, ext...)
		return nil
	}

	// check if error is already registered
	// if so, we will update extensions but return error
	_, exists := j.errors[err]

	j.errors[err] = &registeredError{
		err:        err,
		extensions: ext,
	}

	// if error is already registered, return error
	if exists {
		return fmt.Errorf("%w: overwritten", ErrAlreadyRegistered)
	}

	return nil
}

// RegisterResponse registers extFunc for given object
func (j *jayson) RegisterResponse(what any, extensions ...Extension) error {

	// log caller
	j.debugLogMethod("RegisterResponse", func() []zap.Field {
		return []zap.Field{
			zap.Any("type", reflect.TypeOf(what).String()),
			zap.Int("extensions", len(extensions)),
		}
	})

	j.mutex.Lock()
	defer j.mutex.Unlock()

	// if what is Any, we will register extensions for any response object
	if what == Any {
		j.anyTypeExtensions = append(j.anyTypeExtensions, extensions...)
		return nil
	}

	// get type of what
	typ := reflect.TypeOf(what)

	// check if type exists
	_, exists := j.types[typ]

	// prepare registered type
	j.types[typ] = &registeredType{
		typ:        typ,
		extensions: extensions,
	}

	// now we want to add a special case for pointer types
	switch typ.Kind() {
	case reflect.Ptr:
		// if type is pointer, let's also register non-pointer value (if not already registered)
		typ = typ.Elem()
		if _, ok := j.types[typ]; !ok {
			j.types[typ] = &registeredType{
				typ:        typ,
				extensions: extensions,
			}
		}
	case reflect.Struct:
		// if type is struct, let's also register pointer value (if not already registered)
		typ = reflect.PointerTo(typ)
		if _, ok := j.types[typ]; !ok {
			j.types[typ] = &registeredType{
				typ:        typ,
				extensions: extensions,
			}
		}
	default:
		// do nothing
	}

	// if type is already registered, return error
	if exists {
		return fmt.Errorf("%w: overwritten", ErrAlreadyRegistered)
	}

	return nil
}

// Response writes response to the client
func (j *jayson) Response(ctx context.Context, rw http.ResponseWriter, what any, ext ...Extension) {
	// add object value to the context along with settings
	ctx = contextWithObjectValue(
		contextWithSettingsValue(ctx, j.settings),
		what,
	)

	// rwInternal is a response writer that will be used to collect response
	rwInternal := newResponseWriter(j.settings.DefaultResponseStatus)

	// if what is an ext, we will be having object automatically
	if extension, ok := what.(Extension); ok {
		// create object
		obj := make(map[string]any)

		// prepare executor
		exec := newExecutor(j.anyTypeExtensions, []Extension{extension}, ext)

		// extend response writer
		exec.ExtendResponseWriter(ctx, rwInternal)

		// extend object
		exec.ExtendResponseObject(ctx, obj)

		// now clear buffer
		rwInternal.buffer = bytes.Buffer{}

		// json marshal object
		if err := json.NewEncoder(rwInternal).Encode(obj); err != nil {
			panic(err)
		}
	} else {
		var (
			reg *registeredType
		)

		// get registered type
		j.mutex.RLock()
		reg, _ = j.types[reflect.TypeOf(what)]
		j.mutex.RUnlock()

		// prepare extensions that will be applied
		var regExtensions []Extension

		// if we found registered type, use its extensions
		if reg != nil {
			regExtensions = reg.extensions
		}

		// prepare executor with extensions (first what we found in registered types, then what is passed in function)
		exec := newExecutor(j.anyTypeExtensions, regExtensions, ext)

		// now extend response, no object here
		exec.ExtendResponseWriter(ctx, rwInternal)

		// assign buffer
		rwInternal.buffer = bytes.Buffer{}

		// now json encode object
		if err := json.NewEncoder(rwInternal).Encode(what); err != nil {
			panic(err)
		}
	}

	// set content type
	rwInternal.Header()["Content-Type"] = []string{"application/json"}
	rwInternal.WriteTo(rw)
}

// debugLogMethod logs caller info
func (j *jayson) debugLogMethod(method string, fn ...func() []zap.Field) {
	if j.debug == nil {
		return
	}

	// check if we are on debug level
	if ch := j.debug.Named("jayson").Check(zap.DebugLevel, "caller info"); ch != nil {
		// try to get caller info
		caller := newCallerInfo(DebugMaxCallerDepth)

		fields := []zap.Field{
			zap.String("method", method),
			zap.String("function", caller.fn),
			zap.Int("line", caller.line),
			zap.String("file", caller.file),
		}

		for _, fetch := range fn {
			fields = append(fields, fetch()...)
		}

		ch.Write(fields...)
	}
}

// registeredError holds error and its extensions
type registeredError struct {
	err        error
	extensions []Extension
}

// registeredType holds type and its extensions
type registeredType struct {
	typ        reflect.Type
	extensions []Extension
}
