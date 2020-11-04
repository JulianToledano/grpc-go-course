package main

import (
	"context"
	"io"
	"log"

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
	doClientStream(c)
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
