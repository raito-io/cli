# To try different version of Go
GO := go

protobuf-generate:
	buf generate

generate: protobuf-generate
	go generate ./...

build: generate
	go build main.go

test:
	go test -mod=readonly -tags=integration -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt ./...
	go tool cover -html=coverage.txt -o coverage.html

lint:
	golangci-lint run ./...
	go fmt ./...


protobuf-lint:
	buf lint

protobuf-format:
	buf format -w

