package repository

import "github.com/climblive/platform/backend/internal/domain"

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

func nillableIntToResourceID[T domain.ResourceIDType](value *int) *T {
	if value == nil {
		return nil
	}

	var out T = T(*value)
	return &out
}
