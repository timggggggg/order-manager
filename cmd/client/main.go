package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"gitlab.ozon.dev/timofey15g/homework/internal/client"
	desc "gitlab.ozon.dev/timofey15g/homework/pkg/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
		return
	}

	_, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	conn, err := grpc.NewClient("localhost:5252", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create grpc client: %v", err)
	}

	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		<-signalChan
		conn.Close()
		cancel()
		os.Exit(0)
	}()

	clnt := desc.NewOrderServiceClient(conn)

	app := client.NewApp(clnt)
	app.Run()
}
