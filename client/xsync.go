package client

import (
	"github.com/puzpuzpuz/xsync/v4"
)

type Xsync[K comparable, V any] struct {
	client *xsync.Map[K, V]
}

func (c *Xsync[K, V]) Init(capacity int) {
	client := xsync.NewMap[K, V](xsync.WithPresize(capacity))
	c.client = client
}

func (c *Xsync[K, V]) Name() string {
	return "xsync"
}

func (c *Xsync[K, V]) Get(key K) (V, bool) {
	return c.client.Load(key)
}

func (c *Xsync[K, V]) Set(key K, value V) {
	c.client.Store(key, value)
}

func (c *Xsync[K, V]) Close() {
	c.client.Clear()
	c.client = nil
}
