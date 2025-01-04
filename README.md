# Compile protobuf:
```bash
protoc -I=./internal/protobuf/ --go_out=./internal/protobuf/ ./internal/protobuf/message.proto
```


## future place:
  - add TLS support
  - old message auto deletion
  - filter on messages
  - gui (most likely bubble)