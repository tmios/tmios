package device

import (
	"fmt"
	"sync"
)

type deviceManager struct {
	env   map[string]interface{}
	metas map[string]*DeviceMeta

	mutex sync.RWMutex
}

var (
	m *deviceManager
)

func init() {
	m = &deviceManager{
		metas: make(map[string]*DeviceMeta),
		env:   make(map[string]interface{}),
	}
}

func SetEnv(name string, val interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.env[name] = val
}

func GetEnv(name string) interface{} {
	m.mutex.RLock()
	val := m.env[name]
	m.mutex.RUnlock()

	return val
}

func Register(meta *DeviceMeta) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.metas[meta.Model]; ok {
		panic(fmt.Sprintf("device model '%s' register twice.", meta.Model))
	}

	m.metas[meta.Model] = meta
}

func Metas() []*DeviceMeta {
	var metas []*DeviceMeta

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, meta := range m.metas {
		metas = append(metas, meta)
	}
	return metas
}

func GetMeta(modelName string) *DeviceMeta {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.metas[modelName]
}
