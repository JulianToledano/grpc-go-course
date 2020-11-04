package main

import (
	"context"
	"io"
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

func (*server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	log.Println("Starting Average stream calculation...")
	var result int32 = 0
	var c int32 = 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("Finished reading the client stream")
			return stream.SendAndClose(&calculatorpb.AverageResponse{
				Number: result / c,
			})
		}
		if err != nil {
			log.Fatalln("Error while reading client stream: %v", err)
		}
		n := req.GetNumber()
		result += n
		c++
	}
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
