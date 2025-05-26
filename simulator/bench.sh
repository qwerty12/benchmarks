#!/bin/sh

for file in configs/*.toml; do
    echo "$file":
    go run ./cmd/main.go -config "$file"
done