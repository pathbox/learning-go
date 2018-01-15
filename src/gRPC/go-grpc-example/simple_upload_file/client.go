package main

import (
	"context"
	"io"
	"log"
	"os"

	pb "./messaging"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

type ClientGRPC struct{}

func main() {
	var (
		writing = true
		buf     []byte
		n       int
		file    *os.File
		status  *pb.UploadStatus
	)
	addr := "127.0.0.1:12345"

	path := "/home/user/download/leetcode-solution.pdf"
	file, err := os.Open(path)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		grpclog.Fatalln(err)
	}
	defer conn.Close()
	ctx := context.Background()

	// 初始化客户端, conn 作为参数
	client := pb.NewGuploadServiceClient(conn)
	stream, err := client.Upload(ctx)

	if err != nil {
		log.Panic(err)
	}

	defer stream.CloseSend()

	buf = make([]byte, 1<<22)
	log.Println("Start upload")
	for writing {
		n, err = file.Read(buf)
		if err != nil {
			if err == io.EOF {
				writing = false
				err = nil
				continue
			}

			err = errors.Wrapf(err,
				"errored while copying from file to buf")
			return
		}

		err = stream.Send(&pb.Chunk{
			Content: buf[:n],
		})
		if err != nil {
			err = errors.Wrapf(err,
				"failed to send chunk via stream")
			return
		}
	}

	status, err = stream.CloseAndRecv()

	if err != nil {
		log.Println(err)
	}

	if status.Code != pb.UploadStatusCode_Ok {
		log.Println("Upload failed:", status.Message)
	}
}
