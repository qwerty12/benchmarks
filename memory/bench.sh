#!/bin/bash

set -e

caches=(
  "otter"
  "theine"
  "ristretto"
  "hashicorp"
)

# capacity = 25000 used in caffeine benchmarks
capacities=(10000 25000 100000 1000000)

result_path="./results/memory.txt"

echo -n "" > "$result_path"

for cache in "${caches[@]}"
do
  for capacity in "${capacities[@]}"
  do
    GOMAXPROCS=8 go run main.go "$cache" "$capacity" >> "$result_path"
  done
done

go run ./cmd/main.go "$result_path"
