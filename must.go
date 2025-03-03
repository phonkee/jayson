package jayson

// Must panics if err is not nil, otherwise it returns the value.
func Must(err ...error) {
	for _, e := range err {
		if e != nil {
			panic(e)
		}
	}
}
