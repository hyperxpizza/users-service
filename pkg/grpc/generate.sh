#! /bin/sh

echo "generating proto from users-service.proto"
protoc --go_out=.  --proto_path=. --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative users-service.proto