# Generate model

- update model (/model/model.proto)
- generate model
```shell
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative model/model.proto
```