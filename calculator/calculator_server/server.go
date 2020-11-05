package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"github.com/JulianToledano/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (*server) Maximum(stream calculatorpb.CalculatorService_MaximumServer) error {
	log.Println("Starting Maximum stream...")
	var max int32 = 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("Finished reading client stream")
			return nil
		}
		if err != nil {
			log.Fatalln("Error while reading client stream: %v", err)
			return err
		}
		n := req.GetNumber()
		if n > max {
			max = n
		}
		err = stream.Send(&calculatorpb.MaxResponse{
			Number: max,
		})
		if err != nil {
			log.Fatalln("Error while sending data to client: %v", err)
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	log.Println("Received SquareRoot RPC")
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number: %v", number),
		)
	}
	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
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
