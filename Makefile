# To try different version of Go
GO := go

gotestsum := go run gotest.tools/gotestsum@latest

protobuf-generate:
	buf generate

generate: protobuf-generate
	go generate ./...

build: generate
	go build -o raito main.go

test:
	$(gotestsum) --debug --format pkgname -- -mod=readonly -tags=integration -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt ./...
	go tool cover -html=coverage.txt -o coverage.html

lint:
	golangci-lint run ./...
	go fmt ./...


protobuf-lint:
	buf lint

protobuf-format:
	buf format -w

