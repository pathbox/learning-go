.PHONY: instal
install: protos/chat.pb.go
				 go install .
				 which grpc-chat

protos/chat.pb.go:
				 protoc --go_out="plugins=grpc:." protos/chat.go
.PHONY: docker
docker:
				docker build --rm -t rodaine/grpc-chat .