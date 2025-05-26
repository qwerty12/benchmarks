package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/llxisdsh/pb"
	"github.com/maypok86/otter"
	"github.com/puzpuzpuz/xsync/v4"
	"github.com/viccon/sturdyc"
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
		"otter":      newOtter,
		"sturdyc": newSturdyc,
		"xsync": newXsync,
		"pb": newPb,
	}[name]
	if !ok {
		log.Fatalf("not found cache %s\n", name)
	}

	var o runtime.MemStats
	runtime.ReadMemStats(&o)

	constructor(capacity)

	// runtime.GC()

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
	cache, err := otter.MustBuilder[string, struct{}](capacity).
		WithTTL(time.Hour).
		Build()
	if err != nil {
		log.Fatal(err)
	}
	for _, key := range keys {
		cache.Set(key, struct{}{})
		for i := 0; i < 10; i++ {
			cache.Get(key)
		}
		time.Sleep(5 * time.Microsecond)
	}
}

func newSturdyc(capacity int) {
	cache := sturdyc.New[struct{}](capacity, 16, time.Hour, 10)
	for _, key := range keys {
		cache.Set(key, struct{}{})
		for i := 0; i < 10; i++ {
			cache.Get(key)
		}
		time.Sleep(5 * time.Microsecond)
	}
}

func newXsync(capacity int) {
	cache := xsync.NewMap[string, struct{}](xsync.WithPresize(capacity))
	for _, key := range keys {
		cache.Store(key, struct{}{})
		for i := 0; i < 10; i++ {
			cache.Load(key)
		}
		time.Sleep(5 * time.Microsecond)
	}
}

func newPb(capacity int) {
	cache := pb.NewMapOf[string, struct{}](pb.WithShrinkEnabled(), pb.WithPresize(capacity))
	for _, key := range keys {
		cache.Store(key, struct{}{})
		for i := 0; i < 10; i++ {
			cache.Load(key)
		}
		time.Sleep(5 * time.Microsecond)
	}
}
