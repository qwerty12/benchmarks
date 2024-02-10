package main

import (
	"fmt"
	"github.com/karlseguin/ccache/v3"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/Yiling-J/theine-go"
	"github.com/dgraph-io/ristretto"
	hashicorp "github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/maypok86/otter"
)

var keys []string

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func toMB(bytes uint64) float64 {
	return float64(bytes) / 1024 / 1024
}

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

	runtime.GC()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("%s\t%d\t%v MB\t%v MB\n",
		name,
		capacity,
		toFixed(toMB(m.Alloc-o.Alloc), 2),
		toFixed(toMB(m.TotalAlloc-o.TotalAlloc), 2),
	)
}

func newOtter(capacity int) {
	cache, err := otter.MustBuilder[string, string](capacity).
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
	cache, err := theine.NewBuilder[string, string](int64(capacity)).Build()
	if err != nil {
		log.Fatal(err)
	}
	for _, key := range keys {
		cache.SetWithTTL(key, key, 1, time.Hour)
	}
}

func newCcache(capacity int) {
	cache := ccache.New(ccache.Configure[string]().MaxSize(int64(capacity)))
	for _, key := range keys {
		cache.Set(key, key, time.Hour)
	}
}

func newHashicorp(capacity int) {
	cache := hashicorp.NewLRU[string, string](capacity, nil, time.Hour)
	for _, key := range keys {
		cache.Add(key, key)
	}
}
