package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/snapshot-chromedp/render"
)

type cacheInfo struct {
	cacheName string
	opsPerSec int
}

func main() {
	path := os.Args[1]
	dir := filepath.Dir(path)

	perfFile, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer perfFile.Close()

	scanner := bufio.NewScanner(perfFile)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	workloadToCaches := make(map[string][]cacheInfo)
	for _, line := range lines[3 : len(lines)-2] {
		fields := strings.Fields(line)
		opsPerSec, err := strconv.Atoi(fields[4])
		if err != nil {
			log.Fatal("can not parse benchmark output")
		}

		benchInfo := strings.Split(fields[0], "/")[1]
		benchParts := strings.Split(benchInfo, "_")
		cacheName := benchParts[1]
		workload := strings.Split(benchParts[2], "-")[0]

		workloadToCaches[workload] = append(workloadToCaches[workload], cacheInfo{
			cacheName: cacheName,
			opsPerSec: opsPerSec,
		})
	}

	for workload, caches := range workloadToCaches {
		bar := charts.NewBar()
		bar.SetGlobalOptions(
			charts.WithYAxisOpts(opts.YAxis{
				Name: "ops/s",
			}),
			charts.WithTitleOpts(opts.Title{
				Title: workload,
				Right: "40%",
			}),
			charts.WithLegendOpts(opts.Legend{
				Orient: "vertical",
				Right:  "0%",
				Top:    "10%",
			}),
			// for png render
			charts.WithAnimation(false),
		)

		bar = bar.SetXAxis([]string{"cache"})
		bar.BackgroundColor = "white"
		for _, cache := range caches {
			bar = bar.AddSeries(cache.cacheName, []opts.BarData{
				{
					Value: cache.opsPerSec,
				},
			})
		}

		outputName := strings.Join(strings.Split(workload, "%"), "")
		imagePath := filepath.Join(dir, fmt.Sprintf("%s.png", outputName))
		render.MakeChartSnapshot(bar.RenderContent(), imagePath)
	}
}
