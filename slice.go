package jayson

// reverseSlice reverses the order of elements in a slice.
func reverseSlice[T any](s []T) []T {
	r := make([]T, len(s))
	for i := range s {
		r[len(s)-1-i] = s[i]
	}
	return r
}
