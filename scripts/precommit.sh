#!/bin/sh
go mod tidy
golangci-lint run --tests=false
go test ./...