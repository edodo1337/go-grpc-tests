package main

import (
	"context"
	"flag"
	"grpc-client/pb"
	"log"
	"time"

	"google.golang.org/grpc"
)

var (
	addr = flag.String("addr", "localhost:9000", "the address to connect to")
	name = flag.String("name", "John", "Name to greet")
)

func main() {
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("%v", err.Error())
	}

	defer conn.Close()

	client := pb.NewKVStorageServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	makeRequest := func(key, value string) {
		response, err := client.Put(ctx, &pb.PutRequest{Key: key, Value: value})
		if err != nil {
			log.Fatalf("%v", err.Error())
		}

		log.Println(response)
	}

	go makeRequest("1", "100")
	go makeRequest("2", "200")
	go makeRequest("3", "300")

	time.Sleep(7 * time.Second)
}
