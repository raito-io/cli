# To try different version of Go
GO := go

generate:
	go generate ./...

build: generate
	go build ./...

test:
	go test -mod=readonly -tags=integration -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt ./...
	go tool cover -html=coverage.txt -o coverage.html

lint:
	golangci-lint run ./...
	go fmt ./...


protobuf-lint:
	buf format -w
	buf lint

protobuf-generate:
	go generate  -run protoc ./...