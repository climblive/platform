package repository

func emptyAsNil[T comparable](value T) *T {
	var empty T

	if value == empty {
		return nil
	}

	return &value
}

func nilAsEmpty[T comparable](value *T) T {
	var empty T

	if value == nil {
		return empty
	}

	return *value
}
