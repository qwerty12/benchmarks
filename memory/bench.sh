#!/bin/bash

set -e

caches=(
  "otter"
  "ristretto"
  "theine"
  "hashicorp"
  "ccache"
)

capacities=(10000 100000 1000000)

for cache in "${caches[@]}"
do
  for capacity in "${capacities[@]}"
  do
    go run main.go "$cache" "$capacity"
  done
done

printf "\nDONE!\n"
