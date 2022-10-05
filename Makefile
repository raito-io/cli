# To try different version of Go
GO := go

generate:
	go generate ./...

build: generate
	 go build ./...