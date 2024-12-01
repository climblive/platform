package scores

type DiffMap[K comparable, V any] struct {
	committed  map[K]V
	dirty      map[K]V
	comparator func(v1, v2 V) bool
}

func NewDiffMap[K comparable, V any](comparator func(v1, v2 V) bool) *DiffMap[K, V] {
	return &DiffMap[K, V]{
		committed:  make(map[K]V),
		dirty:      make(map[K]V),
		comparator: comparator,
	}
}

func (d *DiffMap[K, V]) Commit() []V {
	var diff []V

	for k, v := range d.dirty {
		diff = append(diff, v)
		d.committed[k] = v
	}

	clear(d.dirty)

	return diff
}

func (d *DiffMap[K, V]) Set(key K, val V) {
	existing, found := d.committed[key]
	if found && d.comparator(existing, val) {
		delete(d.dirty, key)
		return
	}

	d.dirty[key] = val
}
