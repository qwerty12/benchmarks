package client

import (
	"github.com/llxisdsh/pb"
)

type Pb[K comparable, V any] struct {
	client *pb.MapOf[K, V]
}

func (c *Pb[K, V]) Init(capacity int) {
	client := pb.NewMapOf[K, V](pb.WithShrinkEnabled(), pb.WithPresize(capacity))
	c.client = client
}

func (c *Pb[K, V]) Name() string {
	return "pb"
}

func (c *Pb[K, V]) Get(key K) (V, bool) {
	return c.client.Load(key)
}

func (c *Pb[K, V]) Set(key K, value V) {
	c.client.Store(key, value)
}

func (c *Pb[K, V]) Close() {
	c.client.Clear()
	c.client = nil
}
