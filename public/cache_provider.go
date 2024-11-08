package public

import "sync"

type Store interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Del(key string) error
	SID() string
}

type CacheProvider interface {
	Init(string) (Store, error)
	Get(string) (Store, error)
	Remove(string) error
	GC(maxLifeTime int64)
	Update(string)
}

type MemoryCacheProvider struct {
	lock sync.Mutex
}
