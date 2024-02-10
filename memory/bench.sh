#!/bin/bash

set -e

caches=(
  "otter"
  "theine"
  "ristretto"
  "ccache"
  "hashicorp"
)

# capacity = 25000 used in caffeine benchmarks
capacities=(10000 25000 100000 1000000)

result_path="./results/memory.txt"

echo -n "" > "$result_path"

for capacity in "${capacities[@]}"
do
  for cache in "${caches[@]}"
  do
    go run main.go "$cache" "$capacity" >> "$result_path"
  done
done

go run ./cmd/main.go "$result_path"
