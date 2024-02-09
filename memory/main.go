package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/Yiling-J/theine-go"
	"github.com/dgraph-io/ristretto"
	hashicorp "github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/maypok86/otter"
)

var keys []uint32

func main() {
	name := os.Args[1]
	stringCapacity := os.Args[2]
	capacity, err := strconv.Atoi(stringCapacity)
	if err != nil {
		log.Fatal(err)
	}

	keys = make([]uint32, 0, capacity)
	for i := 0; i < capacity; i++ {
		keys = append(keys, uint32(i))
	}

	constructor, ok := map[string]func(int){
		"otter":     newOtter,
		"ristretto": newRistretto,
		"theine":    newTheine,
		"hashicorp": newHashicorp,
	}[name]
	if !ok {
		log.Fatalf("not found cache %s\n", name)
	}

	var o runtime.MemStats
	runtime.ReadMemStats(&o)

	constructor(capacity)

	runtime.GC()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("%s\t%d\t%v KiB\t%v KiB\n",
		name,
		capacity,
		(m.Alloc-o.Alloc)/1024,
		(m.TotalAlloc-o.TotalAlloc)/1024,
	)
}

func newOtter(capacity int) {
	cache, err := otter.MustBuilder[uint32, uint32](capacity).
		WithTTL(time.Hour).
		Build()
	if err != nil {
		log.Fatal(err)
	}
	for _, key := range keys {
		cache.Set(key, key)
	}
}

func newRistretto(capacity int) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 10 * int64(capacity),
		MaxCost:     int64(capacity),
		BufferItems: 64,
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, key := range keys {
		cache.SetWithTTL(key, key, 1, time.Hour)
	}
}

func newTheine(capacity int) {
	cache, err := theine.NewBuilder[uint32, uint32](int64(capacity)).Build()
	if err != nil {
		log.Fatal(err)
	}
	for _, key := range keys {
		cache.SetWithTTL(key, key, 1, time.Hour)
	}
}

func newHashicorp(capacity int) {
	cache := hashicorp.NewLRU[uint32, uint32](capacity, nil, time.Hour)
	for _, key := range keys {
		cache.Add(key, key)
	}
}
