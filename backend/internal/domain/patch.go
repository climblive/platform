package domain

import (
	"encoding/json"

	"github.com/go-errors/errors"
)

type Patch[T any] struct {
	Present bool
	Value   T
}

func NewPatch[T any](v T) Patch[T] {
	return Patch[T]{
		Present: true,
		Value:   v,
	}
}

func (p Patch[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Value)
}

func (p *Patch[T]) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &p.Value)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	p.Present = true

	return nil
}

func (p *Patch[T]) IsZero() bool {
	return !p.Present
}
