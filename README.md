# jayson

Simple json response/error writer on steroids.

# Design

Jayson defines Extension as interface that does 2 things
    - it can extend http response writer with additional headers, status codes, etc.
    - it can extend object with additional fields


## Errors

Jayson provides way to register errors and provide `extensions` to them.
When error is written to response writer, jayson looks up error and all its ancestors in registry and applies all extensions to it.
This gives ability to define error once and use it in multiple places with different extensions.
You can build error hierarchy and add extensions to each error in hierarchy.

## Responses

Responses are similar to errors, but they are not errors. You can register responses and provide extensions to them.
However when response is passed as object, all extensions that alter "object" can not do anything with it.
It will just run extensions that alter response writer (status, headers).
There is a way to provide extensions to object, but you need to wrap it in `jayson.ExtObjectKeyValue` which marks it as key of reponse object.
This makes it easier to provide multiple extensions to single object.

# Global state

Jayson provides global instance of jayson that can be used to register errors, responses and extensions.
You can swap global instance with your own instance if you need to maintain separate jayson instance.
You can also instantiate your own jayson instance and use it in your code.

# Usage

## Simple response

Jayson can directly marshal values provided to `Response` and write them to response writer.
In this case jayson does not do any additional processing of object.

```go
type User struct {
    ID string `json:"id"`
}

func Handler(rw http.ResponseWriter, r *http.Request) {
    // Write response with default status code
    jayson.G().Response(r.Context(), w, User{})
}

func HandlerCustomStatus(rw http.ResponseWriter, r *http.Request) {
    // Write response with custom status code
    jayson.G().Response(r.Context(), w, User{}, jayson.ExtStatus(http.StatusCreated))
}
```

## Simple Error

When you only want to write error with default status code it's as easy as:

```go
func HandlerSimpleError(rw http.ResponseWriter, r *http.Request) {
    // Write response with http.StatusInternalServerError status code
    jayson.G().Error(r.Context(), w, errors.New("boom internal error"))
}
```

## Advanced Error usage

You can register errors with status codes and other extensions.

```go
var (
    ErrNotFound = errors.New("not found")
)

func init() {
    // register error with status code
    jayson.Must(
        jayson.G().RegisterError(ErrNotFound, jayson.ExtStatus(http.StatusNotFound))
    )   
}

func HandlerAdvancedError(rw http.ResponseWriter, r *http.Request) {
    // Write error with http.StatusNotFound status code
    jayson.G().Error(r.Context(), w, ErrNotFound)
}
```

You can also add additional extensions to error

```go
var (
    // ErrNotFound is generic not found error
    ErrNotFound = errors.New("not found")
    // ErrUserNotFound is returned when user is not found, it wraps ErrNotFound
    ErrUserNotFound = fmt.Errorf("%w: user", ErrNotFound)
)

func init() {
    jayson.Must(
        jayson.G().RegisterError(ErrNotFound, jayson.ExtStatus(http.StatusNotFound)),
    )
}


func HandlerUnwrap(rw http.ResponseWriter, r *http.Request) {
    // Write error with http.StatusNotFound status code (inherits from ErrNotFound) and additional error detail
    jayson.G().Error(r.Context(), w, ErrUserNotFound, jayson.ExtErrorDetail("user not found"))
}
```

Error can also be wrapped with additional extensions, and they inherit all extensions from parent error

```go
var (
    Error1 = errors.New("error1")
    Error2 = fmt.Errorf("%w: error2", Error1)
)

func init() {
    jayson.Must(
        jayson.G().RegisterError(Error1, jayson.ExtStatus(http.StatusNotFound)),
        jayson.G().RegisterError(Error2, jayson.ExtErrorDetail("error2 detail")),
    )
}

func Handler(rw http.ResponseWriter, r *http.Request) {
    // Write error with http.StatusNotFound status code and error detail "error2 detail"
    jayson.G().Error(r.Context(), w, Error2)
}
```

## Advanced Response usage

```go

// CreatedUser is returned when user is created
type CreatedUser struct {
    ID string `json:"id"`
}

// Register all errors, structures with status and other extensions
func init() {
    // RegisterResponse registers response with status code (also non pointer type if not registered already)
    jayson.G().RegisterResponse(&CreatedUser{}, Status(http.StatusCreated))
}

// CreateUser http handler
func CreateUser(w http.ResponseWriter, r *http.Request) {
    // Write response
    jayson.G().Response(r.Context(), w, CreatedUser{}, jayson.Header("X-Request-ID", "123"))
}
```

## Custom jayson instance

For special cases where you want to maintain separate jayson instance you can create your own instance.

```go
type CreatedUser struct {
    ID string `json:"id"`
}

var (
	// jay is custom jayson instance
    jay = jayson.New(jayson.DefaultSettings())
)

func init() {
    jay.RegisterResponse(&CreatedUser{}, jayson.ExtStatus(http.StatusCreated))
}

func Handler(w http.ResponseWriter, r *http.Request) {
    jay.Response(r.Context(), w, CreatedUser{})
}
```

# Author

Peter Vrba <phonkee@phonkee.eu>