package stores

import (
	"sync"
	"time"
)

type MapSyncStore struct {
	m sync.Map
}

func (x *MapSyncStore) Save(key string, data string, expires time.Duration) error {
	time.AfterFunc(expires, func() {
		x.m.Delete(key)
	})
	x.m.Store(key, data)
	return nil
}

func (x *MapSyncStore) Load(key string) (string, error) {
	value, ok := x.m.Load(key)
	if !ok {
		return "", nil
	}
	return value.(string), nil
}

func (x *MapSyncStore) Delete(key string) error {
	x.m.Delete(key)
	return nil
}
