package main

import (
	"container/heap"

	"github.com/maypok86/benchmarks/client"
)

type bench interface {
	Name() string
	Set(k string)
	Ratio() float64
	Close()
}

type benchClient struct {
	c      client.Client[string, string]
	hits   uint64
	misses uint64
}

func newBenchClient(c client.Client[string, string]) bench {
	return &benchClient{
		c: c,
	}
}

func (b *benchClient) Set(k string) {
	v, ok := b.c.Get(k)
	if ok {
		if v != k {
			panic("not valid value")
		}
		b.hits++
	} else {
		b.c.Set(k, k)
		b.misses++
	}
}

func (b *benchClient) Ratio() float64 {
	return 100 * (float64(b.hits) / float64(b.hits+b.misses))
}

func (b *benchClient) Name() string {
	return b.c.Name()
}

func (b *benchClient) Close() {
	b.c.Close()
}

type benchOptimal struct {
	capacity uint64
	hits     map[string]uint64
	access   []string
}

func newBenchOptimal(capacity int) bench {
	return &benchOptimal{
		capacity: uint64(capacity),
		hits:     make(map[string]uint64),
		access:   make([]string, 0),
	}
}

func (b *benchOptimal) Set(key string) {
	b.hits[key]++
	b.access = append(b.access, key)
}

func (b *benchOptimal) Ratio() float64 {
	hits := uint64(0)
	misses := uint64(0)
	look := make(map[string]struct{}, b.capacity)
	data := &optimalHeap{}
	heap.Init(data)
	for _, key := range b.access {
		if _, has := look[key]; has {
			hits++
			continue
		}
		if uint64(data.Len()) >= b.capacity {
			victim := heap.Pop(data)
			delete(look, victim.(*optimalItem).key)
		}
		misses++
		look[key] = struct{}{}
		heap.Push(data, &optimalItem{key, b.hits[key]})
	}

	return 100 * (float64(hits) / float64(hits+misses))
}

func (b *benchOptimal) Name() string {
	return "Optimal"
}

func (b *benchOptimal) Close() {
}

type optimalItem struct {
	key  string
	hits uint64
}

type optimalHeap []*optimalItem

func (h optimalHeap) Len() int           { return len(h) }
func (h optimalHeap) Less(i, j int) bool { return h[i].hits < h[j].hits }
func (h optimalHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *optimalHeap) Push(x any) {
	*h = append(*h, x.(*optimalItem))
}

func (h *optimalHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
