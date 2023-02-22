package proto

//go:generate protoc --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --go_out=../util/plugin --go-grpc_out=../util/plugin -I . info.proto
//go:generate protoc --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --go_out=../util/config --go-grpc_out=../util/config -I . config.proto
//go:generate protoc --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --go_out=../util/error --go-grpc_out=../util/error -I . error.proto
//go:generate protoc --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --go_out=../data_source --go-grpc_out=../data_source -I . data_source.proto
