package protocol

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative transactions.proto
//go:generate protoc --descriptor_set_out=./transactions_descriptor.pb transactions.proto
