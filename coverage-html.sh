#!/bin/bash
set -e
trap 'rm -f coverage.out coverage_filtered.out' ERR

go test -coverprofile=coverage.out ./...
if [ -f .coverignore ]; then
  grep -v -f .coverignore coverage.out > coverage_filtered.out
  go tool cover -html=coverage_filtered.out -o coverage.html
else
  go tool cover -html=coverage.out -o coverage.html
fi
rm -f coverage.out coverage_filtered.out
google-chrome coverage.html &
