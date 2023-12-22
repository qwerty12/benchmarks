package client

import (
	"github.com/VictoriaMetrics/fastcache"
)

type Fastcache struct {
	client *fastcache.Cache
}

func (c *Fastcache) Init(cap int) {
	client := fastcache.New(128 * cap)
	c.client = client
}

func (c *Fastcache) Name() string {
	return "Fastcache"
}

func (c *Fastcache) Get(key string) (string, bool) {
	v := c.client.Get(nil, s2b(key))
	if len(v) == 0 {
		return "", false
	}
	return b2s(v), true
}

func (c *Fastcache) Set(key string, value string) {
	c.client.Set(s2b(key), s2b(value))
}

func (c *Fastcache) Close() {
	c.client.Reset()
	c.client = nil
}
