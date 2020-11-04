package main

import (
	"context"
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
	doUnary(c)
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
