package main

import (
	"fmt"
	"math"
	"runtime"
	"strings"

	"github.com/maypok86/benchmarks/client"
)

func getGeneratorByName(name string) generator {
	switch strings.ToLower(name) {
	case "zipf":
		return newZipf()
	case "s3":
		return newARC("s3")
	case "ds1":
		return newARC("ds1")
	case "p3":
		return newARC("p3")
	case "p8":
		return newARC("p8")
	case "loop":
		return newLIRS("loop")
	case "oltp":
		return newARC("oltp")
	default:
		panic("not found generator")
	}
}

func runBench(name string, caps []int) {
	benchClients := []client.Client[string, string]{
		&client.Otter[string, string]{},
		&client.Theine[string, string]{},
		&client.Ristretto[string, string]{},
		&client.LRU[string, string]{},
		&client.ARC[string, string]{},
	}

	fmt.Printf("%s:\n", name)
	for _, capacity := range caps {
		fmt.Printf("\tCapacity: %d\n", capacity)
		benches := []bench{ /*newBenchOptimal(capacity)*/ }
		for _, c := range benchClients {
			c.Init(capacity)
			benches = append(benches, newBenchClient(c))
		}

		runtime.GC()

		for _, b := range benches {
			var gen generator
			operations := uint64(math.MaxUint64)
			gen = getGeneratorByName(name)
			if name == "Zipf" {
				operations = 1_000_000
			}

			for i := uint64(0); i < operations; i++ {
				k, err := gen()
				if err != nil {
					if err == errDone {
						break
					}
					panic(err)
				}

				b.Set(fmt.Sprintf("%d", k))
			}

			fmt.Printf("\t\t%s simulator: %0.2f\n", b.Name(), b.Ratio())
			b.Close()
		}
	}
	runtime.GC()
}

func main() {
	runBench("Zipf", []int{500, 1000, 2000, 5000, 10000, 20000, 40000, 80000})
	runBench("S3", []int{100_000, 200_000, 300_000, 400_000, 500_000, 600_000, 700_000, 800_000})
	runBench("DS1", []int{1_000_000, 2_000_000, 3_000_000, 4_000_000, 5_000_000, 6_000_000, 7_000_000, 8_000_000})
	runBench("P3", []int{25_000, 50_000, 100_000, 200_000, 300_000, 400_000, 500_000, 600_000})
	runBench("P8", []int{10_000, 20_000, 30_000, 40_000, 50_000, 60_000, 70_000, 80_000})
	runBench("LOOP", []int{250, 500, 750, 1000, 1250, 1500, 1750, 2000})
	runBench("OLTP", []int{250, 500, 750, 1000, 1250, 1500, 1750, 2000})
}
