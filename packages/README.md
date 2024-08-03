# Packages

## Struct generation

```bash
protoc \
--go_out=gen \
--go_opt=paths=source_relative \
--go-grpc_out=gen \
--go-grpc_opt=paths=source_relative \
--experimental_allow_proto3_optional \
protos/svc_account.proto
```
