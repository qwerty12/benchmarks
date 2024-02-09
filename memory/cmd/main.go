package main

import (
	"bufio"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/snapshot-chromedp/render"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type memoryResult struct {
	cacheName string
	alloc     int
}

func main() {
	path := os.Args[1]
	dir := filepath.Dir(path)

	memoryFile, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer memoryFile.Close()

	scanner := bufio.NewScanner(memoryFile)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	capacityToResults := make(map[int][]memoryResult)
	for _, line := range lines {
		fields := strings.Fields(line)
		cacheName := fields[0]
		capacity, err := strconv.Atoi(fields[1])
		if err != nil {
			log.Fatal("can not parse benchmark output", err)
		}
		alloc, err := strconv.Atoi(fields[2])
		if err != nil {
			log.Fatal("can not parse benchmark output", err)
		}

		capacityToResults[capacity] = append(capacityToResults[capacity], memoryResult{
			cacheName: cacheName,
			alloc:     alloc,
		})
	}

	for capacity, results := range capacityToResults {
		bar := charts.NewBar()
		bar.SetGlobalOptions(
			charts.WithYAxisOpts(opts.YAxis{
				Name: "alloc",
				AxisLabel: &opts.AxisLabel{
					Formatter: "{value} KiB",
				},
			}),
			charts.WithTitleOpts(opts.Title{
				Title: fmt.Sprintf("Memory consumption (%d)", capacity),
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
		for _, res := range results {
			bar = bar.AddSeries(res.cacheName, []opts.BarData{
				{
					Value: res.alloc,
				},
			})
		}

		outputName := fmt.Sprintf("memory_%d", capacity)
		imagePath := filepath.Join(dir, fmt.Sprintf("%s.png", outputName))
		render.MakeChartSnapshot(bar.RenderContent(), imagePath)
	}
}
