# Generating code using protoc

- Install the `protoc` compiler
- Build `protoc-gen-go` 
  ```bash
  go build google.golang.org/protobuf/cmd/protoc-gen-go
  ```
- Generate code
  ```bash
  protoc --go_out=. ./proto/message/message.proto
  ```
