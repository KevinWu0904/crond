.PHONY: proto

# make proto: provides the quick command to generate go code.
# instructions about protoc installation:
#   1. protoc: install release binaries (current is 3.17.3) directly from https://github.com/protocolbuffers/protobuf/releases.
#   2. grpc plugin: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.0.0 .
#   3. go plugin: go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1 .
proto:
	protoc -I proto/ --go_out=proto/types --go_opt=paths=source_relative --go-grpc_out=proto/types --go-grpc_opt=paths=source_relative proto/crond.proto