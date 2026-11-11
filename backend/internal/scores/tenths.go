package scores

import "math"

type Tenths int64

func tenthsFromFloat64(points float64) Tenths {
	return Tenths(math.Round(points * 10))
}

func (t Tenths) Float64() float64 {
	return float64(t) / 10
}
