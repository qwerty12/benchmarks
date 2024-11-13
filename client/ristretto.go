package client

import (
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/dgraph-io/ristretto/z"
)

type Ristretto[K z.Key, V any] struct {
	client *ristretto.Cache[K, V]
}

func (c *Ristretto[K, V]) Init(capacity int) {
	client, err := ristretto.NewCache[K, V](&ristretto.Config[K, V]{
		NumCounters:        int64(capacity * 10),
		MaxCost:            int64(capacity),
		BufferItems:        64,
		IgnoreInternalCost: true,
	})
	if err != nil {
		panic(err)
	}
	c.client = client
}

func (c *Ristretto[K, V]) Name() string {
	return "ristretto"
}

func (c *Ristretto[K, V]) Get(key K) (V, bool) {
	return c.client.Get(key)
}

func (c *Ristretto[K, V]) Set(key K, value V) {
	c.client.SetWithTTL(key, value, 1, time.Hour)
}

func (c *Ristretto[K, V]) Close() {
	c.client.Close()
	c.client = nil
}
