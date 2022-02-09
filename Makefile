all: internal/eth/types_easyjson.go \
internal/protobuf/storage/storage.pb.go \
internal/protobuf/storage/storage_vtproto.pb.go \
cmd/config_easyjson.go

internal/eth/types_easyjson.go: internal/eth/types.go
	easyjson -all internal/eth/types.go

cmd/config_easyjson.go: cmd/config.go
	easyjson -all cmd/config.go

internal/protobuf/storage/storage.pb.go internal/protobuf/storage/storage_vtproto.pb.go: proto/storage.proto
	protoc \
    --go_out=. \
    --go_opt=module=github.com/lukasmoellerch/mensa-cli \
    --go-vtproto_out=. \
    --go-vtproto_opt=module=github.com/lukasmoellerch/mensa-cli \
    --go-vtproto_opt=features=marshal+unmarshal+size \
    proto/storage.proto; \