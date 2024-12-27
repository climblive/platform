package domain

import "encoding/json"

type Patch[T any] struct {
	Value T
}

func NewPatch[T any](v T) *Patch[T] {
	return &Patch[T]{Value: v}
}

func (p *Patch[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Value)
}

func (p *Patch[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &p.Value)
}
