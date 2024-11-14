package client

import (
	"time"

	"github.com/cespare/xxhash/v2"
	"github.com/elastic/go-freelru"
)

type Elastic[V any] struct {
	client freelru.Cache[string, V]
}

// more hash function in https://github.com/elastic/go-freelru/blob/main/bench/hash.go
func hashStringXXHASH(s string) uint32 {
	return uint32(xxhash.Sum64String(s)) //nolint:golint,gosec,gocritic
}

func (c *Elastic[V]) Init(capacity int) {
	client, err := freelru.NewSynced[string, V](uint32(capacity), hashStringXXHASH) //nolint:golint,gosec,gocritic
	if err != nil {
		panic(err)
	}
	client.SetLifetime(1 * time.Hour)
	c.client = client
}

func (c *Elastic[V]) Name() string {
	return "elastic"
}

func (c *Elastic[V]) Get(key string) (V, bool) {
	return c.client.Get(key)
}

func (c *Elastic[V]) Set(key string, value V) {
	c.client.Add(key, value)
}

func (c *Elastic[V]) Close() {
	c.client.Purge()
	c.client = nil
}
