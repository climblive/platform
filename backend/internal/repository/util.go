package repository

func e2n[T comparable](value T) *T {
	var empty T

	if value == empty {
		return nil
	}

	return &value
}

func n2e[T comparable](value *T) T {
	var empty T

	if value == nil {
		return empty
	}

	return *value
}
