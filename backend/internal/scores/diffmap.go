package scores

import "sync"

type DiffMap[K comparable, V any] struct {
	mu         sync.Mutex
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
	d.mu.Lock()
	defer d.mu.Unlock()

	var diff []V

	for k, v := range d.dirty {
		diff = append(diff, v)
		d.committed[k] = v
	}

	clear(d.dirty)

	return diff
}

func (d *DiffMap[K, V]) Set(key K, val V) {
	d.mu.Lock()
	defer d.mu.Unlock()

	existing, found := d.committed[key]
	if found && d.comparator(existing, val) {
		delete(d.dirty, key)
		return
	}

	d.dirty[key] = val
}
