# Generating code using protoc ([Quick start](https://grpc.io/docs/languages/go/quickstart/))

- Install the `protoc` compiler
- Build `protoc-gen-go` and `protoc-gen-go-grpc` 
  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
  ```
- Generate code
  ```bash
  protoc --go_out=. --go-grpc_out=. ./proto/sorting/sorting.proto
  ```