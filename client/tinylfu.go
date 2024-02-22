package client

import (
	"strconv"

	"github.com/dgryski/go-tinylfu"
)

type TinyLFU struct {
	client *tinylfu.T[uint64]
}

func (c *TinyLFU) Init(capacity int) {
	client := tinylfu.New[uint64](capacity, 10*capacity)
	c.client = client
}

func (c *TinyLFU) Name() string {
	return "tinylfu"
}

func (c *TinyLFU) Get(key uint64) (uint64, bool) {
	return c.client.Get(strconv.FormatUint(key, 10))
}

func (c *TinyLFU) Set(key, value uint64) {
	c.client.Add(strconv.FormatUint(key, 10), value)
}

func (c *TinyLFU) Close() {
	c.client = nil
}
