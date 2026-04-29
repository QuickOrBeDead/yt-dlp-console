#!/bin/bash
set -e
trap 'rm -f coverage.out coverage_filtered.out' ERR

go test -coverprofile=coverage.out ./...
if [ -f .coverignore ]; then
  grep -v -f .coverignore coverage.out > coverage_filtered.out
  go tool cover -func=coverage_filtered.out | tail -n 1
else
  go tool cover -func=coverage.out | tail -n 1
fi
rm -f coverage.out coverage_filtered.out
