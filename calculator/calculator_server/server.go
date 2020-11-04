package main

import (
	"context"
	"log"
	"net"

	"github.com/JulianToledano/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Calculate(ctx context.Context, in *calculatorpb.CalculateRequest) (*calculatorpb.CalculateResponse, error) {
	log.Println("Calculate function was invoked with %v", in)
	a := in.GetNumbers().GetFirstNumber()
	b := in.GetNumbers().GetSecondNumber()

	result := a + b
	res := &calculatorpb.CalculateResponse{
		Result: result,
	}
	return res, nil
}

func (*server) PrimeDecomposition(req *calculatorpb.PrimeRequest, stream calculatorpb.CalculatorService_PrimeDecompositionServer) error {
	log.Println("Starting Prime computation...")
	pn := req.GetPrimeNumber().GetNumber()
	var k int32 = 2
	for pn > 1 {
		log.Println("Prime computation iteration...")
		if pn%k == 0 {
			res := &calculatorpb.PrimeResponse{
				DecomposedPrime: k,
			}
			stream.Send(res)
			pn /= k
		} else {
			k++
		}
	}
	return nil
}

func main() {
	log.Println("Starting gRPC server...")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalln("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalln("Failed to serve: %v", err)
	}
}
