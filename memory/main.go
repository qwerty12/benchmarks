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
	"github.com/karlseguin/ccache/v3"
	"github.com/maypok86/otter"
)

var keys []string

func main() {
	name := os.Args[1]
	stringCapacity := os.Args[2]
	capacity, err := strconv.Atoi(stringCapacity)
	if err != nil {
		log.Fatal(err)
	}

	keys = make([]string, 0, capacity)
	for i := 0; i < capacity; i++ {
		keys = append(keys, strconv.Itoa(i))
	}

	constructor, ok := map[string]func(int){
		"otter":     newOtter,
		"ristretto": newRistretto,
		"theine":    newTheine,
		"hashicorp": newHashicorp,
		"ccache":    newCcache,
	}[name]
	if !ok {
		log.Fatalf("not found cache %s\n", name)
	}

	var o runtime.MemStats
	runtime.ReadMemStats(&o)

	constructor(capacity)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("%s\t%v KiB\t%v KiB\n",
		name,
		(m.Alloc-o.Alloc)/1024,
		(m.TotalAlloc-o.TotalAlloc)/1024,
	)
}

func newOtter(capacity int) {
	cache, err := otter.MustBuilder[string, int](capacity).
		InitialCapacity(capacity).
		WithTTL(time.Hour).
		Build()
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < capacity; i++ {
		cache.Set(keys[i], i)
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
	for i := 0; i < capacity; i++ {
		cache.SetWithTTL(keys[i], i, 1, time.Hour)
	}
}

func newTheine(capacity int) {
	cache, err := theine.NewBuilder[string, int](int64(capacity)).Build()
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < capacity; i++ {
		cache.SetWithTTL(keys[i], i, 1, time.Hour)
	}
}

func newHashicorp(capacity int) {
	cache := hashicorp.NewLRU[string, int](capacity, nil, time.Hour)
	for i := 0; i < capacity; i++ {
		cache.Add(keys[i], i)
	}
}

func newCcache(capacity int) {
	cache := ccache.New(ccache.Configure[int]().MaxSize(int64(capacity)))
	for i := 0; i < capacity; i++ {
		cache.Set(keys[i], i, time.Hour)
	}
}
