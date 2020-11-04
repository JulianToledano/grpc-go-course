package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/JulianToledano/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	defer cc.Close()
	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}

	c := greetpb.NewGreetServiceClient(cc)
	// doUnary(c)

	// doServerStreaming(c)
	doClientStreaming(c)

}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Printf("Start to do a Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Julián",
			LastName:  "Toledano",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Printf("error while calling Greet RPC: %v", err)
	}

	log.Printf("Response from Greet: %v", res)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	log.Println("Starting to do a Server streaming")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Julián",
			LastName:  "Toledano",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalln("Error while calling GreatmanyTimes RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalln("Error while reading stream")
		}
		log.Println("Response from GreetManyTimes: %v", msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	log.Println("Starting to do a client streaming rpc...")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Julián",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Lucy",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Mark",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Piper",
			},
		},
	}
	sstream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalln("Error while calling LongGreet: %v", err)
	}
	for _, req := range requests {
		log.Println("Sending req: %v", req)
		sstream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}
	res, err := sstream.CloseAndRecv()
	if err != nil {
		log.Fatalln("Error while receiving")
	}
	log.Println("LongGreet response: %v", res)
}
