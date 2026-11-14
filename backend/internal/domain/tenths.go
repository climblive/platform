package domain

import (
	"encoding/json/v2"
	"math"
	"strconv"

	"github.com/go-errors/errors"
)

func TenthsFromFloat64(points float64) Tenths {
	return Tenths(math.Round(points * 10))
}

func (t Tenths) Float64() float64 {
	return float64(t) / 10
}

func (t Tenths) String() string {
	whole, fraction := int64(t/10), int64(t%10)
	data := make([]byte, 0, 21)
	if t < 0 {
		data = append(data, '-')
		whole = -whole
		fraction = -fraction
	}
	data = strconv.AppendInt(data, whole, 10)
	if fraction != 0 {
		data = append(data, '.')
		data = strconv.AppendInt(data, fraction, 10)
	}
	return string(data)
}

func (t Tenths) MarshalJSON() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t *Tenths) UnmarshalJSON(data []byte) error {
	var points float64
	if err := json.Unmarshal(data, &points); err != nil {
		return errors.Wrap(err, 0)
	}

	scaled := points * 10
	rounded := math.Round(scaled)
	if rounded < -0x1p63 || rounded >= 0x1p63 || math.Abs(scaled-rounded) >= 1e-6 {
		return errors.Errorf("%w: value cannot be represented in tenths", ErrInvalidData)
	}

	*t = Tenths(rounded)
	return nil
}
