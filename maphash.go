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

// Sum64 implement bigcache's Hasher func
func (h *Hash) Sum64(key string) uint64 {
	h.once.Do(h.initHash)
	item := h.pool.Get().(*maphash.Hash)
	item.WriteString(key)
	res := item.Sum64()
	item.Reset()
	h.pool.Put(item)
	return res
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
