package sendo

type ExtraOptionsToken struct {
	Tokens string
}

func extraOption[T any](extra []interface{}) (t T) {
	for _, o := range extra {
		if t, ok := o.(T); ok {
			return t
		}
	}

	return
}
