package topics

import "sync"

type Map struct {
	sync.RWMutex
	internal map[int]interface{}
}

func NewMap() *Map {
	return &Map{
		internal: make(map[int]interface{}),
	}
}

func (m *Map) Load(key int) (interface{}, bool) {
	m.RLock()
	defer m.RUnlock()
	value, found := m.internal[key]
	return value, found
}

func (m *Map) Delete(key int) {
	m.Lock()
	defer m.Unlock()
	delete(m.internal, key)
}

func (m *Map) Store(key int, value interface{}) {
	m.Lock()
	defer m.Unlock()
	m.internal[key] = value
}
