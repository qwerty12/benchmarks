package client

import (
	"github.com/dgryski/go-tinylfu"
	"strconv"
)

type TinyLFU struct {
	client *tinylfu.T[uint64]
}

func (c *TinyLFU) Init(cap int) {
	client := tinylfu.New[uint64](cap, 10*cap)
	c.client = client
}

func (c *TinyLFU) Name() string {
	return "tinylfu"
}

func (c *TinyLFU) Get(key uint64) (uint64, bool) {
	return c.client.Get(strconv.FormatUint(key, 10))
}

func (c *TinyLFU) Set(key uint64, value uint64) {
	c.client.Add(strconv.FormatUint(key, 10), value)
}

func (c *TinyLFU) Close() {
	c.client = nil
}
