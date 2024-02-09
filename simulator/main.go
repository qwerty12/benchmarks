package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/snapshot-chromedp/render"
	"github.com/olekukonko/tablewriter"

	"github.com/maypok86/benchmarks/client"
)

type ratioResult struct {
	capacity int
	ratio    float64
}

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
	dir := "results"
	imagePath := filepath.Join(dir, fmt.Sprintf("%s.png", strings.ToLower(name)))
	if _, err := os.Stat(imagePath); err == nil {
		fmt.Printf("Cached results on %s trace. See %s.\n", name, imagePath)
		return
	}

	fmt.Printf("\nStarted hit ratio simulation on %s trace.\n\n", name)

	clients := []client.Client[string, string]{
		&client.Otter[string, string]{},
		&client.Theine[string, string]{},
		&client.Ristretto[string, string]{},
		&client.LRU[string, string]{},
		&client.ARC[string, string]{},
	}

	var cacheNames []string
	for _, c := range clients {
		cacheNames = append(cacheNames, c.Name())
	}

	cacheToResults := make(map[string][]ratioResult)
	for _, capacity := range caps {
		benches := []bench{ /*newBenchOptimal(capacity)*/ }
		for _, c := range clients {
			c.Init(capacity)
			benches = append(benches, newBenchClient(c))
		}

		tempResults := make(map[string]float64)

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

			cacheToResults[b.Name()] = append(cacheToResults[b.Name()], ratioResult{
				capacity: capacity,
				ratio:    b.Ratio(),
			})
			tempResults[b.Name()] = b.Ratio()
			b.Close()
		}

		stringCapacity := strconv.Itoa(capacity)
		tw := tablewriter.NewWriter(os.Stdout)
		tw.SetHeader([]string{"Cache", "Capacity", "Hit ratio"})
		for _, cacheName := range cacheNames {
			tw.Append([]string{cacheName, stringCapacity, fmt.Sprintf("%0.2f%%", tempResults[cacheName])})
		}
		tw.Render()
		fmt.Println()
	}
	runtime.GC()

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithXAxisOpts(opts.XAxis{
			Name: "capacity",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "hit ratio",
			AxisLabel: &opts.AxisLabel{
				Formatter: "{value}%",
			},
		}),
		charts.WithTitleOpts(opts.Title{
			Title: name,
			Right: "50%",
		}),
		charts.WithLegendOpts(opts.Legend{
			Orient: "vertical",
			Right:  "0%",
			Top:    "10%",
		}),
		// for png render
		charts.WithAnimation(false),
	)

	line = line.SetXAxis(caps)
	line.BackgroundColor = "white"
	for _, cacheName := range cacheNames {
		results := cacheToResults[cacheName]
		var lineData []opts.LineData
		for _, res := range results {
			lineData = append(lineData, opts.LineData{
				Value: res.ratio,
			})
		}
		line = line.AddSeries(cacheName, lineData)
	}

	line.SetSeriesOptions(charts.WithLineChartOpts(
		opts.LineChart{
			Smooth: opts.Bool(true),
		}),
	)

	render.MakeChartSnapshot(line.RenderContent(), imagePath)
	fmt.Println()
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
