package cache

import (
	"fmt"
	//"github.com/muesli/cache2go"
	"github.com/bluele/gcache"
	"sync"
	"time"
)

type MemCache struct {
	cacheTables map[string]gcache.Cache
	mutex       sync.RWMutex
}

func NewMemCache() *MemCache {
	return &MemCache{
		cacheTables: make(map[string]gcache.Cache),
	}
}

func (m *MemCache) Add(table string, key string, sec time.Duration, iTerm interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	t, ok := m.cacheTables[table]
	if !ok {
		cacheTable := gcache.New(20).LFU().Build()
		cacheTable.SetWithExpire(key, iTerm, sec)
		m.cacheTables[table] = cacheTable
	} else {
		t.SetWithExpire(key, iTerm, sec)
	}
}

func (m *MemCache) Value(table, key string) (interface{}, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	t, ok := m.cacheTables[table]
	if !ok {
		return nil, fmt.Errorf("not table %s", table)
	} else {
		iTerm, err := t.Get(key)
		if err != nil {
			return nil, err
		}
		return iTerm, nil
	}
}
