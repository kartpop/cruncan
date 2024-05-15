package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	pb "github.com/kartpop/cruncan/backend/reference/grpc/model"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedPaymentsServiceServer
}

func (s *server) DoTransaction(ctx context.Context, in *pb.TransactionRequest) (*pb.TransactionResponse, error) {
	log.Printf("Received transaction request id: %v", in.GetId())

	log.Println("Processing transaction...")
	log.Printf("Amount: %v transferred from source account id: %v to target accound id: %v", in.GetAmount(), in.GetSourceAccountId(), in.GetTargetAccountId())

	return &pb.TransactionResponse{Success: true}, nil
}

func main() {

	lis, err := net.Listen("tcp", "localhost:8443")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPaymentsServiceServer(s, &server{})

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	log.Println("grpc server started...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("stopping grpc server...")
	s.GracefulStop()
}
