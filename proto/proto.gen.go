package proto

//go:generate protoc --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --go_out=../base --go-grpc_out=../base -I . util/plugin/info.proto
//go:generate protoc --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --go_out=../base --go-grpc_out=../base -I . util/config/config.proto
//go:generate protoc --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --go_out=../base --go-grpc_out=../base -I . util/error/error.proto
//go:generate protoc --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --go_out=../base --go-grpc_out=../base -I . data_source/data_source.proto
//go:generate protoc --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --go_out=../base --go-grpc_out=../base -I . identity_store/identity_store.proto
//go:generate protoc --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --go_out=../base --go-grpc_out=../base -I . data_usage/data_usage.proto
//go:generate protoc --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --go_out=../base --go-grpc_out=../base -I . access_provider/access_provider.proto
