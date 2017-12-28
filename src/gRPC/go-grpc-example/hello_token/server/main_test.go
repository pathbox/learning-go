package main

import (
	"testing"

	pb "../../proto/hello_tls"
	"golang.org/x/net/context"
)

func HelloTest(t *testing.T) {
	s := helloService{}

	// set up test cases
	tests := []struct {
		name string
		want string
	}{
		{
			name: "world",
			want: "Hello world",
		},
		{
			name: "123",
			want: "Hello 123",
		},
	}

	for _, tt := range tests {
		req := &pb.HelloRequest{Name: tt.name}
		resp, err := s.SayHello(context.Background(), req)
		if err != nil {
			t.Errorf("HelloTest(%v) got unexpected error")
		}
		if resp.Message != tt.want {
			t.Errorf("HelloText(%v)=%v, wanted %v", tt.name, resp.Message, tt.want)
		}
	}
}
