package jayson

// anyTarget is a type used to register any response/error object.
// it is used as base type for all response/error objects
type anyTarget int

// Error implements error interface to be used in RegisterError
func (a anyTarget) Error() string { return "" }

var (
	// Any is used to RegisterResponse and RegisterError for any response/error object,
	// It is used as base type for all response/error objects
	// regardless if it's registered or not.
	Any = anyTarget(0)
)
