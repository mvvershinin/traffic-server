package main

import (
	"context"
	"fmt"
	"github/mvvershinin/writer_server/config"
	"github/mvvershinin/writer_server/pkg/proto/hello"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type Server struct {
	hello.UnimplementedGreeterServer
}

func (s *Server) SendData(ctx context.Context, in *hello.DataRequest) (*hello.DataReply, error) {
	return &hello.DataReply{Message: "Received data: " + in.Data}, nil
}

func (s *Server) SayHello(ctx context.Context, in *hello.HelloRequest) (*hello.HelloReply, error) {
	fmt.Println("message")
	fmt.Println(in.String())
	fmt.Println(in.Name)

	data, _ := proto.Marshal(&hello.HelloReply{Message: "Hello, " + in.Name})
	fmt.Println(data)
	return &hello.HelloReply{Message: "Hello, " + in.Name}, nil
}

func (s *Server) FetchData(req *hello.DataRequest, srv hello.Greeter_FetchDataServer) error {
	log.Printf("fetch response for id : %d", req.Data)

	var wg sync.WaitGroup
	for i := 0; i < 15; i++ {
		wg.Add(10)
		go func(count int64) {
			defer wg.Done()

			//time sleep to simulate server process time
			time.Sleep(time.Duration(count) * time.Second)
			resp := hello.DataReply{Message: fmt.Sprintf("Request #%d For Id:%d", count, req.Data)}
			if err := srv.Send(&resp); err != nil {
				log.Printf("send error %v", err)
			}
			log.Printf("finishing request number : %d", count)
		}(int64(i))
	}

	wg.Wait()
	return nil
}

func (s *Server) ReceiveData(srv hello.Greeter_ReceiveDataServer) error {
	for {
		//var in hello.DataRequest
		in, err := srv.Recv()
		if in != nil {
			log.Printf("Get data : %s", in)
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	cfg := config.Configure()
	fmt.Println("we run writer-server!")
	lis, err := net.Listen("tcp", cfg.SrvAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	srv := &Server{}

	hello.RegisterGreeterServer(s, srv)
	log.Printf("Server listening at %v", lis.Addr())
	//err = s.Serve(lis)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
