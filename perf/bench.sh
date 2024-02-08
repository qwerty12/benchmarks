#!/bin/bash

set -e

perf_path="./results/perf.txt"

go test -run='^$' -cpu=8 -bench . -timeout=0 > "$perf_path"
go run ./cmd/main.go "$perf_path"

