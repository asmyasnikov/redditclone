all:
		protoc -I/usr/local/include -I. \
			-I${GOPATH}/src \
			-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
			-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway \
			--grpc-gateway_out=logtostderr=true:. \
			--swagger_out=allow_merge=true,merge_file_name=api:./api \
			--go_out=plugins=grpc:. ./api/*.proto