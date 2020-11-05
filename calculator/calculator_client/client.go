package main

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/JulianToledano/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	log.Println("gRPC client...")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Could not connect: %v", err)
	}

	c := calculatorpb.NewCalculatorServiceClient(cc)
	// doUnary(c)
	// doStream(c)
	// doClientStream(c)
	doBiNiStream(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	log.Println("Starting to do a Unary RPC...")
	req := &calculatorpb.CalculateRequest{
		Numbers: &calculatorpb.Calculate{FirstNumber: 3,
			SecondNumber: 10},
	}
	res, err := c.Calculate(context.Background(), req)
	if err != nil {
		log.Fatalln("Error while calling Calculate RPC: %v", err)
	}

	log.Println("Response from Calculate: %v", res)
}

func doStream(c calculatorpb.CalculatorServiceClient) {
	log.Println("Starting to do a server streaming...")

	req := &calculatorpb.PrimeRequest{
		PrimeNumber: &calculatorpb.Prime{
			Number: 120,
		},
	}
	resStream, err := c.PrimeDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalln("Error while calling PrimeDecomposition RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("Error while calculating primes")
		}
		log.Println("Response from PrimeDecomposition: %v", msg.GetDecomposedPrime())
	}
}

func doClientStream(c calculatorpb.CalculatorServiceClient) {
	log.Println("Starring to do a client streaming rpc...")

	requests := []*calculatorpb.AverageRequest{
		&calculatorpb.AverageRequest{
			Number: 1,
		},
		&calculatorpb.AverageRequest{
			Number: 2,
		},
		&calculatorpb.AverageRequest{
			Number: 3,
		},
		&calculatorpb.AverageRequest{
			Number: 4,
		},
	}

	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalln("Error while calling Average: %v", err)
	}

	for _, req := range requests {
		log.Println("Sending req: %v", req)
		stream.Send(req)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln("Error while receiving")
	}
	log.Println("Average result: %v", res)
}

func doBiNiStream(c calculatorpb.CalculatorServiceClient) {
	log.Println("Starring to do a client streaming rpc...")

	requests := []*calculatorpb.MaxRequest{
		&calculatorpb.MaxRequest{
			Number: 1,
		},
		&calculatorpb.MaxRequest{
			Number: 5,
		},
		&calculatorpb.MaxRequest{
			Number: 3,
		},
		&calculatorpb.MaxRequest{
			Number: 6,
		},
		&calculatorpb.MaxRequest{
			Number: 2,
		},
		&calculatorpb.MaxRequest{
			Number: 20,
		},
	}
	stream, err := c.Maximum(context.Background())
	if err != nil {
		log.Fatalln("Error while creating the stream: %v", err)
		return
	}
	waitc := make(chan struct{})

	go func() {
		for _, req := range requests {
			log.Println("Sending message: %v", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				break
			}
			if err != nil {
				log.Fatalln("Error while receiving: %v", err)
				close(waitc)
				break
			}
			log.Println("Received: %v", res)
		}
	}()
	<-waitc
}
