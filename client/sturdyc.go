package client

import (
	"time"

	"github.com/viccon/sturdyc"
)

type Sturdyc[V any] struct {
	client *sturdyc.Client[V]
}

func (c *Sturdyc[V]) Init(capacity int) {
	client := sturdyc.New[V](capacity, 16, time.Hour, 10)
	c.client = client
}

func (c *Sturdyc[V]) Name() string {
	return "sturdyc"
}

func (c *Sturdyc[V]) Get(key string) (V, bool) {
	return c.client.Get(key)
}

func (c *Sturdyc[V]) Set(key string, value V) {
	c.client.Set(key, value)
}

func (c *Sturdyc[V]) Close() {
	c.client = nil
}
