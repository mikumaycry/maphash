package maphash

import (
	"hash/maphash"
	"sync"
)

// Hash contains Seed and hasher pool
type Hash struct {
	once sync.Once
	pool *sync.Pool
	seed maphash.Seed
}

// Sum implement bigcache's Hasher func
func (h *Hash) Sum(key string) uint64 {
	h.once.Do(h.initHash)
	item := h.pool.Get().(*maphash.Hash)
	defer h.pool.Put(item)
	item.WriteString(key)
	return item.Sum64()
}

func (h *Hash) initHash() {
	h.seed = maphash.MakeSeed()
	h.pool = &sync.Pool{
		New: func() interface{} {
			var item maphash.Hash
			item.SetSeed(h.seed)
			return &item
		},
	}
}

// New allocate and init Hash
func New() *Hash {
	var h Hash
	h.once.Do(h.initHash)
	return &h
}
