# jayson

Simple json response/error writer on steroids.

# Design

Jayson defines Extension as interface that does 2 things
    - it can extend http response writer with additional headers, status codes, etc.
    - it can extend object with additional fields

If value is passed directly to `Response`, jayson just marshal it to json and writes it to response writer.
However if the value is `Extension`, jayson prepares object and gives ability to alter its keys and values.
We will show this behavior in [Usage](#Usage) section.

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
    jayson.Response(r.Context(), w, User{})
}

func HandlerCustomStatus(rw http.ResponseWriter, r *http.Request) {
    // Write response with custom status code
    jayson.Response(r.Context(), w, User{}, jayson.Status(http.StatusCreated))
}
```

## Simple Error

When you only want to write error with default status code it's as easy as:

```go
func HandlerSimpleError(rw http.ResponseWriter, r *http.Request) {
    // Write response with http.StatusInternalServerError status code
    jayson.Error(r.Context(), w, errors.New("boom internal error"))
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
        jayson.RegisterError(ErrNotFound, jayson.ExtStatus(http.StatusNotFound))
    )   
}

func HandlerAdvancedError(rw http.ResponseWriter, r *http.Request) {
    // Write error with http.StatusNotFound status code
    jayson.Error(r.Context(), w, ErrNotFound)
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
        jayson.RegisterError(ErrNotFound, jayson.ExtStatus(http.StatusNotFound)),
    )
}


func HandlerUnwrap(rw http.ResponseWriter, r *http.Request) {
    // Write error with http.StatusNotFound status code (inherits from ErrNotFound) and additional error detail
    jayson.Error(r.Context(), w, ErrUserNotFound, jayson.ExtErrorDetail("user not found"))
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
        jayson.RegisterError(Error1, jayson.ExtStatus(http.StatusNotFound)),
        jayson.RegisterError(Error2, jayson.ExtErrorDetail("error2 detail")),
    )
}

func Handler(rw http.ResponseWriter, r *http.Request) {
    // Write error with http.StatusNotFound status code and error detail "error2 detail"
    jayson.Error(r.Context(), w, Error2)
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
    RegisterResponse(&CreatedUser{}, Status(http.StatusCreated))
}

// CreateUser http handler
func CreateUser(w http.ResponseWriter, r *http.Request) {
    // Write response
    Response(r.Context(), w, CreatedUser{}, jayson.Header("X-Request-ID", "123"))
}
```

## Custom jayson instance

For special cases where you want to maintain separate jayson instance you can create your own instance.

```go
type CreatedUser struct {
    ID string `json:"id"`
}

var (
    jay = jayson.New(jayson.DefaultSettiongs())
    jay.RegisterResponse(&CreatedUser{}, jayson.Status(http.StatusCreated))
)

func Handler(w http.ResponseWriter, r *http.Request) {
    jay.Response(r.Context(), w, CreatedUser{})
}
```

# Author

Peter Vrba <phonkee@phonkee.eu>