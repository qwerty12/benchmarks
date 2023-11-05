package client

import (
	"github.com/Yiling-J/theine-go"
)

type Theine[K comparable, V any] struct {
	client *theine.Cache[K, V]
}

func (c *Theine[K, V]) Init(cap int) {
	client, err := theine.NewBuilder[K, V](int64(cap)).Build()
	if err != nil {
		panic(err)
	}
	c.client = client
}

func (c *Theine[K, V]) Name() string {
	return "Theine"
}

func (c *Theine[K, V]) Get(key K) (V, bool) {
	return c.client.Get(key)
}

func (c *Theine[K, V]) Set(key K, value V) {
	c.client.Set(key, value, 1)
}

func (c *Theine[K, V]) Close() {
	c.client.Close()
}
