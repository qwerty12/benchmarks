package client

import (
	"github.com/dgraph-io/ristretto"
)

type Ristretto[K comparable, V any] struct {
	client *ristretto.Cache
}

func (c *Ristretto[K, V]) Init(cap int) {
	client, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: int64(cap * 10),
		MaxCost:     int64(cap),
		BufferItems: 64,
	})
	if err != nil {
		panic(err)
	}
	c.client = client
}

func (c *Ristretto[K, V]) Name() string {
	return "Ristretto"
}

func (c *Ristretto[K, V]) Get(key K) (V, bool) {
	v, ok := c.client.Get(key)
	if ok {
		return v.(V), true
	}
	var zero V
	return zero, false
}

func (c *Ristretto[K, V]) Set(key K, value V) {
	c.client.Set(key, value, 1)
}

func (c *Ristretto[K, V]) Close() {
	c.client.Close()
	c.client = nil
}
