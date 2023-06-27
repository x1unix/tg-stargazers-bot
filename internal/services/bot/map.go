package bot

import "sync"

type Map[K, V any] struct {
	m *sync.Map
}

func NewMap[K, V any]() Map[K, V] {
	return Map[K, V]{
		m: &sync.Map{},
	}
}

func (m Map[K, V]) Get(key K) (val V, ok bool) {
	v, ok := m.m.Load(key)
	if !ok {
		return val, false
	}

	val, ok = v.(V)
	return val, ok
}

func (m Map[K, V]) Set(key K, val V) {
	m.m.Store(key, val)
}

func (m Map[K, V]) Delete(key K) {
	m.m.Delete(key)
}
