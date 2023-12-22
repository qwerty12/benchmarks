package client

import (
	"github.com/allegro/bigcache"
	"reflect"
	"unsafe"
)

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func s2b(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}

type Bigcache struct {
	client *bigcache.BigCache
}

func (c *Bigcache) Init(cap int) {
	client, err := bigcache.NewBigCache(bigcache.Config{
		Shards:             256,
		LifeWindow:         0,
		MaxEntriesInWindow: cap,
		MaxEntrySize:       128,
		Verbose:            false,
	})
	if err != nil {
		panic(err)
	}
	c.client = client
}

func (c *Bigcache) Name() string {
	return "Bigcache"
}

func (c *Bigcache) Get(key string) (string, bool) {
	v, err := c.client.Get(key)
	if err != nil {
		return "", false
	}
	return b2s(v), true
}

func (c *Bigcache) Set(key string, value string) {
	c.client.Set(key, s2b(value))
}

func (c *Bigcache) Close() {
	c.client.Close()
	c.client.Reset()
	c.client = nil
}
