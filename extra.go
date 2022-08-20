package sendo

// ExtraOptionToken provides a token as ane extra option.
type ExtraOptionToken struct {
	Tokens string
}

// extraOption extracts an extra options from the list if it's provided, otherwise returns its zero value.
func extraOption[T any](extra []interface{}) (t T) {
	for _, o := range extra {
		if t, ok := o.(T); ok {
			return t
		}
	}

	return
}
