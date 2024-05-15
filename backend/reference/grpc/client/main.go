package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/kartpop/cruncan/backend/reference/grpc/client/conn"
	"github.com/kartpop/cruncan/backend/reference/grpc/model"
)

func main() {
	log.Println("starting grpc client...")

	grpcConn := conn.NewGenericGRPCConnWithContext(context.Background(), "localhost:8443", 10*time.Second)
	grpcClient := model.NewPaymentsServiceClient(grpcConn.Conn)

	resp, err := grpcClient.DoTransaction(context.Background(), &model.TransactionRequest{
		Id:              "sdflj-2342-sdf",
		Amount:          105,
		SourceAccountId: "123",
		TargetAccountId: "456",
	})
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	log.Printf("response: %v\n", resp)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("stopping grpc client...")
}
