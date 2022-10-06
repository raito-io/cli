# To try different version of Go
GO := go

generate:
	go generate ./...

build: generate
	 go build ./...

test:
	go test -mod=readonly -tags=integration -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.out.tmp ./...
	cat coverage.out.tmp | grep -v "/mock_" | grep -v "/mocks/" > coverage.txt #IGNORE MOCKS
	go tool cover -html=coverage.txt -o coverage.html