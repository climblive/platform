package scores

type DiffMap[K comparable, V comparable] struct {
	committed map[K]V
	dirty     map[K]V
}

func NewDiffMap[K comparable, V comparable]() DiffMap[K, V] {
	return DiffMap[K, V]{
		committed: make(map[K]V),
		dirty:     make(map[K]V),
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
	if found && existing == val {
		return
	}

	d.dirty[key] = val
}
