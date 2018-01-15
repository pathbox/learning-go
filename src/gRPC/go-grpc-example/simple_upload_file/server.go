package main

import (
	"io"
	"log"
	"net"
	"os"

	"github.com/pkg/errors"

	pb "./messaging"
	"google.golang.org/grpc"
)

type ServerGRPC struct{}

func main() {
	addr := "127.0.0.1:12345"

	listener, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	pb.RegisterGuploadServiceServer(server, &ServerGRPC{})

	err = server.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *ServerGRPC) Upload(stream pb.GuploadService_UploadServer) (err error) {
	log.Println("Start upload")
	path := "./receive.pdf"
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	for {
		chunk, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				goto END
			}
			err = errors.Wrapf(err,
				"failed unexpectadely while reading chunks from stream")
			return err
		}
		if _, err := file.Write(chunk.Content); err != nil {
			log.Panic(err)
			return err
		}
	}
	log.Println("Upload Received...")

END:
	err = stream.SendAndClose(&pb.UploadStatus{
		Message: "Upload received with success",
		Code:    pb.UploadStatusCode_Ok,
	})

	if err != nil {
		err = errors.Wrapf(err,
			"failed to send status code")
		return
	}
	return
}
